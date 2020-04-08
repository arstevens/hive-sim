package controller

import (
	"fmt"
	"testing"
)

func TestBasicNode(t *testing.T) {
	fmt.Println("\n-------- START OF BASIC NODE TESTING --------")
	bn := NewRandomBasicNode(nil, nil)
	id := bn.Id()
	fmt.Println(id)
	fmt.Println("-------- END OF BASIC NODE TESTING --------")
}

/*
func TestContract(t *testing.T) {
	fmt.Println("\n-------- START OF CONTRACT TESTING --------")
	bn1 := NewRandomBasicNode(nil, nil)
	bn2 := NewRandomBasicNode(nil, nil)

	id := "transaction_1"
	action := 0 // simple token transfer
	transMap := make(map[string]float64)
	transMap[bn1.Id()] = 0.01
	transMap[bn2.Id()] = -0.01

	contract := NewContract(id, action, transMap)
	contract.SignContract(bn1)
	contract.SignContract(bn2)
	serial := contract.Marshal()
	fmt.Println(contract.Marshal())
	contract2 := NewEmptyContract()
	contract2.Unmarshal(serial)
	fmt.Println(contract2.Marshal())
	fmt.Println("-------- END OF CONTRACT TESTING --------")
}
*/

func TestGen(t *testing.T) {
	ints := generateConsecutiveIntegers(0, 50, 20)
	fmt.Println(ints)
}
