package blockchain

import (
	"context"
	"fmt"
	"github.com/patrickmao1/gosig/types"
	"github.com/patrickmao1/gosig/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/stats"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type OutboundMsgPool interface {
	Pack() *types.Envelope
}

type InboundMsgPool interface {
	Enqueue(msgs *types.Envelope)
}

type Network struct {
	MyIP   net.IP
	MyPort int

	peers    []*Validator
	clients  []types.RPCClient
	fanout   int
	interval time.Duration

	out OutboundMsgPool
	in  InboundMsgPool
	mu  sync.Mutex

	stats *BandwidthStatsHandler
}

func NewNetwork(out OutboundMsgPool, in InboundMsgPool, fanout int, interval time.Duration, vals []*Validator) *Network {
	n := &Network{
		peers:    vals,
		fanout:   fanout,
		interval: interval,
		out:      out,
		in:       in,
		stats:    newBandwidthStatsHandler(),
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

func (n *Network) dialPeers() {
	for i, val := range n.peers {
		dialOpts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(32 * 1024 * 1024)),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(32 * 1024 * 1024)),
			grpc.WithStatsHandler(n.stats),
		}
		cc, err := grpc.NewClient(val.GetURL(), dialOpts...)
		if err != nil {
			log.Fatalf("failed to dial node %d: %s", i, err.Error())
		}
		n.clients = append(n.clients, types.NewRPCClient(cc))
	}
}

func (n *Network) StartGossip() {
	log.Infof("start sending gossips: interval %s, fanout %d, peers %+v", n.interval, n.fanout, n.peers)
	n.dialPeers()

	//go n.stats.Start()

	t := time.Tick(n.interval)
	for {
		ts := <-t

		envelope := n.out.Pack()

		if envelope == nil {
			continue
		}

		// send to randomized peers every time
		targets := pickRandN(n.clients, n.fanout)

		//ctx, cancel := context.WithTimeout(context.Background(), n.interval)
		done := make(chan struct{})
		var wg sync.WaitGroup
		for _, peer := range targets[:n.fanout] {
			wg.Add(1)
			client := peer
			go func() {
				defer wg.Done()
				_, err := client.Send(context.Background(), envelope)
				if err != nil {
					log.Errorf("failed to send gossip to peer: %s", err.Error())
					return
				}
			}()
		}

		go func() {
			wg.Wait()
			done <- struct{}{}
		}()

		// wait for either all Sends are done or gossip round is over
		select {
		case <-done:
			log.Debugf("gossip round done: sent %d msgs, took %s", len(envelope.Msgs.Msgs), time.Since(ts))
			//cancel()
			continue
			//case <-ctx.Done():
			//	go func() {
			//		bs, err := proto.Marshal(toSend)
			//		if err != nil {
			//			log.Errorf("failed to marshal gossip to peer: %s", err.Error())
			//		}
			//		log.Infof("gossip round timed out: start_ts %s, %d msgs, %d bytes", ts, len(msgs), len(bs))
			//	}()
			//	continue
		}
	}
}

func (n *Network) QueryTxs(txHashes [][]byte) ([]*types.SignedTransaction, error) {
	tick := time.NewTicker(n.interval)

	m := make(map[[32]byte]*types.SignedTransaction)
	var mu sync.Mutex
	for i := 0; i < 10; i++ {
		<-tick.C

		targets := pickRandN(n.clients, n.fanout)
		var wg sync.WaitGroup

		for _, target := range targets[:n.fanout] {
			wg.Add(1)
			req := &types.QueryTXsReq{TxHashes: txHashes}

			go func() {
				defer wg.Done()
				res, err := target.QueryTXs(context.Background(), req)
				if err != nil {
					log.Errorf("failed to query txs: %s", err.Error())
					return
				}
				for _, tx := range res.Txs {
					var h [32]byte
					copy(h[:], tx.Hash())
					mu.Lock()
					m[h] = tx
					mu.Unlock()
				}
			}()
		}
		wg.Wait()
		if len(m) == len(txHashes) {
			break
		}
	}

	if len(m) != len(txHashes) {
		return nil, fmt.Errorf("failed to acquire missing txs after 10 attempts")
	}

	ret := make([]*types.SignedTransaction, 0, len(txHashes))
	for _, tx := range m {
		ret = append(ret, tx)
	}
	return ret, nil
}

func pickRandN[T any](input []T, n int) []T {
	if n > len(input) {
		n = len(input)
	}
	indices := rand.Perm(len(input))[:n]
	result := make([]T, n)
	for i, idx := range indices {
		result[i] = input[idx]
	}
	return result
}

type BandwidthStatsHandler struct {
	recording *atomic.Bool

	SentBytes     *atomic.Uint64
	ReceivedBytes *atomic.Uint64
}

func newBandwidthStatsHandler() *BandwidthStatsHandler {
	return &BandwidthStatsHandler{
		recording:     new(atomic.Bool),
		SentBytes:     new(atomic.Uint64),
		ReceivedBytes: new(atomic.Uint64),
	}
}

func (h *BandwidthStatsHandler) Start() {
	tick := time.NewTicker(1 * time.Second)
	for {
		h.record()
		<-tick.C
		h.stop()
	}
}

func (h *BandwidthStatsHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	return ctx
}

func (h *BandwidthStatsHandler) HandleConn(ctx context.Context, connStats stats.ConnStats) {}

func (h *BandwidthStatsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	return ctx
}

func (h *BandwidthStatsHandler) HandleRPC(ctx context.Context, s stats.RPCStats) {
	if !h.recording.Load() {
		return
	}
	switch st := s.(type) {
	case *stats.OutPayload:
		h.SentBytes.Add(uint64(st.WireLength))
	case *stats.InPayload:
		h.ReceivedBytes.Add(uint64(st.WireLength))
	}
}

func (h *BandwidthStatsHandler) record() {
	if h.recording.Load() {
		return
	}
	h.SentBytes.Store(0)
	h.ReceivedBytes.Store(0)
	h.recording.Store(true)
}

func (h *BandwidthStatsHandler) stop() {
	h.recording.Store(false)
	sent, received := h.SentBytes.Load(), h.ReceivedBytes.Load()
	log.Infof("Bandwidth stats: send %d, recv %d", sent, received)
}
