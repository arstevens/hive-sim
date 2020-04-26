package main

import (
	"fmt"

	"github.com/arstevens/hive-sim/src/interactor"
	"github.com/arstevens/hive-sim/src/simulator"
)

func main() {
	wdsSize := 5
	nodeSize := 10
	contractLimit := 10
	transLimit := 1

	nodeDist := make([]int, wdsSize)
	for i := 0; i < wdsSize; i++ {
		nodeDist[i] = nodeSize
	}
	fmt.Println(len(nodeDist))
	generator := interactor.NewBasicGenerator(wdsSize, nodeSize, contractLimit, transLimit, nodeDist)

	network := simulator.NewHiveNet()
	fmt.Println(generator.WDSLeft())
	simulator.AllocateResources(network, generator)
	network.Run()
}
