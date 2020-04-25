package simulator

import "fmt"

func AllocateResources(hn *HiveNet, gen StateGenerator) {
	totalServers := gen.WDSLeft()
	distribution := gen.GetNodeDistribution()

	servers := make([]WDS, totalServers)
	for i, serverLoad := range distribution {
		wds := gen.NextWDS()
		for j := 0; j < serverLoad; j++ {
			newNode := gen.NextNode()
			wds.AssignNode(newNode)
		}
		servers[i] = wds
	}

	for _, server := range servers {
		server.SetMasterKeyList(gen.GetAllNodes())
		hn.AddWDS(server)
	}
	fmt.Println("allocated")

}
