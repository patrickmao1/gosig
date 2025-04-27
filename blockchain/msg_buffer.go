package blockchain

import (
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

	msgs      map[string]*types.Message
	deleteDDL map[string]time.Time
	relayDDL  map[string]time.Time
	mu        sync.RWMutex
}

func NewOutboundMsgBuffer(myPrivKey []byte, myValIndex uint32) *OutboundMsgBuffer {
	return &OutboundMsgBuffer{
		myPrivkey:  myPrivKey,
		myValIndex: myValIndex,
		msgs:       make(map[string]*types.Message),
		deleteDDL:  make(map[string]time.Time),
		relayDDL:   make(map[string]time.Time),
	}
}

func (b *OutboundMsgBuffer) Put(key string, msg *types.Message, relayDDL time.Time, deleteDDL ...time.Time) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.msgs[key] = msg

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
	return msg, ok
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

func (b *OutboundMsgBuffer) Pack() *types.Envelope {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.msgs) == 0 {
		return nil
	}

	msgs := &types.Messages{}
	for id, msg := range b.msgs {
		now := time.Now()
		if b.deleteDDL[id].Before(now) {
			delete(b.msgs, id)
			delete(b.deleteDDL, id)
		} else {
			if !b.relayDDL[id].Before(now) {
				msgs.Msgs = append(msgs.Msgs, msg)
			}
		}
	}
	bs, err := proto.Marshal(msgs)
	if err != nil {
		log.Errorf("Error marshalling messages: %v", err)
	}
	sig := crypto.SignBytes(b.myPrivkey, bs)

	return &types.Envelope{
		Msgs:           msgs,
		Sig:            sig,
		ValidatorIndex: b.myValIndex,
	}
}

type InboundMsgBuffer struct {
	vals Validators

	msgs       []*types.Message
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

func (b *InboundMsgBuffer) Enqueue(envelope *types.Envelope) {
	utils.LogExecTime(time.Now(), "Enqueue %d msgs", len(envelope.Msgs.Msgs))
	ok, err := b.checkSig(envelope)
	if err != nil {
		log.Errorf("failed to enqueue: err %s", err.Error())
		return
	}
	if !ok {
		log.Infof("msg from vi %d, pubkey %x..: sig verification failed: sig %x..",
			envelope.ValidatorIndex, b.vals.PubKeys()[envelope.ValidatorIndex][:8], envelope.Sig[:8])
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	b.msgs = append(b.msgs, envelope.Msgs.Msgs...)
	b.hasMsgCond.Signal()
}

func (b *InboundMsgBuffer) DequeueAll() []*types.Message {
	b.mu.Lock()
	defer b.mu.Unlock()
	for len(b.msgs) == 0 {
		b.hasMsgCond.Wait()
	}
	msgs := b.msgs
	b.msgs = nil
	return msgs
}

func (b *InboundMsgBuffer) checkSig(envelope *types.Envelope) (bool, error) {
	pubKey := b.vals[envelope.ValidatorIndex].GetPubKey()
	msgBytes, err := proto.Marshal(envelope.Msgs)
	if err != nil {
		return false, err
	}
	pass := crypto.VerifySigBytes(pubKey, msgBytes, envelope.Sig)
	return pass, err
}
