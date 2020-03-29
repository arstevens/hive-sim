package simulator

type HiveNet struct {
	servers []WDS
}

func (hn HiveNet) Run() {
	closeChans := make([]chan bool, len(hn.servers))
	for i, wds := range hn.servers {
		wdsClose := wds.Run()
		closeChans[i] = wdsClose
	}

	for _, finished := range closeChans {
		<-finished
	}
}

func (hn HiveNet) AddWDS(s WDS) {
	netServers := len(hn.servers)
	if netServers > 0 {
		if netServers == 1 {
			s.EstablishLink(hn.servers[0])
		} else {
			s.EstablishLink(hn.servers[0], hn.servers[netServers-1])
		}
	}
	hn.servers = append(hn.servers, s)
}

func (hn HiveNet) NetworkLog() map[string]interface{} {
	netLog := make(map[string]interface{})
	for _, wds := range hn.servers {
		netLog[wds.Id()] = wds.ActivityLog()
	}
	return netLog
}
