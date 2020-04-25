package simulator

import (
	"testing"

	"github.com/arstevens/hive-sim/src/interactor"
)

func SimulatorTest(t *testing.T) {
	wdsSize := 5
	nodeSize := 10
	contractLimit := 10
	transLimit := 1

	nodeDist := make([]int, nodeSize*wdsSize)
	for i := 0; i < wdsSize*nodeSize; i++ {
		nodeDist[i] = nodeSize
	}
	generator := interactor.NewBasicGenerator(nodeSize, wdsSize, contractLimit, transLimit, nodeDist)

	network := NewHiveNet()
	AllocateResources(network, generator)
}
