package main

import (
	"encoding/csv"
	"flag"
	"os"
	"strconv"

	"github.com/willbarkoff/ranked-choice-voting/read-prm/election"
)

var configFile = flag.String("in", "election.in", "The input file to be used for counting")
var outFile = flag.String("out", "out.csv", "The output file to be used for storing tallies")

func main() {
	resultsFile, err := os.Create(*outFile)
	if err != nil {
		panic(err)
	}
	defer resultsFile.Close()

	resultsData := csv.NewWriter(resultsFile)
	defer resultsData.Flush()

	flag.Parse()

	election, err := election.ElectionFromConfig(*configFile)
	if err != nil {
		panic(err)
	}

	counts := election.Tally()

	for _, c := range counts {
		line := []string{c.Name, c.ID}
		for _, tally := range c.Rankings {
			line = append(line, strconv.Itoa(tally))
		}
		resultsData.Write(line)
	}
}
