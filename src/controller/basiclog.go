package controller

import "fmt"

type BasicLog struct {
	wdsId                  string
	totalTransactions      int
	successfulTransactions int
	totalSnapshots         int
	successfulSnapshots    int
	wdsTokenLog            map[string]float64
	nodeTokenLog           map[string]float64
}

func (bl BasicLog) Print() {
	fmt.Printf(`Id: %s {
		Total Transactions: %d,
		Successful Transactions: %d,
		Total Snapshots: %d,
		Successful Snapshot: %d,
		}`, bl.wdsId, bl.totalTransactions,
		bl.successfulTransactions, bl.totalSnapshots,
		bl.successfulSnapshots)
	fmt.Println()
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
