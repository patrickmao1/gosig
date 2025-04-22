package main

import (
	"context"
	"encoding/hex"
	"fmt"
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
	network *blockchain.Network
	out     *blockchain.OutboundMsgBuffer
	in      *blockchain.InboundMsgBuffer
}

func NewTestService() *TestService {
	peers := blockchain.Validators{
		{IP: "172.16.0.1", Port: 9090},
		{IP: "172.16.0.2", Port: 9090},
		{IP: "172.16.0.3", Port: 9090},
		{IP: "172.16.0.4", Port: 9090},
		{IP: "172.16.0.5", Port: 9090},
	}
	s := &TestService{}
	var mySk []byte
	for _, peer := range peers {
		sk, pk := crypto.GenKeyPairBytes()
		peer.PubKeyHex = hex.EncodeToString(pk)
		myip, err := getPrivateIP()
		if err != nil {
			log.Fatal(err)
		}
		if peer.IP == myip.String() {
			mySk = sk
		}
	}
	if mySk == nil {
		log.Fatal("i dont have a key!")
	}
	s.out = blockchain.NewOutboundMsgBuffer(mySk, 0)
	s.in = blockchain.NewInboundMsgBuffer(peers)
	s.network = blockchain.NewNetwork(s.out, s.in, 2, 100*time.Millisecond, peers)
	return s
}

func (s *TestService) Run() {
	go s.network.StartGossip()

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

func (s *TestService) GetMsgs(_ context.Context, _ *types.GetReq) (*types.GetRes, error) {
	msgs := s.in.DequeueAll()

	var ss = make(map[string]bool)
	for _, msg := range msgs {
		ss[string(msg.Msg.GetBytes())] = true
	}
	sss := make([]string, 0, len(ss))
	for m := range ss {
		sss = append(sss, m)
	}

	slices.Sort(sss)

	log.Infoln("msgs", ss)

	return &types.GetRes{Values: sss}, nil
}

func (s *TestService) Broadcast(_ context.Context, req *types.BroadcastReq) (*types.BroadcastRes, error) {
	msg := &types.Message{
		Message:  &types.Message_Bytes{Bytes: []byte(req.Value)},
		Deadline: time.Now().Add(2 * time.Second).UnixMilli(),
	}
	s.out.Put(req.Value, msg)
	return &types.BroadcastRes{}, nil
}

func getPrivateIP() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return net.IP{}, err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil && isPrivateIP(ipNet.IP) {
				return ipNet.IP, nil
			}
		}
	}
	return net.IP{}, fmt.Errorf("no private IP found. all addrs:%s", addrs)
}

func isPrivateIP(ip net.IP) bool {
	privateBlocks := []net.IPNet{
		{IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(8, 32)},
		{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)},
		{IP: net.IPv4(192, 168, 0, 0), Mask: net.CIDRMask(16, 32)},
	}
	for _, block := range privateBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}
