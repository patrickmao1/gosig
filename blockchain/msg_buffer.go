package blockchain

import (
	"fmt"
	"github.com/patrickmao1/gosig/crypto"
	"github.com/patrickmao1/gosig/types"
	"google.golang.org/protobuf/proto"
	"sync"
	"time"
)

type OutboundMsgBuffer struct {
	myPrivkey  []byte
	myValIndex uint32

	msgs map[string]*types.Envelope
	mu   sync.Mutex
}

func NewOutboundMsgBuffer(myPrivKey []byte, myValIndex uint32) *OutboundMsgBuffer {
	return &OutboundMsgBuffer{
		myPrivkey:  myPrivKey,
		myValIndex: myValIndex,
		msgs:       make(map[string]*types.Envelope),
	}
}

func (b *OutboundMsgBuffer) Put(key string, msg *types.Message) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.msgs[key] = b.sign(msg)
}

func (b *OutboundMsgBuffer) Get(key string) (*types.Message, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	msg, ok := b.msgs[key]
	return msg.Msg, ok
}

func (b *OutboundMsgBuffer) Tx(doInTx func(buf *OutboundMsgBuffer) error) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return doInTx(b)
}

func (b *OutboundMsgBuffer) Has(key string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	_, ok := b.msgs[key]
	return ok
}

func (b *OutboundMsgBuffer) List() []*types.Envelope {
	b.mu.Lock()
	defer b.mu.Unlock()
	var msgs []*types.Envelope
	for id, msg := range b.msgs {
		if msg.DDL().Before(time.Now()) {
			delete(b.msgs, id)
		} else {
			msgs = append(msgs, msg)
		}
	}
	return msgs
}

func (b *OutboundMsgBuffer) sign(msg *types.Message) *types.Envelope {
	bs, err := proto.Marshal(msg)
	if err != nil {
		panic(fmt.Errorf("populateSignature: failed to marshal msg type %T: %s", msg.Message, err.Error()))
	}
	sig := crypto.SignBytes(b.myPrivkey, bs)
	return &types.Envelope{
		Sig:            sig,
		ValidatorIndex: b.myValIndex,
		Msg:            msg,
	}
}

type InboundMsgBuffer struct {
	vals Validators

	msgs       []*types.Envelope
	mu         sync.Mutex
	hasMsgCond *sync.Cond
}

func NewInboundMsgBuffer(vals Validators) *InboundMsgBuffer {
	b := &InboundMsgBuffer{
		vals: vals,
	}
	b.hasMsgCond = sync.NewCond(&b.mu)
	return b
}

func (b *InboundMsgBuffer) Enqueue(msgs []*types.Envelope) {
	var signedMsgs []*types.Envelope
	for _, msg := range msgs {
		ok, err := b.checkSig(msg)
		if err != nil {
			log.Errorf("failed to enqueue: sig verification failed for msg %+v, err: %s", msg, err.Error())
		}
		if ok {
			signedMsgs = append(signedMsgs, msg)
		}
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.msgs = append(b.msgs, msgs...)
	b.hasMsgCond.Signal()
}

func (b *InboundMsgBuffer) DequeueAll() []*types.Envelope {
	b.mu.Lock()
	defer b.mu.Unlock()

	for len(b.msgs) == 0 {
		b.hasMsgCond.Wait()
	}

	msgs := b.msgs
	b.msgs = nil

	return msgs
}

func (b *InboundMsgBuffer) checkSig(signedMsg *types.Envelope) (bool, error) {
	pubKey := b.vals[signedMsg.ValidatorIndex].GetPubKey()
	msgBytes, err := proto.Marshal(signedMsg.Msg)
	if err != nil {
		return false, err
	}
	return crypto.VerifySigBytes(pubKey, msgBytes, signedMsg.Sig), nil
}
