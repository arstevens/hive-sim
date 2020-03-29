package controller

import (
	"encoding/base64"
	"log"

	"github.com/arstevens/hive-sim/src/simulator"
)

type BasicWDS struct {
	id       string
	tokens   float64
	wdsConn  chan string
	nodeConn chan string
	inWDS    chan string
	outWDS   chan string
	nodes    map[string]simulator.Node
	tokenMap map[string]float64
}

func (bw *BasicWDS) Assign(n simulator.Node) {
	bw.nodes[n.Id()] = n
}

func (bw BasicWDS) Id() string {
	return bw.id
}

func (bw BasicWDS) Conn() chan string {
	return bw.wdsConn
}

func (bw BasicWDS) Tokens(id string) float64 {
	if id == "" {
		return bw.tokens
	}
	return bw.tokenMap[id]
}

func (bw *BasicWDS) EstablishLink(servers ...simulator.WDS) {
	bw.outWDS = servers[0].Conn()
}

func basicWDSConnListener(bw BasicWDS) {
	for {
		select {
		case rawSnap := <-bw.inWDS:
			snapshot := NewEmptyContract()
			snapshot.Unmarshal(rawSnap)
		}
	}
}

func basicVerifySnapshot(snapshot Contract, bw BasicWDS) bool {
	signatures := snapshot.GetSignatures()
	hash := []byte(snapshot.hashTransaction())
	for id, signature := range signatures {
		decodedSignature, err := base64.StdEncoding.DecodeString(signature)
		if err != nil {
			log.Fatal(err)
		}

		nodePk := bw.nodes[id].PublicKey()
		valid := Verify(hash, decodedSignature, nodePk)
		if !valid {
			return false
		}
	}
	return true
}
