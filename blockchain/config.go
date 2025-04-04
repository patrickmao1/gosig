package blockchain

import (
	"encoding/hex"
	"fmt"
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

type NodeConfig struct {
	DbPath                   string       `yaml:"db_path"`
	PrivKeyHex               string       `yaml:"priv_key_hex"`
	ProposalThreshold        uint32       `yaml:"proposal_threshold"`
	GossipIntervalMs         int          `yaml:"gossip_interval"`
	GossipDegree             int          `yaml:"gossip_degree"`
	ProposalStageDurationMs  int64        `yaml:"proposal_stage_duration_ms"`
	AgreementStateDurationMs int64        `yaml:"agreement_state_duration_ms"`
	Validators               []*Validator `yaml:"validators"`

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
	myIP, err := getPrivateIP()
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
	log.Fatal("my IP %s is not in the validator set", myIP)
	return nil
}

func (c *NodeConfig) MyValidatorIndex() uint32 {
	c.Me()
	return *c.myIndex
}

func (c *NodeConfig) ProposalStageDuration() time.Duration {
	return time.Duration(c.ProposalStageDurationMs) * time.Millisecond
}
