package main

import (
	"context"
	"github.com/patrickmao1/gosig/blockchain"
	"github.com/patrickmao1/gosig/crypto"
	"github.com/patrickmao1/gosig/types"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"slices"
	"time"
)

func main() {
	service := NewTestService()
	service.Run()
}

type TestService struct {
	types.UnimplementedNetworkTestServer
	network  *blockchain.Network
	msgs     *blockchain.MsgBuffer
	received map[string]string
}

func NewTestService() *TestService {
	peers := []*blockchain.Validator{
		{IP: "172.16.0.1", Port: 9090},
		{IP: "172.16.0.2", Port: 9090},
		{IP: "172.16.0.3", Port: 9090},
		{IP: "172.16.0.4", Port: 9090},
		{IP: "172.16.0.5", Port: 9090},
	}
	sk, _ := crypto.GenKeyPairBytesFromSeed(1)
	msgBuf := blockchain.NewMsgBuffer(sk, 0)
	network := blockchain.NewNetwork(msgBuf, 2, 300*time.Millisecond, peers)
	return &TestService{
		network:  network,
		msgs:     msgBuf,
		received: make(map[string]string),
	}
}

func (s *TestService) Run() {
	go s.network.StartGossip()
	go s.network.Listen(s.handleMessage)
	svr := grpc.NewServer()
	lis, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("grpc listening on 0.0.0.0:8080")
	types.RegisterNetworkTestServer(svr, s)
	err = svr.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *TestService) handleMessage(signedMsg *types.Envelope) {
	b := signedMsg.Msg.GetBytes()
	log.Infof("got message %s", b)
	s.msgs.Put(string(b), signedMsg.Msg)
	s.received[string(b)] = string(b)
}

func (s *TestService) GetMsgs(_ context.Context, _ *types.GetReq) (*types.GetRes, error) {
	var msgs []string
	for _, msg := range s.received {
		msgs = append(msgs, msg)
	}
	slices.Sort(msgs)
	return &types.GetRes{Values: msgs}, nil
}

func (s *TestService) Broadcast(_ context.Context, req *types.BroadcastReq) (*types.BroadcastRes, error) {
	msg := &types.Message{
		Message:  &types.Message_Bytes{Bytes: []byte(req.Value)},
		Deadline: time.Now().Add(time.Second).UnixMilli(),
	}
	s.msgs.Put(req.Value, msg)
	return &types.BroadcastRes{}, nil
}
