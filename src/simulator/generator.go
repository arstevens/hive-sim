package simulator

type StateGenerator interface {
	NextNode() Node
	NextWDS() WDS
	NodeLeft() int
	GetAllNodes() []Node
	WDSLeft() int
	GetNodeDistribution() []int
}
