package controller

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"log"
	mrand "math/rand"
	"time"
)

type BasicNode struct {
	RsaNode
	id             string
	wds            chan string
	conn           chan string
	tokens         float64
	contracts      []Contract
	activeContract *Contract
}

func NewBasicNode(id string, sk *rsa.PrivateKey, tk float64, contracts []Contract, wds chan string) BasicNode {
	bn := BasicNode{RsaNode: RsaNode{secretKey: sk}, id: id, tokens: tk, contracts: contracts, wds: wds, activeContract: nil}
	bn.conn = make(chan string)
	return bn
}

func NewRandomBasicNode(contracts []Contract, wds chan string) BasicNode {
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

	bn := BasicNode{RsaNode: RsaNode{secretKey: sk}, id: id, tokens: tokens, contracts: contracts, wds: wds, activeContract: nil}
	return bn
}

func (bn BasicNode) Id() string {
	return bn.id
}

func (bn BasicNode) Tokens() float64 {
	return bn.tokens
}

func (bn BasicNode) Conn() chan string {
	return bn.conn
}

func (bn BasicNode) ExecuteNextContract() {
	// Ask WDS to setup the verification subnet
	nextContract := bn.contracts[0]
	bn.contracts = bn.contracts[1:]

	serialContract := nextContract.Marshal()
	bn.wds <- serialContract
	bn.activeContract = &nextContract
}

func recvContract(in chan string) Contract {
	var contract Contract
	serialContract := <-in
	contract.Unmarshal(serialContract)
	return contract
}

func (bn BasicNode) EnterContract(transNode chan string, verfNode chan string) chan bool {
	if bn.activeContract != nil {
		activeEnterContract(bn, transNode, verfNode)
		return nil
	}
	inactiveEnterContract(bn, transNode, verfNode)
	return nil
}

func activeEnterContract(bn BasicNode, transNode chan string, verfNode chan string) {
	serialContract := bn.activeContract.Marshal()
	transNode <- serialContract

	signedContract := recvContract(verfNode)
	snapshot := signedContract.Marshal()
	bn.wds <- snapshot
}

func inactiveEnterContract(bn BasicNode, transNode chan string, verfNode chan string) {
	contract := recvContract(transNode)

	// Verify that you are gaining tokens in the transaction
	nodeId := bn.Id()
	if contract.GetAmount(nodeId) > 0.0 {
		contract.SignContract(bn)
		verfNode <- contract.Marshal()
		return
	}

	// If transaction is bad, kill the subnet
	verfNode <- "kill"
}

// Simply signs contract and passes it along
func (bn BasicNode) JoinVerification(leftNode chan string, rightNode chan string) chan bool {
	contract := recvContract(leftNode)
	contract.SignContract(bn)
	rightNode <- contract.Marshal()
	return nil
}
