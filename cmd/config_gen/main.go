package main

import (
	"github.com/patrickmao1/gosig/blockchain"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"strings"
)

var nodeURLs = []string{
	"172.16.0.1:8080",
	"172.16.0.2:8080",
	"172.16.0.3:8080",
	"172.16.0.4:8080",
	"172.16.0.5:8080",
}

func main() {
	cfgs := blockchain.GenTestConfigs(buildValidators(nodeURLs))
	bs, err := yaml.Marshal(cfgs)
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile("test_config.yaml", os.O_RDWR|os.O_CREATE, 0600)
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
