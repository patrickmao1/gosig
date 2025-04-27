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

func (s *Service) SubmitTransaction(ctx context.Context, req *types.SubmitTransactionReq) (*types.Empty, error) {
	txHash := req.Tx.Tx.Hash()
	err := s.txPool.AddTransaction(req.Tx)
	if err != nil {
		return nil, err
	}
	msg := &types.Message{Message: &types.Message_Tx{Tx: req.Tx}}
	s.outMsgs.Put(s.txKey(txHash), msg, s.gossipDDLFromNow())
	return &types.Empty{}, nil
}

func (s *Service) SubmitTransactions(ctx context.Context, req *types.SubmitTransactionsReq) (*types.Empty, error) {
	for _, tx := range req.Txs {
		err := s.txPool.AddTransaction(tx)
		if err != nil {
			return nil, err
		}
	}
	for _, tx := range req.Txs {
		msg := &types.Message{Message: &types.Message_Tx{Tx: tx}}
		s.outMsgs.Put(s.txKey(tx.Hash()), msg, s.gossipDDLFromNow())
	}
	return &types.Empty{}, nil
}

func (s *Service) GetBalance(ctx context.Context, req *types.GetBalanceReq) (*types.GetBalanceRes, error) {
	balance, err := s.db.GetBalance(req.Account)
	if err != nil {
		return nil, err
	}
	return &types.GetBalanceRes{Amount: balance}, nil
}

func (s *Service) Send(ctx context.Context, req *types.Envelope) (*types.Empty, error) {
	s.inMsgs.Enqueue(req)
	return &types.Empty{}, nil
}

func (s *Service) QueryTXs(ctx context.Context, req *types.QueryTXsReq) (*types.QueryTXsRes, error) {
	txs, _ := s.txPool.BatchGet(req.TxHashes)
	return &types.QueryTXsRes{Txs: txs}, nil
}
