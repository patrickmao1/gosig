package blockchain

import (
	"fmt"
	"github.com/patrickmao1/gosig/crypto"
	"github.com/patrickmao1/gosig/types"
	"github.com/patrickmao1/gosig/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"sync"
	"time"
)

type OutboundMsgBuffer struct {
	myPrivkey  []byte
	myValIndex uint32

	msgs      map[string]*types.Envelope
	deleteDDL map[string]time.Time
	relayDDL  map[string]time.Time
	mu        sync.RWMutex
}

func NewOutboundMsgBuffer(myPrivKey []byte, myValIndex uint32) *OutboundMsgBuffer {
	return &OutboundMsgBuffer{
		myPrivkey:  myPrivKey,
		myValIndex: myValIndex,
		msgs:       make(map[string]*types.Envelope),
		deleteDDL:  make(map[string]time.Time),
		relayDDL:   make(map[string]time.Time),
	}
}

func (b *OutboundMsgBuffer) Put(key string, msg *types.Message, relayDDL time.Time, deleteDDL ...time.Time) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.msgs[key] = b.sign(msg)

	if len(deleteDDL) > 0 {
		b.deleteDDL[key] = deleteDDL[0]
	} else {
		b.deleteDDL[key] = relayDDL
	}
	b.relayDDL[key] = relayDDL
}

func (b *OutboundMsgBuffer) Get(key string) (*types.Message, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	msg, ok := b.msgs[key]
	if !ok {
		return nil, false
	}
	return msg.Msg, ok
}

func (b *OutboundMsgBuffer) Tx(doInTx func(buf *OutboundMsgBuffer) error) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	tx := &OutboundMsgBuffer{
		myPrivkey:  b.myPrivkey,
		myValIndex: b.myValIndex,
		msgs:       b.msgs,
		deleteDDL:  b.deleteDDL,
		relayDDL:   b.relayDDL,
	}
	return doInTx(tx)
}

func (b *OutboundMsgBuffer) Has(key string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, ok := b.msgs[key]
	return ok
}

func (b *OutboundMsgBuffer) List() []*types.Envelope {
	b.mu.Lock()
	defer b.mu.Unlock()
	var msgs []*types.Envelope
	for id, msg := range b.msgs {
		now := time.Now()
		if b.deleteDDL[id].Before(now) {
			delete(b.msgs, id)
			delete(b.deleteDDL, id)
		} else {
			if !b.relayDDL[id].Before(now) {
				msgs = append(msgs, msg)
			}
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
	defer utils.LogExecTime(time.Now(), "Enqueue")

	var signedMsgs []*types.Envelope
	for _, msg := range msgs {
		ok, err := b.checkSig(msg)
		if err != nil {
			log.Errorf("failed to enqueue: err %s", err.Error())
			continue
		}
		if !ok {
			log.Infof("msg from vi %d, pubkey %x..: sig verification failed: sig %x.., msg %s",
				msg.ValidatorIndex, b.vals.PubKeys()[msg.ValidatorIndex][:8], msg.Sig[:8], msg.Msg.ToString())
			continue
		}
		signedMsgs = append(signedMsgs, msg)
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	b.msgs = append(b.msgs, msgs...)
	b.hasMsgCond.Signal()
}

func (b *InboundMsgBuffer) DequeueAll() []*types.Envelope {
	// Perf: maybe buffer and deduplicate msgs

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
	pass := crypto.VerifySigBytes(pubKey, msgBytes, signedMsg.Sig)
	return pass, err
}
