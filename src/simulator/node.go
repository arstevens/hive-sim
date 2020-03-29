package simulator

import "crypto/rsa"

type Node interface {
	Id() string
	Conn() chan string
	Tokens() float64
	Sign([]byte) []byte
	Verify([]byte, []byte, *rsa.PublicKey) bool
	JoinVerification(chan string, chan string) chan bool
	EnterContract(chan string, chan string) chan bool
	ExecuteNextContract()
}
