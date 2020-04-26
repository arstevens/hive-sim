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

	nodeDist := make([]int, nodeSize*wdsSize)
	for i := 0; i < wdsSize*nodeSize; i++ {
		nodeDist[i] = nodeSize
	}
	generator := interactor.NewBasicGenerator(nodeSize, wdsSize, contractLimit, transLimit, nodeDist)
	g2 := &generator

	network := simulator.NewHiveNet()
	simulator.AllocateResources(network, g2)
}
