package main

import (
	"context"
	"github.com/patrickmao1/gosig/blockchain"
	"github.com/patrickmao1/gosig/types"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"time"
)

func main() {
	service := NewTestService()
	service.Run()
}

type TestService struct {
	types.UnimplementedNetworkTestServer
	network   *blockchain.Network
	msgs      []string
	processed map[[32]byte]bool
}

func NewTestService() *TestService {
	peers := []*blockchain.Validator{
		{IP: "172.16.0.1", Port: 9090},
		{IP: "172.16.0.2", Port: 9090},
		{IP: "172.16.0.3", Port: 9090},
		{IP: "172.16.0.4", Port: 9090},
		{IP: "172.16.0.5", Port: 9090},
	}
	network := blockchain.NewNetwork(2, 300*time.Millisecond, peers)
	return &TestService{
		network:   network,
		processed: make(map[[32]byte]bool),
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

func (s *TestService) handleMessage(msg *types.SignedMessage) {
	b := msg.GetBytes()
	msgId := blockchain.ComputeMsgId(msg)
	if _, ok := s.processed[msgId]; ok {
		return
	}
	s.processed[msgId] = true
	log.Infof("got message %s", b)
	s.network.Broadcast(msg)
	s.msgs = append(s.msgs, string(b))
}

func (s *TestService) GetMsgs(_ context.Context, _ *types.GetReq) (*types.GetRes, error) {
	return &types.GetRes{Values: s.msgs}, nil
}

func (s *TestService) Broadcast(_ context.Context, req *types.BroadcastReq) (*types.BroadcastRes, error) {
	msg := &types.SignedMessage{
		MessageTypes: &types.SignedMessage_Bytes{Bytes: []byte(req.Value)},
		Deadline:     time.Now().Add(time.Second).UnixMilli(),
	}
	s.network.Broadcast(msg)
	return &types.BroadcastRes{}, nil
}
