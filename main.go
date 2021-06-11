package main

import (
	"encoding/csv"
	"flag"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/Sam-Izdat/govote"
	"github.com/cheggaaa/pb/v3"
)

// A ballot represents a ranking of candidates A, B, and C
type Ballot []string

var numCandidates = flag.Int("candidates", 3, "Sets the number of candidates per race.")
var totalBallots = flag.Int64("total-ballots", int64(10000), "Sets the total number of ballots per election.")
var totalElections = flag.Int64("total-elections", int64(10000), "Sets the total number of elections to run.")
var output = flag.String("output", "results.csv", "The output file for election results.")
var candidates = Ballot{}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func generateCandidates(num int) {
	for i := 0; i < num; i++ {
		candidates = append(candidates, strconv.QuoteRuneToASCII(rune(i+65)))
	}
}

func main() {
	flag.Parse()
	generateCandidates(*numCandidates)

	rand.Seed(time.Now().UnixNano())

	resultsFile, err := os.Create(*output)
	chk(err)
	defer resultsFile.Close()

	resultsData := csv.NewWriter(resultsFile)
	defer resultsData.Flush()

	bar := pb.Start64(*totalElections)

	for j := int64(0); j < *totalElections; j++ {

		plurality, err := govote.Plurality.New(candidates)
		chk(err)
		irv, err := govote.InstantRunoff.New(candidates)
		chk(err)
		condorcet, err := govote.Schulze.New(candidates)
		chk(err)

		for i := int64(0); i < *totalBallots; i++ {
			ballot := generateBallot()
			plurality.AddBallot(ballot[0])
			irv.AddBallot(ballot)
			condorcet.AddBallot(ballot)
		}

		_, pluralityRes, err := plurality.Evaluate()
		chk(err)
		_, irvRes, err := irv.Evaluate()
		chk(err)
		_, condorcetRes, err := condorcet.Evaluate()
		chk(err)

		pluralityWinner := pluralityRes[0].Name
		irvWinner := irvRes[len(irvRes)-1][0].Name
		condorcetResWinner := condorcetRes[0].Name

		resultsData.Write([]string{pluralityWinner, irvWinner, condorcetResWinner})

		bar.Increment()
	}

	bar.Finish()
}

func generateBallot() Ballot {
	ballot := make(Ballot, len(candidates))
	copy(ballot, candidates)
	var swapRankings = func(i, j int) { ballot[i], ballot[j] = ballot[j], ballot[i] }
	rand.Shuffle(len(ballot), swapRankings)
	return ballot
}
