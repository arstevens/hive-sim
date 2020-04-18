package controller

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	mrand "math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/arstevens/hive-sim/src/simulator"
)

const (
	clientKey = "A"
	workerKey = "B"
)

var totalContracts = 0

type BasicContract struct {
	id               string
	action           int
	startingBalances map[string]float64
	transactions     map[string]float64
	transactionHash  string
	signatures       map[string]string
}

func NewBasicContract(id string, act int, trans map[string]float64) *BasicContract {
	transMap := trans
	if transMap == nil {
		transMap = make(map[string]float64)
	}

	contract := BasicContract{id: id, action: act, transactions: transMap,
		signatures: make(map[string]string), startingBalances: make(map[string]float64)}
	return &contract
}

func NewRandomBasicContract(transLimit int) *BasicContract {
	transMap := make(map[string]float64)

	mrand.Seed(time.Now().UnixNano())
	var transVal float64
	if transLimit < 2 {
		transVal = mrand.Float64()
	} else {
		transVal = float64(mrand.Intn(transLimit-1)) + mrand.Float64()
	}
	transMap[clientKey] = -transVal
	transMap[workerKey] = transVal
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

func (c BasicContract) GetStartingBalance(id string) float64 {
	return c.startingBalances[id]
}

func (c *BasicContract) SetStartingBalances(sb map[string]float64) {
	c.startingBalances = sb
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

	fmt.Println("Node Id: ", n.GetId())
	fmt.Println("Hash")
	fmt.Println(hash)
	fmt.Println("Signature")
	fmt.Println(signature)
	fmt.Println()
	encodedSignature := base64.StdEncoding.EncodeToString(signature)
	c.signatures[n.GetId()] = encodedSignature
}

func (c BasicContract) marshalTransaction() string {
	serial := c.id + "," + strconv.Itoa(c.action)

	flipTransactions := make(map[float64]string)
	keys := make([]float64, 0)
	for k, v := range c.transactions {
		flipTransactions[v] = k
		keys = append(keys, v)
	}
	sort.Float64s(keys)

	for _, v := range keys {
		id := flipTransactions[v]
		pair := id + ":" + strconv.FormatFloat(v, 'E', -1, 64)
		serial += "," + pair
	}

	return serial
}

func (c BasicContract) marshalStartingBalances() string {
	serial := ""
	for id, value := range c.startingBalances {
		serial += "," + id + ":" + strconv.FormatFloat(value, 'E', -1, 64)
	}
	if len(serial) > 0 {
		return serial[1:]
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
	startingBalanceSerial := c.marshalStartingBalances()

	serial := transactionSerial + ",snap"
	if len(snapshotSerial) > 0 {
		serial += "," + snapshotSerial
	}
	serial += ",bal," + startingBalanceSerial

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
	for ; fields[i] != "bal"; i++ {
		pair := strings.Split(fields[i], ":")
		if len(pair) > 1 {
			c.signatures[pair[0]] = pair[1]
		}
	}
	i++

	c.startingBalances = make(map[string]float64)
	for ; i < len(fields); i++ {
		pair := strings.Split(fields[i], ":")
		partyId := pair[0]
		amount, err := strconv.ParseFloat(pair[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		c.startingBalances[partyId] = amount
	}
}
