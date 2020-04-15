package controller

import (
	"fmt"
	"testing"
)

func TestBasicNode(t *testing.T) {
	fmt.Println("\n-------- START OF BASIC NODE TESTING --------")
	bn := NewRandomBasicNode()
	id := bn.Id()
	fmt.Println(id)
	fmt.Println("-------- END OF BASIC NODE TESTING --------")
}

func TestContract(t *testing.T) {
	fmt.Println("\n-------- START OF CONTRACT TESTING --------")
	bn1 := NewRandomBasicNode()
	bn2 := NewRandomBasicNode()

	id := "transaction_1"
	action := 0 // simple token transfer
	transMap := make(map[string]float64)
	transMap[bn1.Id()] = 0.01
	transMap[bn2.Id()] = -0.01

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
