package main

import (
	"flag"
	"github.com/patrickmao1/gosig/blockchain"
	"github.com/patrickmao1/gosig/crypto"
	"github.com/patrickmao1/gosig/utils"
	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"gopkg.in/yaml.v3"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
)

var (
	logLevel   string
	configPath string
)

func init() {
	flag.StringVar(&logLevel, "logLevel", log.InfoLevel.String(), "log level")
	flag.StringVar(&configPath, "configPath", "/app/gosig_config.yaml", "config path")
}

func main() {
	flag.Parse()

	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)

	// parse configs
	cfg := getMyConfig(configPath)

	// init db
	err = os.RemoveAll(cfg.DbPath) // start from a clean db
	if err != nil {
		log.Fatal(err)
	}
	db, err := leveldb.OpenFile(cfg.DbPath, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// dependency modules
	d := blockchain.NewDB(db)
	genesis := blockchain.DefaultGenesisConfig()
	txPool := blockchain.NewTxPool()

	// init test state
	testAmount := 10_000_000
	testAccount := crypto.MarshalECDSAPublic(&utils.TestECDSAKeys[0].PublicKey)
	log.Infof("funding test account %x with %d", testAccount, testAmount)
	err = d.PutBalance(testAccount, 10_000_000)
	if err != nil {
		return
	}

	go func() {
		err := http.ListenAndServe("0.0.0.0:6060", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// starts the blockchain
	chain := blockchain.NewService(cfg, genesis, d, txPool)
	chain.Start()
}

func getMyConfig(path string) *blockchain.NodeConfig {
	f, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	cfg := &blockchain.NodeConfigs{}
	err = yaml.Unmarshal(f, cfg)
	if err != nil {
		log.Fatal(err)
	}
	nodeIndex, ok := os.LookupEnv("node_index")
	if !ok {
		log.Fatal("env variable NODE_INDEX not set")
	}
	ni, err := strconv.Atoi(nodeIndex)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("my node index: %d, my privkey %s", ni, cfg.Configs[ni].PrivKeyHex)
	return cfg.Configs[ni]
}
