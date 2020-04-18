package simulator

type Contract interface {
	GetAmount(string) float64
	GetSignatures() map[string]string
	GetTransactions() map[string]float64
	GetStartingBalance(string) float64
	SetStartingBalances(map[string]float64)
	AddTransaction(string, float64)
	HashTransaction() string
	DeleteTransaction(string)
	SignContract(Node)
	Marshal() string
	Unmarshal(serial string)
}
