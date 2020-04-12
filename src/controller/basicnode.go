package controller

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"log"
	mrand "math/rand"
	"strconv"
	"time"

	"github.com/arstevens/hive-sim/src/simulator"
)

const (
	clientCode   = 0
	workerCode   = 1
	verifierCode = 2
)

type BasicNode struct {
	RsaNode
	id     string
	tokens float64
}

func NewBasicNode(id string, sk *rsa.PrivateKey, tk float64) BasicNode {
	bn := BasicNode{RsaNode: RsaNode{secretKey: sk}, id: id, tokens: tk}
	return bn
}

func NewRandomBasicNode() BasicNode {
	// Generate random key
	sk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	// Generate Id from key
	pk := (sk.Public()).(*rsa.PublicKey)
	crypt, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pk, []byte("identification"), []byte{})
	if err != nil {
		log.Fatal(err)
	}
	nodeId := sha256.Sum256(crypt)
	id := base64.StdEncoding.EncodeToString(nodeId[:])

	// Generate random token count
	mrand.Seed(time.Now().UnixNano())
	tokens := float64(mrand.Intn(10)) + mrand.Float64()

	bn := BasicNode{RsaNode: RsaNode{secretKey: sk}, id: id, tokens: tokens}
	return bn
}

func (bn BasicNode) Id() string {
	return bn.id
}

func (bn BasicNode) Tokens() float64 {
	return bn.tokens
}

func (bn BasicNode) EvaluateContract(contract simulator.Contract, job int) bool {
	if job == clientCode {
		return bn.clientEvaluateContract(contract)
	} else if job == workerCode {
		return bn.workerEvaluateContract(contract)
	}
	return bn.verifierEvaluateContract(contract)
}

func (bn BasicNode) clientEvaluateContract(contract simulator.Contract) bool {
	return contract.GetAmount(bn.Id()) < 0.0
}

func (bn BasicNode) workerEvaluateContract(contract simulator.Contract) bool {
	return contract.GetAmount(bn.Id()) > 0.0
}

func (bn BasicNode) verifierEvaluateContract(contract simulator.Contract) bool {
	verifierKey := strconv.Itoa(verifierCode)
	return contract.GetAmount(verifierKey) > 0.0
}
