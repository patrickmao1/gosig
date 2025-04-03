package network

import (
	"context"
	"fmt"
	"github.com/patrickmao1/gosig/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"
)

var nodeURLs = []string{
	"localhost:8081",
	"localhost:8082",
	"localhost:8083",
	"localhost:8084",
	"localhost:8085",
}

var clients []types.NetworkTestClient

func init() {
	for _, url := range nodeURLs {
		clients = append(clients, newClient(url))
	}
}

func newClient(url string) types.NetworkTestClient {
	dialOpt := grpc.WithTransportCredentials(insecure.NewCredentials())
	cc, err := grpc.NewClient(url, dialOpt)
	if err != nil {
		log.Fatal(err)
	}
	return types.NewNetworkTestClient(cc)
}

func TestNetwork(t *testing.T) {
	ctx := context.Background()

	msg := "hello 1"
	_, err := clients[0].Broadcast(ctx, &types.BroadcastReq{Value: msg})
	require.NoError(t, err)

	time.Sleep(time.Second)

	for i, client := range clients {
		res, err := client.GetMsgs(ctx, &types.GetReq{})
		require.NoError(t, err, i)
		require.Len(t, res.Values, 1)
		require.Equal(t, msg, res.Values[0])
	}
}

func TestParallel(t *testing.T) {
	ctx := context.Background()

	var msgs []string
	for i, client := range clients {
		msg := fmt.Sprintf("hello from %d", i)
		_, err := client.Broadcast(ctx, &types.BroadcastReq{Value: msg})
		require.NoError(t, err, i)
		msgs = append(msgs, msg)
	}

	time.Sleep(time.Second)

	for i, client := range clients {
		res, err := client.GetMsgs(ctx, &types.GetReq{})
		require.NoError(t, err, i)
		require.EqualValues(t, msgs, res.Values, i)
	}
}
