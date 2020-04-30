package controller

var (
	TotalTransactions      = "transactions_seen"
	SuccessfulTransactions = "transactions_passed"
	TotalSnapshots         = "snapshots_seen"
	SuccessfulSnapshots    = "snapshots_passed"
)

type BasicLog struct {
	wdsId                  string
	totalTransactions      int
	successfulTransactions int
	totalSnapshots         int
	successfulSnapshots    int
	wdsTokenLog            map[string]float64
	nodeTokenLog           map[string]float64
}

func (bl BasicLog) GetStats() map[string]int {
	statsMap := make(map[string]int)
	statsMap[TotalTransactions] = bl.totalTransactions
	statsMap[SuccessfulTransactions] = bl.successfulTransactions
	statsMap[TotalSnapshots] = bl.totalSnapshots
	statsMap[SuccessfulSnapshots] = bl.successfulSnapshots
	return statsMap
}

func (bl *BasicLog) IncTotalTransactions() {
	bl.totalTransactions++
}

func (bl *BasicLog) IncTotalSnapshots() {
	bl.totalSnapshots++
}

func (bl *BasicLog) IncSuccessfulTransactions() {
	bl.successfulTransactions++
}

func (bl *BasicLog) IncSuccessfulSnapshots() {
	bl.successfulSnapshots++
}
