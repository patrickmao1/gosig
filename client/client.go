package client

import (
	"context"
	"github.com/patrickmao1/gosig/crypto"
	"github.com/patrickmao1/gosig/types"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"sync"
)

type Client struct {
	clients []types.RPCClient
	privkey []byte
	pubkey  []byte
}

func New(privkey, pubkey []byte, nodes []string) *Client {
	c := new(Client)
	c.privkey = privkey
	c.pubkey = pubkey
	for _, url := range nodes {
		dialOpt := grpc.WithTransportCredentials(insecure.NewCredentials())
		cc, err := grpc.NewClient(url, dialOpt)
		if err != nil {
			log.Fatal(err)
		}
		c.clients = append(c.clients, types.NewRPCClient(cc))
	}
	return c
}

func (c *Client) SubmitTx(tx *types.Transaction) error {
	bs, err := proto.Marshal(tx)
	if err != nil {
		return err
	}
	sig := crypto.SignBytes(c.privkey, bs)
	signedTx := &types.SignedTransaction{
		Tx:  tx,
		Sig: sig,
	}
	pass := crypto.VerifySigBytes(c.pubkey, bs, sig)
	if !pass {
		panic("failed to verify my own signature")
	}
	req := &types.SubmitTransactionReq{Tx: signedTx}

	c.callClients(func(idx int, client types.RPCClient) {
		_, err := client.SubmitTransaction(context.Background(), req)
		if err != nil {
			log.Errorf("failed to submit transaction to node %d: %s", idx, err.Error())
		}
	})
	return nil
}

func (c *Client) callClients(call func(idx int, client types.RPCClient)) {
	var wg = new(sync.WaitGroup)
	for i, client := range c.clients {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			call(i, client)
			wg.Done()
		}(wg)
	}
	wg.Wait()
}
