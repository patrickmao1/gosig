package gosig

import (
	"github.com/patrickmao1/gosig/types"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

var clientURLs = []string{
	"localhost:8081",
	"localhost:8082",
	"localhost:8083",
	"localhost:8084",
	"localhost:8085",
}

var clients []types.TransactionServiceClient

func init() {
	for _, url := range clientURLs {
		clients = append(clients, newClient(url))
	}
}

func newClient(url string) types.TransactionServiceClient {
	dialOpt := grpc.WithTransportCredentials(insecure.NewCredentials())
	cc, err := grpc.NewClient(url, dialOpt)
	if err != nil {
		log.Fatal(err)
	}
	return types.NewTransactionServiceClient(cc)
}

func Test(t *testing.T) {

}
