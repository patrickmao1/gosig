package rpc

import (
	"context"
	"github.com/patrickmao1/gosig/blockchain"
	"github.com/patrickmao1/gosig/types"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	types.UnsafeRPCServer
	pool *blockchain.TxPool
}

func NewServer(txPool *blockchain.TxPool) *Server {
	return &Server{
		pool: txPool,
	}
}

func (s *Server) Start() {
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

func (s *Server) SubmitTransaction(ctx context.Context, req *types.SubmitTransactionReq) (*types.SubmitTransactionRes, error) {
	tx := req.GetTx().GetTx()
	log.Infof("Got transaction: from %x, to %x, amount %d", tx.From, tx.To, tx.Amount)

	err := s.pool.AppendTx(req.Tx)
	if err != nil {
		return nil, err
	}
	return &types.SubmitTransactionRes{}, nil
}

func (s *Server) TransactionStatus(ctx context.Context, req *types.TransactionStatusReq) (*types.TransactionStatusRes, error) {
	//TODO implement me
	panic("implement me")
}
