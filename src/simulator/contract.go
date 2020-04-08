package simulator

type Contract interface {
	GetAmount(string) float64
	GetSignatures() map[string]string
	GetTransactions() map[string]float64
	AddTransaction(string, float64)
	DeleteTransaction(string)
	SignContract(Node)
	Marshal() string
	Unmarshal(serial string)
}
