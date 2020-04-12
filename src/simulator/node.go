package simulator

import "crypto/rsa"

type Node interface {
	Id() string
	Tokens() float64
	Sign([]byte) []byte
	PublicKey() *rsa.PublicKey
	EvaluateContract(Contract, int) bool
}
