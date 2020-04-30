package simulator

type Log interface {
	GetStats() map[string]int
}
