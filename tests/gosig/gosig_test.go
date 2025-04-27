package gosig

import (
	"github.com/patrickmao1/gosig/blockchain"
	"github.com/patrickmao1/gosig/client"
	"github.com/patrickmao1/gosig/types"
	"github.com/patrickmao1/gosig/utils"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"
	"testing"
	"time"
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
		time.Sleep(5 * time.Millisecond)
	}

	//pass := false
	//for i := 0; i < 20; i++ {
	//	time.Sleep(1 * time.Second)
	//	balance, err := c.GetBalance(pubkeys[0])
	//	if err != nil {
	//		return
	//	}
	//	log.Infof("balance: %v", balance)
	//	if balance == 1000-123 {
	//		pass = true
	//		break
	//	}
	//}
	//require.True(t, pass)
}

func TestBalance(t *testing.T) {
	c := client.New(privkeys[0], pubkeys[0], rpcURLs)
	balance, err := c.GetBalance(pubkeys[0])
	if err != nil {
		return
	}
	log.Infof("balance: %v", balance)
}

func TestDB(t *testing.T) {
	d, err := leveldb.OpenFile("./tests/gosig/docker_volumes/node1/gosig.db", nil)
	require.NoError(t, err)
	db := blockchain.NewDB(d)

	head, err := db.GetHeadBlock()
	require.NoError(t, err)

	log.Println(head)
}
