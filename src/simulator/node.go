package simulator

import "crypto/rsa"

type Node interface {
	Id() string
	Conn() chan string
	Tokens() float64
	Sign([]byte) []byte
	PublicKey() *rsa.PublicKey
	EvaluateContract(Contract) bool
}
