package controller

import (
	"crypto/sha256"
	"encoding/base64"
	"log"
	mrand "math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/arstevens/hive-sim/src/simulator"
)

var totalContracts = 0

type BasicContract struct {
	id              string
	action          int
	transactions    map[string]float64
	transactionHash string
	signatures      map[string]string
}

func NewBasicContract(id string, act int, trans map[string]float64) BasicContract {
	transMap := trans
	if transMap == nil {
		transMap = make(map[string]float64)
	}

	contract := BasicContract{id: id, action: act, transactions: transMap, signatures: make(map[string]string)}
	return contract
}

func NewRandomBasicContract(transLimit int) BasicContract {
	transMap := make(map[string]float64)

	mrand.Seed(time.Now().UnixNano())
	transVal := float64(mrand.Intn(transLimit-1)) + mrand.Float64()
	transMap["A"] = transVal
	transMap["B"] = -transVal
	transId := "T" + strconv.Itoa(totalContracts)
	totalContracts++

	return NewBasicContract(transId, 0, transMap)
}

func (c BasicContract) GetAmount(id string) float64 {
	return c.transactions[id]
}

func (c BasicContract) GetSignatures() map[string]string {
	return c.signatures
}

func (c BasicContract) GetTransactions() map[string]float64 {
	return c.transactions
}

func (c *BasicContract) AddTransaction(id string, amount float64) {
	c.transactions[id] = amount
}

func (c *BasicContract) DeleteTransaction(id string) {
	delete(c.transactions, id)
}

func (c *BasicContract) SignContract(n simulator.Node) {
	hash := []byte(c.HashTransaction())
	signature := n.Sign(hash)
	encodedSignature := base64.StdEncoding.EncodeToString(signature)
	c.signatures[n.GetId()] = encodedSignature
}

func (c BasicContract) marshalTransaction() string {
	serial := c.id + "," + strconv.Itoa(c.action)
	for k, v := range c.transactions {
		pair := k + ":" + strconv.FormatFloat(v, 'E', -1, 64)
		serial += "," + pair
	}
	return serial
}

func (c BasicContract) marshalSignatures() string {
	serial := ""
	for nodeId, signature := range c.signatures {
		serial += "," + nodeId + ":" + signature
	}
	if len(serial) > 0 {
		serial = serial[1:]
	}
	return serial
}

func (c BasicContract) HashTransaction() string {
	serial := c.marshalTransaction()
	checksum := sha256.Sum256([]byte(serial))
	return string(checksum[:])
}

func (c BasicContract) Marshal() string {
	transactionSerial := c.marshalTransaction()
	snapshotSerial := c.marshalSignatures()

	serial := transactionSerial + ",snap"
	if len(snapshotSerial) > 0 {
		serial += "," + snapshotSerial
	}

	return serial
}

func (c *BasicContract) Unmarshal(serial string) {
	fields := strings.Split(serial, ",")
	c.id = fields[0]
	action, err := strconv.ParseInt(fields[1], 10, 32)
	if err != nil {
		log.Fatal(err)
	}

	c.action = int(action)
	c.transactions = make(map[string]float64)
	i := 2
	for ; fields[i] != "snap"; i++ {
		pair := strings.Split(fields[i], ":")
		partyId := pair[0]
		amount, err := strconv.ParseFloat(pair[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		c.transactions[partyId] = amount
	}
	i++

	c.signatures = make(map[string]string)
	for ; i < len(fields); i++ {
		pair := strings.Split(fields[i], ":")
		if len(pair) > 1 {
			c.signatures[pair[0]] = pair[1]
		}
	}
}
