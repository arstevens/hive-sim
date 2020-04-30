/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"math"

	"github.com/arstevens/hive-sim/src/interactor"
	"github.com/arstevens/hive-sim/src/simulator"
	"github.com/spf13/cobra"
)

// runbasicCmd represents the runbasic command
var runbasicCmd = &cobra.Command{
	Use:   "basic",
	Short: "basic runs a perfect-world hive simulation",
	Long: `basic runs a simulation where all nodes have a standardized
	behavior. This is unlikely to happen in the real world.`,
	Run: func(cmd *cobra.Command, args []string) {
		serverCount, _ := cmd.Flags().GetInt32("servers")
		nodeCount, _ := cmd.Flags().GetInt32("nodes")
		contractCount, _ := cmd.Flags().GetInt32("contracts")
		transactionLimit, _ := cmd.Flags().GetInt32("tlimit")
		runBasicSimulation(int(serverCount), int(nodeCount),
			int(contractCount), int(transactionLimit))
	},
}

func init() {
	runCmd.AddCommand(runbasicCmd)
	runbasicCmd.Flags().Int32P("servers", "s", 1, "Number of WDS servers to simulate")
	runbasicCmd.Flags().Int32P("nodes", "n", 10, "Number of nodes per server to simulate")
	runbasicCmd.Flags().Int32P("contracts", "c", 0, "Number of contracts per server to simulate")
	runbasicCmd.Flags().Int32P("tlimit", "t", 0, "Exchange limit on a single contract")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runbasicCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runbasicCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runBasicSimulation(serverCount int, nodeCount int, conCount int, transLimit int) {
	sCount, nCount, cCount, tLimit := basicSanitizer(serverCount, nodeCount, conCount, transLimit)
	nodeDist := makeEvenNodeDist(sCount, nCount)

	basicGenerator := interactor.NewBasicGenerator(sCount, nCount, cCount, tLimit, nodeDist)
	networkSimulator := simulator.NewHiveNet()

	simulator.AllocateResources(networkSimulator, basicGenerator)
	networkSimulator.Run()

	logs := networkSimulator.NetworkLog()
	for _, log := range logs {
		log.Print()
	}
}

var (
	serverLimit   = 0xf4240
	nodePerMin    = 10
	nodePerLimit  = 0xf4240
	contractLimit = 0xf4240
)

func basicSanitizer(serverCount int, nodeCount int, conCount int, transLimit int) (int, int, int, int) {
	sCount := basicMinMaxIntSanitize(0, serverLimit, serverCount)
	nCount := basicMinMaxIntSanitize(nodePerMin, nodePerLimit, nodeCount)
	cCount := basicMinMaxIntSanitize(0, contractLimit, conCount)
	tLimit := basicMinMaxIntSanitize(0, math.MaxInt32, transLimit)
	return sCount, nCount, cCount, tLimit
}

func basicMinMaxIntSanitize(min int, max int, val int) int {
	cleanCount := val
	if cleanCount < min {
		cleanCount = min
	}
	if cleanCount > max {
		cleanCount = max
	}
	return cleanCount
}

func makeEvenNodeDist(wdsSize int, nodeSize int) []int {
	dist := make([]int, wdsSize)
	for i := 0; i < len(dist); i++ {
		dist[i] = nodeSize
	}
	return dist
}
