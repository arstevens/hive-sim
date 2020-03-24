package simulator

import "crypto/rsa"

type Node interface {
	Activate()
	Id() string
	Tokens() float64
	Sign([]byte) []byte
	Verify([]byte, rsa.PublicKey) bool
	JoinVerification(Node, Node) error
	EnterContract(Node, Node) error
	ExecuteNextContract() error
}
