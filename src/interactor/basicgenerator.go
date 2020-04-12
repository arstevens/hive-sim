package interactor

import (
	"github.com/arstevens/hive-sim/src/controller"
	"github.com/arstevens/hive-sim/src/simulator"
)

type BasicGenerator struct {
	nodesLeft int
	wdsLeft   int
}

func (bg BasicGenerator) NodesLeft() int {
	return bg.nodesLeft
}

func (bg BasicGenerator) WDSLeft() int {
	return bg.wdsLeft
}

func (bg BasicGenerator) NextNode() simulator.Node {
	return controller.NewRandomBasicNode()
}
