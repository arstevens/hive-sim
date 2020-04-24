package simulator

type HiveNet struct {
	servers []WDS
}

func (hn HiveNet) Run() {

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
