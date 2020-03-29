package controller

type BasicNode struct {
	RsaNode
	id             string
	wds            chan string
	conn           chan string
	tokens         float64
	contracts      []Contract
	activeContract *Contract
}

func (bn BasicNode) Id() string {
	return bn.id
}

func (bn BasicNode) Tokens() float64 {
	return bn.tokens
}

func (bn BasicNode) Conn() chan string {
	return bn.conn
}

func (bn BasicNode) ExecuteNextContract() {
	// Ask WDS to setup the verification subnet
	nextContract := bn.contracts[0]
	bn.contracts = bn.contracts[1:]

	serialContract := nextContract.Marshal()
	bn.wds <- serialContract
	bn.activeContract = &nextContract
}

func recvContract(in chan string) Contract {
	var contract Contract
	serialContract := <-in
	contract.Unmarshal(serialContract)
	return contract
}

func (bn BasicNode) EnterContract(transNode chan string, verfNode chan string) {
	if bn.activeContract != nil {
		activeEnterContract(bn, transNode, verfNode)
		return
	}
	inactiveEnterContract(bn, transNode, verfNode)
}

func activeEnterContract(bn BasicNode, transNode chan string, verfNode chan string) {
	serialContract := bn.activeContract.Marshal()
	transNode <- serialContract

	signedContract := recvContract(verfNode)
	snapshot := GenerateSnapshot(signedContract)
	bn.wds <- snapshot
}

func inactiveEnterContract(bn BasicNode, transNode chan string, verfNode chan string) {
	contract := recvContract(transNode)

	// Verify that you are gaining tokens in the transaction
	nodeId := bn.Id()
	if contract.GetAmount(nodeId) > 0.0 {
		contract.SignContract(bn)
		verfNode <- contract.Marshal()
		return
	}

	// If transaction is bad, kill the subnet
	verfNode <- "kill"
}

// Simply signs contract and passes it along
func (bn BasicNode) JoinVerification(leftNode chan string, rightNode chan string) {
	contract := recvContract(leftNode)
	contract.SignContract(bn)
	rightNode <- contract.Marshal()
}
