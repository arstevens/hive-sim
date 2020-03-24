package simulator

type WDS interface {
	Id() string
	Run() chan bool
	Tokens(string) float64
	EstablishLink(...WDS) error
	ActivityLog() map[string]interface{}
	Assign(Node)
}
