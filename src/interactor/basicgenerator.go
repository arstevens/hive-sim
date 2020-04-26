package interactor

import (
	"github.com/arstevens/hive-sim/src/controller"
	"github.com/arstevens/hive-sim/src/simulator"
)

type BasicGenerator struct {
	nodesLeft        int
	wdsLeft          int
	contractLimit    int
	transactionLimit int
	allNodes         []simulator.Node
	nodeDistribution []int
}

func NewBasicGenerator(wdsCount int, nodeCount int, contractLimit int, transLimit int, nodeDist []int) *BasicGenerator {
	return &BasicGenerator{
		nodesLeft:        nodeCount * wdsCount,
		wdsLeft:          wdsCount,
		contractLimit:    contractLimit,
		transactionLimit: transLimit,
		allNodes:         make([]simulator.Node, 0),
		nodeDistribution: nodeDist,
	}
}

func (bg BasicGenerator) NodesLeft() int {
	return bg.nodesLeft
}

func (bg BasicGenerator) WDSLeft() int {
	return bg.wdsLeft
}

func (bg BasicGenerator) GetAllNodes() []simulator.Node {
	allNodes := bg.allNodes
	return allNodes
}

func (bg *BasicGenerator) NextNode() simulator.Node {
	node := controller.NewRandomBasicNode()
	bg.allNodes = append(bg.allNodes, node)
	bg.nodesLeft = bg.nodesLeft - 1
	return node
}

func (bg *BasicGenerator) NextWDS() simulator.WDS {
	wds := controller.NewRandomBasicWDS()
	bg.wdsLeft = bg.wdsLeft - 1
	contractCount := bg.contractLimit

	for i := 0; i < contractCount; i++ {
		contract := controller.NewRandomBasicContract(bg.transactionLimit)
		wds.AssignContract(contract)
	}

	return wds
}

func (bg BasicGenerator) GetNodeDistribution() []int {
	return bg.nodeDistribution
}
