package simulator

type StateGenerator interface {
	NextNode() Node
	NextWDS() WDS
	NodesLeft() int
	WDSLeft() int
	GetAllNodes() []Node
	GetNodeDistribution() []int
}
