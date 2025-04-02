package blockchain

import (
	"encoding/hex"
	"fmt"
	bls12381 "github.com/kilic/bls12-381"
)

type Validator struct {
	PubKeyHex string `yaml:"pub_key"`
	IP        string `yaml:"ip"`
	Port      int    `yaml:"port"`

	// cached values
	pubKey *bls12381.PointG1
}

func (v *Validator) GetPubKey() *bls12381.PointG1 {
	if v.pubKey != nil {
		return v.pubKey
	}
	pubKeyBytes, err := hex.DecodeString(v.PubKeyHex)
	if err != nil {
		panic(err)
	}
	pubKey, err := bls12381.NewG1().FromCompressed(pubKeyBytes)
	if err != nil {
		panic(err)
	}
	v.pubKey = pubKey
	return pubKey
}

func (v *Validator) GetURL() string {
	return fmt.Sprintf("%s:%d", v.IP, v.Port)
}

type NodeConfig struct {
	DbPath                   string       `yaml:"db_path"`
	PrivKeyHex               string       `yaml:"priv_key_hex"`
	ProposalThreshold        uint32       `yaml:"proposal_threshold"`
	GossipInterval           int          `yaml:"gossip_interval"`
	GossipDegree             int          `yaml:"gossip_degree"`
	ProposalStageDurationMs  int64        `yaml:"proposal_stage_duration_ms"`
	AgreementStateDurationMs int64        `yaml:"agreement_state_duration_ms"`
	Validators               []*Validator `yaml:"validators"`

	// cached values
	privKey *bls12381.Fr
}

func (c *NodeConfig) GetPrivKey() *bls12381.Fr {
	if c.privKey != nil {
		return c.privKey
	}
	privKey, err := hex.DecodeString(c.PrivKeyHex)
	if err != nil {
		panic(err)
	}
	c.privKey = bls12381.NewFr().FromBytes(privKey)
	return c.privKey
}
