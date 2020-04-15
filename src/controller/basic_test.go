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

	fmt.Println("Evaluating Contract as Client: ", bn.EvaluateContract(&contract, clientCode))
	fmt.Println("Evaluating Contract as Worker: ", bn.EvaluateContract(&contract, workerCode))
	fmt.Println("Evaluating Contract as Verifier: ", bn.EvaluateContract(&contract, verifierCode))

	id := bn.GetId()
	fmt.Println(id)
	fmt.Println("-------- END OF BASIC NODE TESTING --------")
}

func TestBasicWDS(t *testing.T) {
	fmt.Println("\n-------- START OF BASIC WDS TESTING --------")
	wdsSize := 3
	nodesPerWds := 5

	// Generate WDSs and establish cyclic network
	wds := make([]simulator.WDS, wdsSize)
	wds[0] = NewRandomBasicWDS()
	for i := 1; i < wdsSize; i++ {
		wds[i] = NewRandomBasicWDS()
		wds[i-1].EstablishLink(wds[i])
	}
	wds[wdsSize-1].EstablishLink(wds[0])

	// Generate Nodes and assign to WDSs
	for _, server := range wds {
		for i := 0; i < nodesPerWds; i++ {
			node := NewRandomBasicNode()
			server.AssignNode(node)
		}
	}

	fmt.Println("-------- END OF BASIC WDS TESTING --------")
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
