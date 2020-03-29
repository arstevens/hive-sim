package controller

import "github.com/arstevens/hive-sim/src/simulator"

type BasicWDS struct {
	id          string
	tokens      float64
	wdsConn     chan string
	nodeConn    chan string
	connections map[string]chan string
	nodes       map[string]simulator.Node
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
	return bw.nodes[id].Tokens()
}

func (bw *BasicWDS) EstablishLink(servers ...simulator.WDS) {
	for _, wds := range servers {
		bw.connections[wds.Id()] = wds.Conn()
	}
}
