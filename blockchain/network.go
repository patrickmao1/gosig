package blockchain

import (
	"bufio"
	"context"
	"fmt"
	"github.com/patrickmao1/gosig/types"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

type MessageBuffer interface {
	List() []*types.SignedMessage
}

type Network struct {
	MyIP   net.IP
	MyPort int

	peers    []*Validator
	fanout   int
	interval time.Duration

	msgs MessageBuffer
	mu   sync.Mutex
}

func NewNetwork(msgPool MessageBuffer, fanout int, interval time.Duration, vals []*Validator) *Network {
	n := &Network{
		peers:    vals,
		fanout:   fanout,
		interval: interval,
		msgs:     msgPool,
	}
	var err error
	n.MyIP, err = getPrivateIP()
	if err != nil {
		log.Fatal("failed to get private ip", err)
	}
	// find myself in the vals list and assign my port
	for _, val := range vals {
		if val.IP == n.MyIP.String() {
			n.MyPort = val.Port
			break
		}
	}
	return n
}

func (n *Network) StartGossip() {
	log.Infof("start sending gossips: interval %s, fanout %d, peers %v", n.interval, n.fanout, n.peers)
	t := time.Tick(n.interval)
	for {
		ts := <-t
		log.Debugf("start new gossip round: ts %s", ts)

		msgs := n.msgs.List()

		if len(msgs) == 0 {
			log.Debugf("no messages in buffer")
			continue
		}

		// send to randomized peers every time
		rand.Shuffle(len(n.peers), func(i, j int) {
			n.peers[i], n.peers[j] = n.peers[j], n.peers[i]
		})

		ctx, cancel := context.WithTimeout(context.Background(), n.interval)
		done := make(chan struct{})
		var wg sync.WaitGroup

		for _, peer := range n.peers[:n.fanout] {
			wg.Add(1)
			url := peer.GetURL()
			go func() {
				defer wg.Done()
				conn, err := n.open(url)
				if err != nil {
					log.Errorln(err)
					return
				}
				bs, err := proto.Marshal(&types.SignedMessages{Msgs: msgs})
				if err != nil {
					log.Errorln(err)
					return
				}
				err = n.send(conn, bs)
				if err != nil {
					log.Errorln(err)
				}
			}()
		}

		go func() {
			wg.Wait()
			done <- struct{}{}
		}()

		select {
		case <-done:
			log.Debugf("gossip round done: start_ts %s", ts)
			cancel()
			continue
		case <-ctx.Done():
			log.Infof("gossip round timed out: start_ts %s", ts)
			continue
		}
	}
}

func (n *Network) open(url string) (*net.UDPConn, error) {
	receiver, err := net.ResolveUDPAddr("udp", url)
	if err != nil {
		return nil, fmt.Errorf("error resolving reciever address: %s", err.Error())
	}
	conn, err := net.DialUDP("udp", nil, receiver)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %s", err.Error())
	}
	return conn, nil
}

func (n *Network) send(conn *net.UDPConn, b []byte) error {
	b = append(b, 0x00)
	w := bufio.NewWriter(conn)

	_, err := w.Write(b)

	if err != nil {
		return fmt.Errorf("failed to send msg while writing to buffer: %s", err.Error())
	}
	err = w.Flush()
	if err != nil {
		return fmt.Errorf("failed to send msg while flushing buffer: %s", err.Error())
	}
	return nil
}

type MessageHandler func(msg *types.SignedMessage)

func (n *Network) Listen(handleMsg MessageHandler) error {
	address, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%d", n.MyPort))
	if err != nil {
		log.Fatal("Error resolving address:", err)
	}
	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		log.Errorln("error listening:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	log.Infof("network: listening on 0.0.0.0:%d", n.MyPort)

	r := bufio.NewReader(conn)

	for {
		b, err := r.ReadBytes(0x00)
		if err != nil {
			log.Infoln("error reading message:", err)
			continue
		}
		b = b[:len(b)-1] // remove delimiter
		ms := &types.SignedMessages{}
		err = proto.Unmarshal(b, ms)
		if err != nil {
			log.Errorf("error unmarshalling incoming msg: %s", err.Error())
			continue
		}
		for _, m := range ms.Msgs {
			go handleMsg(m)
		}
	}
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
