package simulator

type WDS interface {
	Id() string
	Run() chan bool
	Conn() chan string
	Tokens(string) float64
	EstablishLink(...WDS)
	ActivityLog() map[string]interface{}
	Assign(Node)
}
