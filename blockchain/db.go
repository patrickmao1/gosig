package blockchain

import (
	"fmt"
	"github.com/patrickmao1/gosig/types"
	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/protobuf/proto"
)

var tcBlockKey = []byte("tc_block")
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

func (db *DB) PutTcCert(blockHash []byte, cert *types.Certificate) error {
	bs, err := proto.Marshal(cert)
	if err != nil {
		return err
	}
	return db.Put(blockTcCertKey(blockHash), bs, nil)
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

func (db *DB) GetTcBlock() (*types.BlockHeader, *types.Certificate, error) {
	blockHash, err := db.Get(tcBlockKey, nil)
	if err != nil {
		return nil, nil, err
	}

	h, err := db.GetBlockHeader(blockHash)
	if err != nil {
		return nil, nil, err
	}
	tcCert, err := db.GetTcCert(blockHash)
	if err != nil {
		return nil, nil, err
	}
	return h, tcCert, err
}

func (db *DB) PutTcBlock(blockHash []byte, certificate *types.Certificate) error {
	bs, err := proto.Marshal(certificate)
	if err != nil {
		return err
	}
	err = db.Put(blockTcCertKey(blockHash), bs, nil)
	if err != nil {
		return err
	}
	return db.Put(tcBlockKey, blockHash, nil)
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
