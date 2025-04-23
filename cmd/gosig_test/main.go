package main

import (
	"flag"
	"github.com/patrickmao1/gosig/blockchain"
	"github.com/patrickmao1/gosig/rpc"
	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"gopkg.in/yaml.v3"
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
	for i, val := range cfg.Validators {
		log.Infof("Validator #%d: pubkey %x..", i, val.GetPubKey()[:8])
	}

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

	// service modules
	txServer := rpc.NewServer(txPool)
	chain := blockchain.NewService(cfg, genesis, d, txPool)

	// start services
	go chain.Start()
	txServer.Start()
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
