package blockchain

import (
	"github.com/patrickmao1/gosig/types"
	"sync"
	"time"
)

type MsgBuffer struct {
	msgs map[string]*types.SignedMessage
	mu   sync.Mutex
}

func NewMsgBuffer() *MsgBuffer {
	return &MsgBuffer{msgs: make(map[string]*types.SignedMessage)}
}

func (b *MsgBuffer) Put(key string, msg *types.SignedMessage) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.msgs[key] = msg
}

func (b *MsgBuffer) Has(key string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	_, ok := b.msgs[key]
	return ok
}

func (b *MsgBuffer) List() []*types.SignedMessage {
	b.mu.Lock()
	defer b.mu.Unlock()
	var msgs []*types.SignedMessage
	for id, msg := range b.msgs {
		if msg.DDL().Before(time.Now()) {
			delete(b.msgs, id)
		} else {
			msgs = append(msgs, msg)
		}
	}
	return msgs
}
