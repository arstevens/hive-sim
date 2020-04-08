package controller

import (
	"encoding/base64"
	"log"

	"github.com/arstevens/hive-sim/src/simulator"
)

type BasicWDS struct {
	id        string
	tokens    float64
	inWDS     chan string
	outWDS    chan string
	contracts []simulator.Contract
	nodes     map[string]simulator.Node
	tokenMap  map[string]float64
}

func (bw *BasicWDS) Assign(n simulator.Node) {
	bw.nodes[n.Id()] = n
	bw.tokenMap[n.Id()] = n.Tokens()
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
	return bw.tokenMap[id]
}

func (bw *BasicWDS) EstablishLink(servers ...simulator.WDS) {
	bw.outWDS = servers[0].Conn()
}

func (bw *BasicWDS) updateTokenMap(snap Contract) {
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
			snapshot := NewEmptyContract()
			snapshot.Unmarshal(rawSnap)
			if basicVerifySnapshot(snapshot, bw) {
				bw.updateTokenMap(snapshot)
				bw.outWDS <- rawSnap
			}
		case <-close:
			return
		}
	}
}

func basicVerifySnapshot(snapshot Contract, bw *BasicWDS) bool {
	transactions := snapshot.GetTransactions()
	for id, amount := range transactions {
		if amount > 0.0 && bw.tokenMap[id] < amount {
			return false
		}
	}

	signatures := snapshot.GetSignatures()
	hash := []byte(snapshot.hashTransaction())
	for id, signature := range signatures {
		decodedSignature, err := base64.StdEncoding.DecodeString(signature)
		if err != nil {
			log.Fatal(err)
		}

		nodePk := bw.nodes[id].PublicKey()
		valid := Verify(hash, decodedSignature, nodePk)
		if !valid {
			return false
		}
	}
	return true
}

func basicContractExecutor(bw *BasicWDS, done chan bool) {
	for _, contract := range bw.contracts {
		subnet := basicSelectNodesForSubnet(bw.nodes)
		snapshot := basicExecuteContract(subnet, contract)
		if basicVerifySnapshot(snapshot, bw) {
			bw.outWDS <- snapshot.Marshal()
		}
	}
	done <- true
}

func basicExecuteContract(nodes []simulator.Node, contract Contract) Contract {
	client := nodes[0]
	worker := nodes[1]
	basicFillContract(client.Id(), worker.Id(), &contract)
	basicAgreeOnContract(client, worker, &contract)
	basicVerifyContract(nodes[2:], &contract)

	return contract
}

/* Possibly have basic outline for sign checking on each node but allow the condition to
come from the nodes
	ex: if client.Evaluate(Contract) { sign }
*/
func basicAgreeOnContract(client simulator.Node, worker simulator.Node, contract *Contract) {
	if worker.EvaluateContract(contract) {
		contract.SignContract(worker)
	}

	if contract.GetAmount(worker.Id()) > 0.0 {
		contract.SignContract(worker)
	}
}

func basicVerifyContract(nodes []simulator.Node, contract *Contract) {
	for _, node := range nodes {
		contract.SignContract(node)
	}
}

func basicFillContract(clientId string, workerId string, c *Contract) {
	c.AddTransaction(clientId, c.GetAmount("1"))
	c.DeleteTransaction("1")

	c.AddTransaction(workerId, c.GetAmount("2"))
	c.DeleteTransaction("2")
}

func basicSelectNodesForSubnet(allNodes map[string]simulator.Node) []simulator.Node {
	nodes := make([]simulator.Node, 8)
	i := 0
	for id, node := range allNodes {
		if i > 7 {
			break
		}

		nodes[i] = node
		i++
	}
	return nodes
}
