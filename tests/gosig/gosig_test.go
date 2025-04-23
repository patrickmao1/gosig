package gosig

import (
	"github.com/patrickmao1/gosig/blockchain"
	"github.com/patrickmao1/gosig/client"
	"github.com/patrickmao1/gosig/crypto"
	"github.com/patrickmao1/gosig/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"
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
	for i := 0; i < 10; i++ {
		sk, pk := crypto.GenKeyPairBytesFromSeed(int64(i))
		pubkeys = append(pubkeys, pk)
		privkeys = append(privkeys, sk)
	}
}

func Test(t *testing.T) {
	c := client.New(privkeys[0], pubkeys[0], rpcURLs)
	err := c.SubmitTx(&types.Transaction{
		From:   pubkeys[0],
		To:     pubkeys[1],
		Amount: 123,
	})
	require.NoError(t, err)
}

func TestDB(t *testing.T) {
	d, err := leveldb.OpenFile("./tests/gosig/docker_volumes/node1/gosig.db", nil)
	require.NoError(t, err)
	db := blockchain.NewDB(d)

	head, err := db.GetHeadBlock()
	require.NoError(t, err)

	log.Println(head)
}
