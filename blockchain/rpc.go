package blockchain

import (
	"context"
	"github.com/patrickmao1/gosig/crypto"
	"github.com/patrickmao1/gosig/types"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
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
	log.Infof("submittx %x.. 11111", txHash[:8])
	err := s.txPool.AddTransaction(req.Tx)
	if err != nil {
		return nil, err
	}
	log.Infof("submittx %x.. 22222", txHash[:8])
	msg := &types.Message{Message: &types.Message_Tx{Tx: req.Tx}}
	s.outMsgs.Put(s.txKey(txHash), msg, s.gossipDDLFromNow())
	log.Infof("submittx %x.. 33333", txHash[:8])
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

func (s *Service) Send(ctx context.Context, req *types.Envelopes) (*types.Empty, error) {
	s.inMsgs.Enqueue(req.Msgs)
	return &types.Empty{}, nil
}

func (s *Service) QueryTXs(ctx context.Context, req *types.QueryTXsReq) (*types.QueryTXsRes, error) {
	txs, _ := s.txPool.BatchGet(req.TxHashes)
	return &types.QueryTXsRes{Txs: txs}, nil
}

func (s *Service) checkSig(signedMsg *types.Envelope) (bool, error) {
	pubKey := s.cfg.Validators[signedMsg.ValidatorIndex].GetPubKey()
	msgBytes, err := proto.Marshal(signedMsg.Msg)
	if err != nil {
		return false, err
	}
	pass := crypto.VerifySigBytes(pubKey, msgBytes, signedMsg.Sig)
	return pass, err
}
