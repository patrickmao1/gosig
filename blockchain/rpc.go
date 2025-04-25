package blockchain

import (
	"context"
	"github.com/patrickmao1/gosig/types"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

func (s *Service) StartRPC() {
	svr := grpc.NewServer()
	lis, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("grpc listening on 0.0.0.0:8080")
	types.RegisterRPCServer(svr, s)
	err = svr.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Service) SubmitTransaction(ctx context.Context, req *types.SubmitTransactionReq) (*types.SubmitTransactionRes, error) {
	tx := req.GetTx().GetTx()
	log.Infof("Got transaction: from %x, to %x, amount %d", tx.From, tx.To, tx.Amount)

	err := s.txPool.AddTransaction(req.Tx)
	if err != nil {
		return nil, err
	}

	msg := &types.Message{Message: &types.Message_Tx{Tx: req.Tx}}
	s.outMsgs.Put(s.txKey(tx.Hash()), msg, s.gossipDDLFromNow())

	return &types.SubmitTransactionRes{}, nil
}

func (s *Service) GetBalance(ctx context.Context, req *types.GetBalanceReq) (*types.GetBalanceRes, error) {
	balance, err := s.db.GetBalance(req.Account)
	if err != nil {
		return nil, err
	}
	return &types.GetBalanceRes{Amount: balance}, nil
}
