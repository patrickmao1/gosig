package main

import (
	"flag"
	"fmt"
	"github.com/patrickmao1/gosig/blockchain"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert/yaml"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
)

var (
	logLevel   string
	configPath string
)

func init() {
	flag.StringVar(&logLevel, "logLevel", log.InfoLevel.String(), "log level")
	flag.StringVar(&configPath, "configPath", "/gosig_config.yaml", "config path")
}

func main() {
	fmt.Println("Hello, World!")
	f, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	config := &blockchain.NodeConfig{}
	err = yaml.Unmarshal(f, config)
	if err != nil {
		log.Fatal(err)
	}
	db, err := leveldb.OpenFile(config.DbPath, nil)
	if err != nil {
		log.Fatal(err)
	}
}
