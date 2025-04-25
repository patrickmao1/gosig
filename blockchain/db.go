package blockchain

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/patrickmao1/gosig/types"
	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/protobuf/proto"
)

var headKey = []byte("head_block")

func blockHeaderKey(blockHash []byte) []byte {
	return []byte(fmt.Sprintf("block/%x/header", blockHash))
}

func txHashesKey(blockHash []byte) []byte {
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

func accountKey(pubkey []byte) []byte {
	return []byte(fmt.Sprintf("account/%x", pubkey))
}

type DB struct {
	*leveldb.DB
}

func NewDB(db *leveldb.DB) *DB {
	return &DB{DB: db}
}

func (db *DB) PutBlock(blockHeader *types.BlockHeader) error {
	return db.put(blockHeaderKey(blockHeader.Hash()), blockHeader)
}

func (db *DB) GetBlock(blockHash []byte) (*types.BlockHeader, error) {
	return get[types.BlockHeader](db, blockHeaderKey(blockHash))
}

func (db *DB) PutHeadBlock(blockHash []byte) error {
	return db.Put(headKey, blockHash, nil)
}

func (db *DB) GetHeadBlock() (*types.BlockHeader, error) {
	headHash, err := db.Get(headKey, nil)
	if err != nil {
		return nil, err
	}
	head, err := db.GetBlock(headHash)
	if err != nil {
		return nil, err
	}
	return head, err
}

func (db *DB) PutPCert(blockHash []byte, cert *types.Certificate) error {
	return db.put(blockPCertKey(blockHash), cert)
}

func (db *DB) GetPCert(blockHash []byte) (*types.Certificate, error) {
	return get[types.Certificate](db, blockPCertKey(blockHash))
}

func (db *DB) PutTcCert(blockHash []byte, cert *types.Certificate) error {
	return db.put(blockTcCertKey(blockHash), cert)
}

func (db *DB) GetTcCert(blockHash []byte) (*types.Certificate, error) {
	return get[types.Certificate](db, blockTcCertKey(blockHash))
}

func (db *DB) PutTxHashes(blockHash []byte, txHashes *types.TransactionHashes) error {
	return db.put(txHashesKey(blockHash), txHashes)
}

func (db *DB) GetTxHashes(blockHash []byte) (*types.TransactionHashes, error) {
	return get[types.TransactionHashes](db, txHashesKey(blockHash))
}

func (db *DB) PutBalance(account []byte, balance uint64) error {
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, balance)
	return db.Put(accountKey(account), bs, nil)
}

func (db *DB) GetBalance(account []byte) (uint64, error) {
	bs, err := db.Get(accountKey(account), nil)
	if err != nil && errors.Is(err, leveldb.ErrNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(bs), nil
}

func (db *DB) GetBlockTxs(blockHash []byte) (txs []*types.SignedTransaction, root []byte, err error) {
	txHashes, err := db.GetTxHashes(blockHash)
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
	return txs, txHashes.Root(), nil
}

func (db *DB) put(key []byte, value proto.Message) error {
	bs, err := proto.Marshal(value)
	if err != nil {
		return err
	}
	return db.Put(key, bs, nil)
}

type ProtoMessagePtr[T any] interface {
	proto.Message
	*T
}

func get[T any, PT ProtoMessagePtr[T]](db *DB, key []byte) (*T, error) {
	t := new(T)
	pt := PT(t)
	bytes, err := db.Get(key, nil)
	if err != nil {
		return t, err
	}
	err = proto.Unmarshal(bytes, pt)
	return t, err
}
