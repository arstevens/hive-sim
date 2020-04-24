package simulator

type WDS interface {
	GetId() string
	RunContracts() []string
	VerifySnapshots(snaps []string)
	GetTokens(string) float64
	GetLog() interface{}
	AssignNode(Node)
	AssignContract(Contract)
	SetMasterKeyList([]Node)
}
