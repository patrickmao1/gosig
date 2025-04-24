package blockchain

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/patrickmao1/gosig/crypto"
	"github.com/patrickmao1/gosig/utils"
	log "github.com/sirupsen/logrus"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"
)

type Validator struct {
	PubKeyHex string `yaml:"pub_key"`
	IP        string `yaml:"ip"`
	Port      int    `yaml:"port"`

	// cached values
	pubKey []byte
}

func (v *Validator) GetPubKey() []byte {
	if v.pubKey != nil {
		return v.pubKey
	}
	pubKeyBytes, err := hex.DecodeString(v.PubKeyHex)
	if err != nil {
		panic(err)
	}
	return pubKeyBytes
}

func (v *Validator) GetURL() string {
	return fmt.Sprintf("%s:%d", v.IP, v.Port)
}

type Validators []*Validator

func (vs Validators) PubKeys() [][]byte {
	var ret [][]byte
	for _, val := range vs {
		ret = append(ret, val.GetPubKey())
	}
	return ret
}

type NodeConfig struct {
	DbPath                   string     `yaml:"db_path"`
	PrivKeyHex               string     `yaml:"priv_key_hex"`
	ProposalThreshold        uint32     `yaml:"proposal_threshold"`
	GossipIntervalMs         int        `yaml:"gossip_interval"`
	GossipDegree             int        `yaml:"gossip_degree"`
	ProposalStageDurationMs  int64      `yaml:"proposal_stage_duration_ms"`
	AgreementStateDurationMs int64      `yaml:"agreement_state_duration_ms"`
	Validators               Validators `yaml:"validators"`

	// cached values
	privKey []byte
	myIndex *uint32
}

func (c *NodeConfig) MyPrivKey() []byte {
	if c.privKey != nil {
		return c.privKey
	}
	var err error
	c.privKey, err = hex.DecodeString(c.PrivKeyHex)
	if err != nil {
		panic(err)
	}
	return c.privKey
}

func (c *NodeConfig) Me() *Validator {
	if c.myIndex != nil {
		return c.Validators[*c.myIndex]
	}
	myIP, err := utils.GetPrivateIP()
	if err != nil {
		log.Fatal("failed to get private ip", err)
	}
	for i, val := range c.Validators {
		if val.IP == myIP.String() {
			j := uint32(i)
			c.myIndex = &j
			return val
		}
	}
	log.Fatalf("my IP %s is not in the validator set", myIP)
	return nil
}

func (c *NodeConfig) Quorum() int {
	return len(c.Validators)*2/3 + 1
}

func (c *NodeConfig) MyValidatorIndex() uint32 {
	c.Me()
	return *c.myIndex
}

func (c *NodeConfig) ProposalStageDuration() time.Duration {
	return time.Duration(c.ProposalStageDurationMs) * time.Millisecond
}

func (c *NodeConfig) AgreementStateDuration() time.Duration {
	return time.Duration(c.AgreementStateDurationMs) * time.Millisecond
}

func (c *NodeConfig) RoundDuration() time.Duration {
	return c.ProposalStageDuration() + c.AgreementStateDuration()
}

func (c *NodeConfig) ProposalThresholdPerc() float64 {
	t := big.NewFloat(float64(c.ProposalThreshold))
	t.Quo(t, big.NewFloat(float64(math.MaxUint32)))
	ret, _ := t.Float64()
	return ret
}

type NodeConfigs struct {
	Configs []*NodeConfig `yaml:"configs"`
}

func GenTestConfigs(nodes Validators) (cfgs *NodeConfigs) {
	cfg := NodeConfig{
		DbPath:                   "/app/runtime/gosig.db",
		ProposalThreshold:        computeThreshold(len(nodes)),
		GossipIntervalMs:         500,
		GossipDegree:             2,
		ProposalStageDurationMs:  4000,
		AgreementStateDurationMs: 8000,
	}
	privkeys := make([]string, 0, len(nodes))
	for _, node := range nodes {
		bytes := strings.Split(node.IP, ".")
		ipBytes := make([]byte, len(bytes))
		for i, bs := range bytes {
			n, err := strconv.Atoi(bs)
			if err != nil {
				panic(err)
			}
			if n > 255 {
				panic("too big")
			}
			ipBytes[i] = byte(n)
		}

		bs := make([]byte, 8)
		copy(bs[4:8], ipBytes)
		seed := binary.BigEndian.Uint64(bs)

		privkey, pubkey := crypto.GenKeyPairBytesFromSeed(int64(seed))
		node.PubKeyHex = hex.EncodeToString(pubkey)
		privkeys = append(privkeys, hex.EncodeToString(privkey))
	}
	cfgs = &NodeConfigs{}
	for i := range nodes {
		cp := cfg
		cp.Validators = nodes
		cp.PrivKeyHex = privkeys[i]
		cfgs.Configs = append(cfgs.Configs, &cp)
	}
	return cfgs
}

func computeThreshold(n int) uint32 {
	// T = f(N) such that the probability of no one proposing a block is 0.01
	// T = 1 - e^(-4.60517/N),
	const constant = 4.60517
	t := 1.0 - math.Exp(-constant/float64(n))
	return uint32(float64(math.MaxUint32) * t)
}
