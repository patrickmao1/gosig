package types

import (
	"golang.org/x/crypto/sha3"
	"google.golang.org/protobuf/proto"
)

func (tx *Transaction) Hash() []byte {
	bs, err := proto.Marshal(tx)
	if err != nil {
		panic(err)
	}
	hash := sha3.Sum256(bs)
	return hash[:]
}

func (h *BlockHeader) Hash() []byte {
	bs, err := proto.Marshal(h)
	if err != nil {
		panic(err)
	}
	hash := sha3.Sum256(bs)
	return hash[:]
}
