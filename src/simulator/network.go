package simulator

type HiveNet struct {
	servers []WDS
}

func (hn HiveNet) Run() {
	executedContracts := make([][]string, 0)
	for _, server := range hn.servers {
		runContracts := server.RunContracts()
		executedContracts = append(executedContracts, runContracts)
	}

	for i, server := range hn.servers {
		for j, snapshots := range executedContracts {
			if i != j {
				server.VerifySnapshots(snapshots)
			}
		}
	}

	for _, server := range hn.servers {
		log := server.GetLog()
		log.Print()
	}
}

func (hn HiveNet) AddWDS(s WDS) {
	hn.servers = append(hn.servers, s)
}

func (hn HiveNet) NetworkLog() map[string]interface{} {
	netLog := make(map[string]interface{})
	for _, wds := range hn.servers {
		netLog[wds.GetId()] = wds.GetLog()
	}
	return netLog
}
