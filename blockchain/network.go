package blockchain

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/patrickmao1/gosig/types"
	"github.com/patrickmao1/gosig/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"io"
	"math/rand"
	"net"
	"os"
	"slices"
	"sync"
	"time"
)

type OutboundMsgPool interface {
	List() []*types.Envelope
}

type InboundMsgPool interface {
	Enqueue(msgs []*types.Envelope)
}

type Network struct {
	MyIP   net.IP
	MyPort int

	peers    []*Validator
	fanout   int
	interval time.Duration

	out OutboundMsgPool
	in  InboundMsgPool
	mu  sync.Mutex
}

func NewNetwork(out OutboundMsgPool, in InboundMsgPool, fanout int, interval time.Duration, vals []*Validator) *Network {
	n := &Network{
		peers:    vals,
		fanout:   fanout,
		interval: interval,
		out:      out,
		in:       in,
	}
	var err error
	n.MyIP, err = utils.GetPrivateIP()
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
	log.Infof("start sending gossips: interval %s, fanout %d, peers %+v", n.interval, n.fanout, n.peers)
	go n.listen()
	t := time.Tick(n.interval)
	for {
		ts := <-t

		msgs := n.out.List()

		if len(msgs) == 0 {
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
				bs, err := proto.Marshal(&types.Envelopes{Msgs: msgs})
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
			log.Debugf("gossip round done: sent %d msgs, start_ts %s", len(msgs), ts)
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
	lb := make([]byte, 8)
	binary.BigEndian.PutUint64(lb, uint64(len(b)))
	b = slices.Concat(lb, b)
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

func (n *Network) listen() {
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
		b := make([]byte, 8)
		read, err := io.ReadFull(r, b)
		if err != nil {
			log.Infoln("error reading message len prefix:", err)
			continue
		}
		if read != 8 {
			log.Errorf("error reading message: expected read 8, actual read %d", read)
			continue
		}
		l := binary.BigEndian.Uint64(b)
		bs := make([]byte, l)
		read, err = io.ReadFull(r, bs)
		if err != nil {
			log.Errorln("error reading message:", err)
			continue
		}
		if uint64(read) != l {
			log.Errorf("error reading message: expected read %d, actual read %d", l, read)
			continue
		}

		ms := &types.Envelopes{}
		err = proto.Unmarshal(bs, ms)
		if err != nil {
			log.Errorf("error unmarshalling incoming msg: %s, read bytes %x", err.Error(), bs)
			continue
		}

		n.in.Enqueue(ms.Msgs)
	}
}
