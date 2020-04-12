package simulator

type StateGenerator interface {
	NextNode() Node
	NextWDS() WDS
	NodeLeft() int
	WDSLeft() int
	GetNodeDistribution() []int
}
