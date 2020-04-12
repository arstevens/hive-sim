package controller

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"log"
	mrand "math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/arstevens/hive-sim/src/simulator"
)

type BasicLog struct {
	totalTransactions      int
	successfulTransactions int
	totalSnapshots         int
	successfulSnapshots    int
	wdsTokenLog            map[string]float64
	nodeTokenLog           map[string]float64
}

type BasicWDS struct {
	id          string
	tokens      float64
	inWDS       chan string
	outWDS      chan string
	nodesMutex  *sync.Mutex
	tokensMutex *sync.Mutex
	log         *BasicLog
	contracts   []simulator.Contract
	nodes       map[string]simulator.Node
	tokenMap    map[string]float64
}

func NewBasicWDS(id string, tokens float64) simulator.WDS {
	inWDS := make(chan string)
	var nodesMutex sync.Mutex
	var tokensMutex sync.Mutex
	log := BasicLog{
		wdsTokenLog:  make(map[string]float64),
		nodeTokenLog: make(map[string]float64),
	}
	contracts := make([]simulator.Contract, 0)
	nodes := make(map[string]simulator.Node)
	tokenMap := make(map[string]float64)

	return BasicWDS{
		id:          id,
		tokens:      tokens,
		inWDS:       inWDS,
		outWDS:      nil,
		nodesMutex:  &nodesMutex,
		tokensMutex: &tokensMutex,
		log:         &log,
		contracts:   contracts,
		nodes:       nodes,
		tokenMap:    tokenMap,
	}
}

func NewRandomBasicWDS() simulator.WDS {
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

	return NewBasicWDS(id, tokens)
}

func (bw *BasicWDS) AssignNode(n simulator.Node) {
	bw.nodes[n.Id()] = n
	bw.tokenMap[n.Id()] = n.Tokens()
}

func (bw *BasicWDS) AssignContract(c simulator.Contract) {
	bw.contracts = append(bw.contracts, c)
}

func (bw BasicWDS) Id() string {
	return bw.id
}

func (bw BasicWDS) Conn() chan string {
	return bw.inWDS
}

func (bw BasicWDS) Tokens(id string) float64 {
	if id == "" {
		return bw.tokens
	}
	bw.tokensMutex.Lock()
	tokens := bw.tokenMap[id]

	bw.tokensMutex.Unlock()
	return tokens
}

func (bw BasicWDS) GetLog() interface{} {
	log := bw.log
	log.wdsTokenLog = bw.tokenMap
	log.nodeTokenLog = make(map[string]float64)
	for nid, node := range bw.nodes {
		log.nodeTokenLog[nid] = node.Tokens()
	}

	return log
}

func (bw BasicWDS) getNode(id string) simulator.Node {
	bw.nodesMutex.Lock()
	node := bw.nodes[id]

	bw.nodesMutex.Unlock()
	return node
}

func (bw *BasicWDS) EstablishLink(servers ...simulator.WDS) {
	bw.inWDS = servers[0].Conn()
	bw.outWDS = servers[1].Conn()
}

func (bw *BasicWDS) updateTokenMap(snap simulator.Contract) {
	transaction := snap.GetTransactions()
	for id, amount := range transaction {
		bw.tokenMap[id] = bw.tokenMap[id] + amount
	}
}

func (bw *BasicWDS) StartListener(close chan bool) {
	go basicWDSConnListener(bw, close)
}

func (bw *BasicWDS) StartExecution() chan bool {
	finished := make(chan bool)
	go basicContractExecutor(bw, finished)
	return finished
}

func basicWDSConnListener(bw *BasicWDS, close chan bool) {
	for {
		select {
		case rawSnap := <-bw.inWDS:
			var snapshot simulator.Contract
			snapshot.Unmarshal(rawSnap)
			if basicVerifySnapshot(snapshot, bw) {
				bw.updateTokenMap(snapshot)
				bw.log.successfulSnapshots++
				bw.outWDS <- rawSnap
			}
			bw.log.totalSnapshots++
		case <-close:
			return
		}
	}
}

func basicVerifySnapshot(snapshot simulator.Contract, bw *BasicWDS) bool {
	transactions := snapshot.GetTransactions()
	for id, amount := range transactions {
		if amount > 0.0 && bw.Tokens(id) < amount {
			return false
		}
	}

	signatures := snapshot.GetSignatures()
	hash := []byte(snapshot.HashTransaction())
	for id, signature := range signatures {
		decodedSignature, err := base64.StdEncoding.DecodeString(signature)
		if err != nil {
			log.Fatal(err)
		}

		nodePk := bw.getNode(id).PublicKey()
		valid := Verify(hash, decodedSignature, nodePk)
		if !valid {
			return false
		}
	}
	return true
}

func basicContractExecutor(bw *BasicWDS, done chan bool) {
	for _, contract := range bw.contracts {
		subnet := basicSelectNodesForSubnet(bw.nodes, bw.nodesMutex)
		snapshot := basicExecuteContract(subnet, contract)
		if basicVerifySnapshot(snapshot, bw) {
			bw.updateTokenMap(snapshot)
			bw.log.successfulTransactions++
			bw.outWDS <- snapshot.Marshal()
		}
		bw.log.totalTransactions++
	}
	done <- true
}

func basicExecuteContract(nodes []simulator.Node, contract simulator.Contract) simulator.Contract {
	client := nodes[0]
	worker := nodes[1]
	filledContract := basicFillContract(client.Id(), worker.Id(), contract)
	agreedContract := basicAgreeOnContract(client, worker, filledContract)
	verifiedContract := basicVerifyContract(nodes[2:], agreedContract)

	return verifiedContract
}

/* Possibly have basic outline for sign checking on each node but allow the condition to
come from the nodes
	ex: if client.Evaluate(Contract) { sign }
*/
func basicAgreeOnContract(client simulator.Node, worker simulator.Node, contract simulator.Contract) simulator.Contract {
	if worker.EvaluateContract(contract, workerCode) {
		contract.SignContract(worker)
	}

	if client.EvaluateContract(contract, clientCode) {
		contract.SignContract(client)
	}

	return contract
}

func basicVerifyContract(nodes []simulator.Node, contract simulator.Contract) simulator.Contract {
	for _, node := range nodes {
		if node.EvaluateContract(contract, verifierCode) {
			contract.SignContract(node)
		}
	}
	return contract
}

func basicFillContract(clientId string, workerId string, c simulator.Contract) simulator.Contract {
	clientKey := strconv.Itoa(clientCode)
	c.AddTransaction(clientId, c.GetAmount(clientKey))
	c.DeleteTransaction(clientKey)

	workerKey := strconv.Itoa(workerCode)
	c.AddTransaction(workerId, c.GetAmount(workerKey))
	c.DeleteTransaction(workerKey)

	return c
}

func basicSelectNodesForSubnet(allNodes map[string]simulator.Node, lock *sync.Mutex) []simulator.Node {
	nodes := make([]simulator.Node, 8)
	i := 0
	lock.Lock()
	for _, node := range allNodes {
		if i > 7 {
			break
		}

		nodes[i] = node
		i++
	}
	lock.Unlock()
	return nodes
}
