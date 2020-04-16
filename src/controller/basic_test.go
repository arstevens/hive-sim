package controller

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/arstevens/hive-sim/src/simulator"
)

func TestBasicNode(t *testing.T) {
	fmt.Println("\n-------- START OF BASIC NODE TESTING --------")
	bn := NewRandomBasicNode()
	bn.SetTokens(0.1)

	transMap := make(map[string]float64)
	transMap[bn.GetId()] = -0.05                   // Client
	transMap["FakeNodeId"] = 0.05                  // Worker
	transMap[strconv.Itoa(verifierCode)] = 0.00001 //Verifier
	contract := NewBasicContract("t1", 0, transMap)

	fmt.Println("Evaluating Contract as Client: ", bn.EvaluateContract(contract, clientCode))
	fmt.Println("Evaluating Contract as Worker: ", bn.EvaluateContract(contract, workerCode))
	fmt.Println("Evaluating Contract as Verifier: ", bn.EvaluateContract(contract, verifierCode))

	id := bn.GetId()
	fmt.Println(id)
	fmt.Println("-------- END OF BASIC NODE TESTING --------")
}

func TestBasicWDS(t *testing.T) {
	fmt.Println("\n-------- START OF BASIC WDS TESTING --------")
	wdsSize := 3
	nodesPerWds := 10
	contractsPerWds := 50
	transLimit := 1

	fmt.Printf("Test Parameters: WDS_SIZE: %d, NODES_PER: %d, CONTRACTS_PER: %d, TRANS_LIM: %d\n", wdsSize, nodesPerWds, contractsPerWds, transLimit)

	// Generate WDSs and establish cyclic network
	wds := make([]simulator.WDS, wdsSize)
	killChannels := make([]chan bool, wdsSize)

	wds[0] = NewRandomBasicWDS()
	killChannels[0] = make(chan bool)
	for i := 1; i < wdsSize; i++ {
		wds[i] = NewRandomBasicWDS()
		wds[i-1].EstablishLink(wds[i])
		killChannels[i] = make(chan bool)
	}
	wds[wdsSize-1].EstablishLink(wds[0])

	// Generate Nodes and Contracts
	for idx, server := range wds {
		for i := 0; i < nodesPerWds; i++ {
			node := NewRandomBasicNode()
			server.AssignNode(node)
		}

		for i := 0; i < contractsPerWds; i++ {
			contract := NewRandomBasicContract(transLimit)
			server.AssignContract(contract)
		}

		server.StartListener(killChannels[idx])
	}

	// Start Network Execution
	finishChannels := make([]chan bool, wdsSize)
	for idx, server := range wds {
		finishChannels[idx] = server.StartExecution()
	}

	// Wait for execution to end
	for _, finishChannel := range finishChannels {
		<-finishChannel
	}

	// Kill WDS Listeners
	for i := 0; i < wdsSize; i++ {
		killChannels[i] <- true
	}

	// Retrieve and Print Logs
	logs := make([]*BasicLog, wdsSize)
	for i := 0; i < wdsSize; i++ {
		log := wds[i].GetLog().(*BasicLog)
		logs[i] = log
		log.Print()
	}

	fmt.Println("-------- END OF BASIC WDS TESTING --------")
}

func TestCrypto(t *testing.T) {
	fmt.Println("\n-------- START OF CRYPTO TESTING --------")
	node := NewRandomBasicNode()
	data := []byte{1, 2, 3}
	s := node.Sign(data)
	fmt.Println(RsaVerify(data, s, node.PublicKey()))
	fmt.Println("\n-------- END OF CRYPTO TESTING --------")
}

func TestContract(t *testing.T) {
	fmt.Println("\n-------- START OF CONTRACT TESTING --------")
	bn1 := NewRandomBasicNode()
	bn2 := NewRandomBasicNode()

	id := "transaction_1"
	action := 0 // simple token transfer
	transMap := make(map[string]float64)
	transMap[bn1.GetId()] = 0.01
	transMap[bn2.GetId()] = -0.01

	contract := NewBasicContract(id, action, transMap)
	contract.SignContract(bn1)
	contract.SignContract(bn2)
	serial := contract.Marshal()
	fmt.Println(contract.Marshal())
	var contract2 BasicContract
	contract2.Unmarshal(serial)
	fmt.Println(contract2.Marshal())
	fmt.Println("-------- END OF CONTRACT TESTING --------")
}
