package simulator

type StateGenerator interface {
	NextNode() Node
	NextWDS() WDS
	NodeLefts() int
	WDSLeft() int
}
