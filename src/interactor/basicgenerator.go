package interactor

import (
	mrand "math/rand"
	"time"

	"github.com/arstevens/hive-sim/src/controller"
	"github.com/arstevens/hive-sim/src/simulator"
)

const (
	RANDOM_CGEN = 0
	EVEN_CGEN   = 1
)

type BasicGenerator struct {
	nodesLeft        int
	wdsLeft          int
	contractGenType  int
	contractLimit    int
	allNodes         []simulator.Node
	nodeDistribution []int
}

func NewBasicGenerator(nodeCount int, wdsCount int, genType int, contractLimit int, nodeDist []int) BasicGenerator {
	return BasicGenerator{
		nodesLeft:        nodeCount,
		wdsLeft:          wdsCount,
		contractGenType:  genType,
		contractLimit:    contractLimit,
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
	if bg.contractGenType == RANDOM_CGEN {
		mrand.Seed(time.Now().UnixNano())
		contractCount = mrand.Intn(bg.contractLimit)
	}

	for i := 0; i < contractCount; i++ {
		contract := controller.NewRandomBasicContract(1)
		wds.AssignContract(contract)
	}

	return wds
}

func (bg BasicGenerator) GetNodeDistribution() []int {
	return bg.nodeDistribution
}
