package blockchain

import (
	"encoding/hex"
	"github.com/patrickmao1/gosig/crypto"
	"github.com/patrickmao1/gosig/types"
	"github.com/stretchr/testify/require"
	"slices"
	"testing"
	"time"
)

func TestInboundMsgBuffer(t *testing.T) {
	sk, pk := crypto.GenKeyPairBytes()
	b := NewInboundMsgBuffer(Validators{{
		PubKeyHex: hex.EncodeToString(pk),
		IP:        "1.1.1.1",
		Port:      1234,
	}})

	count := 255

	var received []byte
	go func() {
		for i := 0; i < count; i++ {
			ms := b.DequeueAll()
			for _, m := range ms {
				received = append(received, m.Msg.GetBytes()...)
			}
		}
	}()

	for i := 1; i <= count; i++ {
		bs := []byte{uint8(i)}
		sig := crypto.SignBytes(sk, bs)

		msgs := []*types.Envelope{{
			Sig:            sig,
			ValidatorIndex: 0,
			Msg: &types.Message{
				Message:  &types.Message_Bytes{Bytes: bs},
				Deadline: time.Now().Add(1 * time.Second).UnixMilli(),
			},
		}}
		go b.Enqueue(msgs)
	}

	slices.Sort(received)

	var expected []byte
	for i := 1; i <= count; i++ {
		expected = append(expected, uint8(i))
	}
	time.Sleep(1 * time.Second)
	require.Equal(t, received, expected)
}
