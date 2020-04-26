package interactor

import (
	"testing"

	"github.com/arstevens/hive-sim/src/simulator"
)

func makeEvenNodeDist(wdsSize int, nodeSize int) []int {
	dist := make([]int, wdsSize)
	for i := 0; i < len(dist); i++ {
		dist[i] = nodeSize
	}
	return dist
}

func TestBasicGenerator(t *testing.T) {
	wdsSize := 5
	nodeSize := 10
	dist := makeEvenNodeDist(wdsSize, nodeSize)
	gen := NewBasicGenerator(wdsSize, nodeSize, 10, 1, dist)

	servers := make([]simulator.WDS, wdsSize)
	for i := 0; gen.WDSLeft() > 0; i++ {
		servers[i] = gen.NextWDS()
	}

	nodes := make([]simulator.Node, nodeSize*wdsSize)
	for i := 0; gen.NodesLeft() > 0; i++ {
		nodes[i] = gen.NextNode()
	}
}
