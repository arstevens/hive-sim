package reporter

import (
	"os"
	"strconv"

	"github.com/arstevens/hive-sim/src/simulator"
)

type CSVFormatter struct {
	entries []string
	outPath string
}

func NewCSVFormatter(out string) *CSVFormatter {
	return &CSVFormatter{
		outPath: out,
		entries: make([]string, 0),
	}
}

func (cf *CSVFormatter) Format(netlog map[string]simulator.Log) {
	cf.entries = make([]string, len(netlog))
	entryCount := 0

	for wdsId, log := range netlog {
		statsMap := log.GetStats()
		entry := wdsId + ","

		for statType, stat := range statsMap {
			entry += statType + "," + strconv.Itoa(stat) + ","
		}
		entry = entry[:len(entry)-1]
		cf.entries[entryCount] = entry

		entryCount++
	}
}

func (cf CSVFormatter) Save() error {
	file, err := os.Create(cf.outPath)
	if err != nil {
		return err
	}

	entries := cf.entries
	for _, entry := range entries {
		out := []byte(entry + "\n")
		_, err := file.Write(out)
		if err != nil {
			return err
		}
	}
	return nil
}
