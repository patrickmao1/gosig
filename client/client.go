package client

import (
	"context"
	"fmt"
	"github.com/patrickmao1/gosig/crypto"
	"github.com/patrickmao1/gosig/types"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"math/rand"
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
		dialOpt := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(32 * 1024 * 1024)),
		}
		cc, err := grpc.NewClient(url, dialOpt...)
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
	req := &types.SubmitTransactionReq{Tx: signedTx}

	i := rand.Intn(len(c.clients))
	log.Infof("sending to node %d", i)
	_, err = c.clients[i].SubmitTransaction(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to submit transaction to node %d: %s", i, err.Error())
	}
	return nil
}

func (c *Client) SubmitTxs(txs []*types.Transaction) error {
	signedTxs := make([]*types.SignedTransaction, len(txs))
	for i, tx := range txs {
		if (i+1)%10000 == 0 {
			log.Infof("progress %d", i)
		}
		bs, err := proto.Marshal(tx)
		if err != nil {
			return err
		}
		sig := crypto.SignBytes(c.privkey, bs)
		signedTx := &types.SignedTransaction{
			Tx:  tx,
			Sig: sig,
		}
		signedTxs[i] = signedTx
	}
	req := &types.SubmitTransactionsReq{Txs: signedTxs}

	i := rand.Intn(len(c.clients))
	log.Infof("sending %d txs to node %d", len(txs), i)
	_, err := c.clients[i].SubmitTransactions(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to submit transaction to node %d: %s", i, err.Error())
	}
	return nil
}

func (c *Client) GetBalance(pubkey []byte) (uint64, error) {
	req := &types.GetBalanceReq{Account: pubkey}
	i := rand.Intn(len(c.clients))
	res, err := c.clients[i].GetBalance(context.Background(), req)
	if err != nil {
		return 0, fmt.Errorf("failed to get balance from node %d: %s", i, err.Error())
	}
	return res.Amount, nil
}
