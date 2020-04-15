package simulator

import "crypto/rsa"

type Node interface {
	GetId() string
	GetTokens() float64
	SetTokens(float64)
	Sign([]byte) []byte
	PublicKey() *rsa.PublicKey
	EvaluateContract(Contract, int) bool
}
