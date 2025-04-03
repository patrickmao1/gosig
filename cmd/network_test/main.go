package main

import (
	"context"
	"github.com/patrickmao1/gosig/blockchain"
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
	network   *blockchain.Network
	msgs      *blockchain.MsgBuffer
	processed []string
}

func NewTestService() *TestService {
	peers := []*blockchain.Validator{
		{IP: "172.16.0.1", Port: 9090},
		{IP: "172.16.0.2", Port: 9090},
		{IP: "172.16.0.3", Port: 9090},
		{IP: "172.16.0.4", Port: 9090},
		{IP: "172.16.0.5", Port: 9090},
	}
	msgBuf := blockchain.NewMsgBuffer()
	network := blockchain.NewNetwork(msgBuf, 2, 300*time.Millisecond, peers)
	return &TestService{
		network: network,
		msgs:    msgBuf,
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
	if slices.Contains(s.processed, string(b)) {
		return
	}
	log.Infof("got message %s", b)
	s.msgs.Put(string(b), msg)
	s.processed = append(s.processed, string(msg.GetBytes()))
}

func (s *TestService) GetMsgs(_ context.Context, _ *types.GetReq) (*types.GetRes, error) {
	slices.Sort(s.processed)
	return &types.GetRes{Values: s.processed}, nil
}

func (s *TestService) Broadcast(_ context.Context, req *types.BroadcastReq) (*types.BroadcastRes, error) {
	msg := &types.SignedMessage{
		MessageTypes: &types.SignedMessage_Bytes{Bytes: []byte(req.Value)},
		Deadline:     time.Now().Add(time.Second).UnixMilli(),
	}
	s.msgs.Put(string(msg.GetBytes()), msg)
	return &types.BroadcastRes{}, nil
}
