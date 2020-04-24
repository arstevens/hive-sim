package interactor

import (
	"github.com/arstevens/hive-sim/src/controller"
	"github.com/arstevens/hive-sim/src/simulator"
)

const (
	RANDOM_CGEN = 0
	EVEN_CGEN   = 1
)

type BasicGenerator struct {
	nodesLeft            int
	wdsLeft              int
	contractsLeft        int
	contractGenType      int
	contractLimit        int
	nodeDistribution     []int
	contractDistribution []int
}

func NewBasicGenerator(nodeCount int, wdsCount int, contractCount int, genType int,
	contractLimit int, nodeDist []int) BasicGenerator {
	return BasicGenerator{
		nodesLeft:        nodeCount,
		wdsLeft:          wdsCount,
		contractsLeft:    contractCount,
		contractGenType:  genType,
		contractLimit:    contractLimit,
		nodeDistribution: nodeDist,
	}
}

func (bg BasicGenerator) NodesLeft() int {
	return bg.nodesLeft
}

func (bg BasicGenerator) WDSLeft() int {
	return bg.wdsLeft
}

func (bg *BasicGenerator) NextNode() simulator.Node {
	node := controller.NewRandomBasicNode()
	bg.nodesLeft = bg.nodesLeft - 1
	return node
}

func (bg *BasicGenerator) NextWDS() simulator.WDS {
	wds := controller.NewRandomBasicWDS()
	bg.wdsLeft = bg.wdsLeft - 1
	return wds
}

func (bg BasicGenerator) GetNodeDistribution() []int {
	return bg.nodeDistribution
}
