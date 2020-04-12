package interactor

import (
	mrand "math/rand"
	"strconv"
	"time"

	"github.com/arstevens/hive-sim/src/controller"
	"github.com/arstevens/hive-sim/src/simulator"
)

var totalContracts = 0

const (
	RANDOM_CGEN = 0
	EVEN_CGEN   = 1
)

type BasicGenerator struct {
	nodesLeft        int
	wdsLeft          int
	contractGenType  int
	contractLimit    int
	nodeDistribution []int
}

func NewBasicGenerator(nodeCount int, wdsCount int, genType int, contractLimit int, nodeDist []int) BasicGenerator {
	return BasicGenerator{
		nodesLeft:        nodeCount,
		wdsLeft:          wdsCount,
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

func (bg BasicGenerator) NextNode() simulator.Node {
	return controller.NewRandomBasicNode()
}

func (bg BasicGenerator) NextWDS() simulator.WDS {
	wds := controller.NewRandomBasicWDS()
	contractCount := bg.contractLimit
	if bg.contractGenType == RANDOM_CGEN {
		mrand.Seed(time.Now().UnixNano())
		contractCount = mrand.Intn(bg.contractLimit)
	}

	for i := 0; i < contractCount; i++ {
		contract := createRandomBasicContract(1)
		wds.AssignContract(contract)
	}

	return wds
}

func (bg BasicGenerator) GetNodeDistribution() []int {
	return bg.nodeDistribution
}

func createRandomBasicContract(transLimit int) controller.BasicContract {
	transMap := make(map[string]float64)

	mrand.Seed(time.Now().UnixNano())
	transVal := float64(mrand.Intn(transLimit-1)) + mrand.Float64()
	transMap["A"] = transVal
	transMap["B"] = -transVal
	transId := "T" + strconv.Itoa(totalContracts)
	totalContracts++

	return controller.NewBasicContract(transId, 0, transMap)
}
