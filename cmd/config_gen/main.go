package main

import (
	"github.com/patrickmao1/gosig/blockchain"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"strings"
)

var nodeURLs = []string{
	"10.0.0.10:8080",
	"10.0.0.11:8080",
	"10.0.0.12:8080",
	"10.0.0.13:8080",
	"10.0.0.14:8080",
}

func main() {
	cfgs := blockchain.GenTestConfigs(buildValidators(nodeURLs), "/app/runtime/gosig.db")
	bs, err := yaml.Marshal(cfgs)
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile("config.yaml", os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	_, err = f.Write(bs)
	if err != nil {
		panic(err)
	}
}

func buildValidators(nodeURLs []string) blockchain.Validators {
	ret := blockchain.Validators{}
	for _, url := range nodeURLs {
		s := strings.Split(url, ":")
		ip, portStr := s[0], s[1]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			panic(err)
		}
		ret = append(ret, &blockchain.Validator{IP: ip, Port: port})
	}
	return ret
}
