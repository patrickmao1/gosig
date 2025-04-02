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

type msgWithDeadline struct {
	msg *types.SignedMessage
	ddl time.Time
}

type Network struct {
	myIP     net.IP
	myPort   int
	peers    []*Validator
	fanout   int
	interval time.Duration
	msgs     []*msgWithDeadline
	mu       sync.Mutex
}

func NewNetwork(fanout int, interval time.Duration, vals []*Validator) *Network {
	n := &Network{
		fanout:   fanout,
		interval: interval,
	}
	var err error
	n.myIP, err = getPrivateIP()
	if err != nil {
		log.Fatal("failed to get private ip", err)
	}
	for _, val := range vals {
		if val.IP == n.myIP.String() {
			n.myPort = val.Port
		}
		n.peers = append(n.peers, val)
	}
	return n
}

func (n *Network) Broadcast(msg *types.SignedMessage, ddl time.Time) {
	n.mu.Lock()
	n.msgs = append(n.msgs, &msgWithDeadline{
		msg: msg,
		ddl: ddl,
	})
	n.mu.Unlock()
}

func (n *Network) StartGossip() {
	log.Infof("start sending gossips: interval %s, fanout %d, peers %v", n.interval, n.fanout, n.peers)
	t := time.Tick(n.interval)
	for {
		ts := <-t
		log.Infof("start new gossip round: ts %s", ts)
		n.mu.Lock()
		newMsgs := make([]*msgWithDeadline, 0)

		// drop deadline-exceeded messages
		for _, msg := range n.msgs {
			if time.Now().Before(msg.ddl) {
				newMsgs = append(newMsgs, msg)
			}
		}
		n.msgs = newMsgs
		n.mu.Unlock()

		if len(newMsgs) == 0 {
			log.Debugf("no messages in buffer")
			continue
		}

		for _, msg := range n.msgs {
			log.Infof("sending messages %v", msg)
		}

		// serialize msgs
		msgs := &types.SignedMessages{}
		for _, msg := range newMsgs {
			msgs.Msgs = append(msgs.Msgs, msg.msg)
		}
		msgsBytes, err := proto.Marshal(msgs)
		if err != nil {
			log.Fatal("unserializable message found in broadcast queue")
		}

		// send to randomized peers every time
		rand.Shuffle(len(n.peers), func(i, j int) {
			n.peers[i], n.peers[j] = n.peers[j], n.peers[i]
		})

		ctx, cancel := context.WithTimeout(context.Background(), n.interval)
		done := make(chan struct{})
		var wg sync.WaitGroup

		for _, val := range n.peers[:n.fanout] {
			wg.Add(1)
			url := val.GetURL()
			go func() {
				defer wg.Done()
				conn, err := n.open(url)
				if err != nil {
					log.Errorln(err)
					return
				}
				err = n.send(conn, msgsBytes)
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
			log.Infof("gossip round done: start_ts %s", ts)
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
	address, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%d", n.myPort))
	if err != nil {
		log.Fatal("Error resolving address:", err)
	}
	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		log.Errorln("error listening:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	log.Infof("network: listening on 0.0.0.0:%d", n.myPort)

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
