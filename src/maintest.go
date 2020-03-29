package main

import (
	"fmt"

	"github.com/arstevens/hive-sim/src/controller"
)

func main() {
	ContractTest()
}

func ContractTest() {
	id := "transaction_1"
	action := 0 // simple token transfer
	transMap := make(map[string]float64)
	transMap["node1"] = 0.01
	transMap["node2"] = -0.01

	contract := controller.NewContract(id, action, transMap)

	// Marshalling test
	serial := contract.Marshal()
	fmt.Println(serial)
	var c2 controller.Contract
	c2.Unmarshal(serial)
	fmt.Println(c2.GetAmount("node1"), ":", c2.GetAmount("node2"))

}
