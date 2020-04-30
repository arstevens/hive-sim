package reporter

import (
	"fmt"

	"github.com/arstevens/hive-sim/src/simulator"
)

type ConsoleFormatter struct {
	entries []string
}

func (cf *ConsoleFormatter) Format(netlog map[string]simulator.Log) {
	cf.entries = make([]string, len(netlog))
	entryCount := 0

	for wdsId, log := range netlog {
		statsMap := log.GetStats()
		entryFmt := "Id: " + wdsId + " {\n"
		stats := make([]int, len(statsMap))
		statCount := 0

		for statType, stat := range statsMap {
			entryFmt += statType + ": %d\n"
			stats[statCount] = stat
			statCount++
		}
		entryFmt += "}\n"

		entry := fmt.Sprintf(entryFmt, stats...)
		cf.entries[entryCount] = entry
		entryCount++
	}
}

func (cf ConsoleFormatter) Save() error {
	for _, entry := range cf.entries {
		fmt.Print(entry)
	}
	return nil
}
