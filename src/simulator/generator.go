package simulator

type StateGenerator interface {
	NextNode() Node
	NextWDS() WDS
	NextContract() Contract
	NodeLeft() int
	WDSLeft() int
	ContractsLeft() int
	GetNodeDistribution() []int
	GetContractDistribution() []int
}
