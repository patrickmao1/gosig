package gosig

import (
	"github.com/patrickmao1/gosig/blockchain"
	"github.com/patrickmao1/gosig/client"
	"github.com/patrickmao1/gosig/types"
	"github.com/patrickmao1/gosig/utils"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"
	"sync"
	"testing"
)

var rpcURLs = []string{
	"localhost:8081",
	"localhost:8082",
	"localhost:8083",
	"localhost:8084",
	"localhost:8085",
}
var pubkeys [][]byte
var privkeys [][]byte

func init() {
	pubkeys, privkeys = utils.GenTestKeyPairs(10)
}

func Test(t *testing.T) {
	c := client.New(privkeys[0], pubkeys[0], rpcURLs)
	for i := 0; i < 1000; i++ {
		log.Infof("sending tx %d", i)
		err := c.SubmitTx(&types.Transaction{
			From:   pubkeys[0],
			To:     pubkeys[1],
			Amount: 1,
			Nonce:  uint32(i),
		})
		require.NoError(t, err)
	}
}

func TestParallel(t *testing.T) {
	c := client.New(privkeys[0], pubkeys[0], rpcURLs)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var txs []*types.Transaction
			for j := 0; j < 500; j++ {
				txs = append(txs, &types.Transaction{
					From:   pubkeys[0],
					To:     pubkeys[1],
					Amount: 1,
					Nonce:  uint32(i*500 + j),
				})
			}
			err := c.SubmitTxs(txs)
			require.NoError(t, err)
		}()
	}
	wg.Wait()
}

func TestBatch(t *testing.T) {
	c := client.New(privkeys[0], pubkeys[0], rpcURLs)
	var txs []*types.Transaction

	for i := 0; i < 10_000; i++ {
		txs = append(txs, &types.Transaction{
			From:   pubkeys[0],
			To:     pubkeys[1],
			Amount: 1,
			Nonce:  uint32(i),
		})
	}
	err := c.SubmitTxs(txs)
	require.NoError(t, err)
}

func TestBalance(t *testing.T) {
	c := client.New(privkeys[0], pubkeys[0], rpcURLs)
	balance0, err := c.GetBalance(pubkeys[0])
	if err != nil {
		return
	}
	balance1, err := c.GetBalance(pubkeys[1])
	if err != nil {
		return
	}
	log.Infof("wallet 0 balance: %d", balance0)
	log.Infof("wallet 1 balance: %d", balance1)
}

func TestDB(t *testing.T) {
	d, err := leveldb.OpenFile("./tests/gosig/docker_volumes/node1/gosig.db", nil)
	require.NoError(t, err)
	db := blockchain.NewDB(d)

	head, err := db.GetHeadBlock()
	require.NoError(t, err)

	log.Println(head)
}
