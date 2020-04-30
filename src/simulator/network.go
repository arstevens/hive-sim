package simulator

type HiveNet struct {
	servers   []WDS
	formatter Formatter
}

func NewHiveNet(format Formatter) *HiveNet {
	return &HiveNet{
		servers:   make([]WDS, 0),
		formatter: format,
	}
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

	hn.formatter.Format(hn.networkLog())
	hn.formatter.Save()
}

func (hn *HiveNet) AddWDS(s WDS) {
	hn.servers = append(hn.servers, s)
}

func (hn HiveNet) networkLog() map[string]Log {
	netLog := make(map[string]Log)
	for _, wds := range hn.servers {
		netLog[wds.GetId()] = wds.GetLog()
	}
	return netLog
}
