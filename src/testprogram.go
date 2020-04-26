package main

import (
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
	generator := interactor.NewBasicGenerator(wdsSize, nodeSize, contractLimit, transLimit, nodeDist)

	network := simulator.NewHiveNet()
	simulator.AllocateResources(network, generator)
	network.Run()
}
