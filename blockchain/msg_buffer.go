package blockchain

import (
	"fmt"
	"github.com/patrickmao1/gosig/crypto"
	"github.com/patrickmao1/gosig/types"
	"google.golang.org/protobuf/proto"
	"sync"
	"time"
)

type MsgBuffer struct {
	myPrivkey  []byte
	myValIndex uint32

	msgs map[string]*types.Envelope
	mu   sync.Mutex
}

func NewMsgBuffer(myPrivKey []byte, myValIndex uint32) *MsgBuffer {
	return &MsgBuffer{
		myPrivkey:  myPrivKey,
		myValIndex: myValIndex,
		msgs:       make(map[string]*types.Envelope),
	}
}

func (b *MsgBuffer) Put(key string, msg *types.Message) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.msgs[key] = b.sign(msg)
}

func (b *MsgBuffer) Get(key string) (*types.Message, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	msg, ok := b.msgs[key]
	return msg.Msg, ok
}

func (b *MsgBuffer) Tx(doInTx func(buf *MsgBuffer) error) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return doInTx(b)
}

func (b *MsgBuffer) Has(key string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	_, ok := b.msgs[key]
	return ok
}

func (b *MsgBuffer) List() []*types.Envelope {
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

func (b *MsgBuffer) sign(msg *types.Message) *types.Envelope {
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
