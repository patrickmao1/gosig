package blockchain

import (
	"fmt"
	"github.com/patrickmao1/gosig/types"
	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/protobuf/proto"
)

var tcStateKey = []byte("tc_state")
var headKey = []byte("head_block")

func blockHeaderKey(blockHash []byte) []byte {
	return []byte(fmt.Sprintf("block/%x/header", blockHash))
}

func blockTxsKey(blockHash []byte) []byte {
	return []byte(fmt.Sprintf("block/%x/txs", blockHash))
}

func blockPCertKey(blockHash []byte) []byte {
	return []byte(fmt.Sprintf("block/%x/p_cert", blockHash))
}

func blockTcCertKey(blockHash []byte) []byte {
	return []byte(fmt.Sprintf("block/%x/tc_cert", blockHash))
}

func txKey(txHash []byte) []byte {
	return []byte(fmt.Sprintf("transaction/%x", txHash))
}

type DB struct {
	*leveldb.DB
}

func NewDB(db *leveldb.DB) *DB {
	return &DB{
		DB: db,
	}
}

func (db *DB) GetBlockHeader(blockHash []byte) (*types.BlockHeader, error) {
	bytes, err := db.Get(blockHeaderKey(blockHash), nil)
	if err != nil {
		return nil, err
	}
	header := &types.BlockHeader{}
	err = proto.Unmarshal(bytes, header)
	return header, err
}

func (db *DB) GetHeadBlock() (*types.BlockHeader, error) {
	headHash, err := db.Get(headKey, nil)
	if err != nil {
		return nil, err
	}
	head, err := db.GetBlockHeader(headHash)
	if err != nil {
		return nil, err
	}
	return head, err
}

func (db *DB) GetTcCert(blockHash []byte) (*types.Certificate, error) {
	bytes, err := db.Get(blockTcCertKey(blockHash), nil)
	if err != nil {
		return nil, err
	}
	cert := &types.Certificate{}
	err = proto.Unmarshal(bytes, cert)
	return cert, err
}

func (db *DB) GetTcState() (*types.BlockHeader, uint32, *types.Certificate, error) {
	bs, err := db.Get(tcStateKey, nil)
	if err != nil {
		return nil, 0, nil, err
	}
	tc := &types.TcState{}
	err = proto.Unmarshal(bs, tc)

	blockHash := tc.BlockHash
	h, err := db.GetBlockHeader(blockHash)
	if err != nil {
		return nil, 0, nil, err
	}
	tcCert, err := db.GetTcCert(blockHash)
	if err != nil {
		return nil, 0, nil, err
	}
	return h, tc.Round, tcCert, err
}

func (db *DB) GetBlockTxHashes(blockHash []byte) (*types.TransactionHashes, error) {
	bs, err := db.Get(blockTxsKey(blockHash), nil)
	if err != nil {
		return nil, err
	}
	tc := &types.TransactionHashes{}
	err = proto.Unmarshal(bs, tc)
	return tc, err
}

func (db *DB) GetBlockTxs(blockHash []byte) (txs []*types.SignedTransaction, root []byte, err error) {
	txHashes, err := db.GetBlockTxHashes(blockHash)
	if err != nil {
		return nil, nil, err
	}
	for _, txHash := range txHashes.TxHashes {
		txBytes, err := db.Get(txKey(txHash), nil)
		if err != nil {
			return nil, nil, err
		}
		tx := &types.SignedTransaction{}
		err = proto.Unmarshal(txBytes, tx)
		if err != nil {
			return nil, nil, err
		}
		txs = append(txs, tx)
	}
	return txs, txHashes.Root, nil
}
