package controller

import (
	"crypto/sha256"
	"log"
	"strconv"
	"strings"

	"github.com/arstevens/hive-sim/src/simulator"
)

type Contract struct {
	id              string
	action          int
	transactions    map[string]float64
	transactionHash string
	signatures      map[string]string
}

func NewEmptyContract() Contract {
	return Contract{}
}

func NewContract(id string, act int, trans map[string]float64) Contract {
	transMap := trans
	if transMap == nil {
		transMap = make(map[string]float64)
	}

	contract := Contract{id: id, action: act, transactions: transMap}
	return contract
}

func (c Contract) GetAmount(id string) float64 {
	return c.transactions[id]
}

func (c Contract) AddTransaction(id string, amount float64) {
	c.transactions[id] = amount
}

func (c Contract) SignContract(n simulator.Node) {
	hash := []byte(c.hashTransaction())
	signature := string(n.Sign(hash))
	c.signatures[n.Id()] = signature
}

func (c Contract) marshalTransaction() string {
	serial := c.id + "," + strconv.Itoa(c.action)
	for k, v := range c.transactions {
		pair := k + ":" + strconv.FormatFloat(v, 'E', -1, 64)
		serial += "," + pair
	}
	return serial
}

func (c Contract) marshalSignatures() string {
	serial := ""
	for nodeId, signature := range c.signatures {
		serial += "," + nodeId + " : " + signature
	}
	if len(serial) > 0 {
		serial = serial[1:]
	}
	return serial
}

func (c Contract) hashTransaction() string {
	serial := c.marshalTransaction()
	checksum := sha256.Sum256([]byte(serial))
	return string(checksum[:])
}

func (c Contract) Marshal() string {
	transactionSerial := c.marshalTransaction()
	snapshotSerial := c.marshalSignatures()
	serial := transactionSerial + ",snap," + snapshotSerial

	return serial
}

func (c Contract) Unmarshal(serial string) {
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

	c.signatures = make(map[string]string)
	for ; i < len(fields); i++ {
		pair := strings.Split(fields[i], ":")
		c.signatures[pair[0]] = pair[1]
	}
}
