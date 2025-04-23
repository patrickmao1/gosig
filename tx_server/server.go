package tx_server

import (
	"context"
	"github.com/patrickmao1/gosig/blockchain"
	"github.com/patrickmao1/gosig/types"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type TxServer struct {
	types.UnsafeTransactionServiceServer
	pool *blockchain.TxPool
}

func NewTxServer(txPool *blockchain.TxPool) *TxServer {
	return &TxServer{
		pool: txPool,
	}
}

func (s *TxServer) Start() {
	svr := grpc.NewServer()
	lis, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("grpc listening on 0.0.0.0:8080")
	types.RegisterTransactionServiceServer(svr, s)
	err = svr.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *TxServer) SubmitTransaction(ctx context.Context, req *types.SubmitTransactionReq) (*types.SubmitTransactionRes, error) {
	err := s.pool.AppendTx(req.Tx)
	if err != nil {
		return nil, err
	}
	return &types.SubmitTransactionRes{}, nil
}

func (s *TxServer) TransactionStatus(ctx context.Context, req *types.TransactionStatusReq) (*types.TransactionStatusRes, error) {
	//TODO implement me
	panic("implement me")
}
