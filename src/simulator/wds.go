package simulator

type WDS interface {
	GetId() string
	StartListener(chan bool)
	StartExecution() chan bool
	Conn() chan string
	GetTokens(string) float64
	EstablishLink(...WDS)
	GetLog() interface{}
	AssignNode(Node)
	AssignContract(Contract)
}
