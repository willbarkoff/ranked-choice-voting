package election

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

type Ballot struct {
	Rankings []string
	ID       string
}

type Election struct {
	Candidates      map[string]string
	Ballots         []Ballot
	Title           string
	Contest         string
	BallotSeperator string
	RankSeperator   string
}

type Candidate struct {
	Name     string
	ID       string
	Rankings []int
}

func removePrefix(s, prefix string) string {
	s = strings.TrimSpace(strings.TrimPrefix(s, prefix))
	s = strings.TrimPrefix(s, "\"")
	s = strings.TrimSuffix(s, "\"")
	return s
}

func ElectionFromConfig(filename string) (Election, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Election{}, err
	}

	scanner := bufio.NewScanner(file)

	election := Election{}
	election.Candidates = map[string]string{}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 || strings.HasPrefix(line, "#") {
			// comment
			continue
		} else if strings.HasPrefix(line, ".TITLE") {
			if len(election.Title) != 0 {
				election.Title += "\n"
			}
			election.Title += removePrefix(line, ".TITLE")
		} else if strings.HasPrefix(line, ".CONTEST") {
			election.Contest += removePrefix(line, ".CONTEST")
		} else if strings.HasPrefix(line, ".BALLOT-FORMAT-SEPS") {
			election.RankSeperator = removePrefix(line, ".BALLOT-FORMAT-SEPS")[0:1]
			election.BallotSeperator = removePrefix(line, ".BALLOT-FORMAT-SEPS")[1:2]
		} else if strings.HasPrefix(line, ".CANDIDATE") {
			line = removePrefix(line, ".CANDIDATE")
			data := strings.Split(line, ",")
			election.Candidates[data[0]] = removePrefix(data[1], "")
		} else if strings.HasPrefix(line, ".INCLUDE") {
			filename := removePrefix(line, ".INCLUDE")
			dir := path.Dir(file.Name())
			err = election.Include(path.Join(dir, filename))
			if err != nil {
				return election, err
			}
		} else {
			fmt.Fprintf(os.Stderr, "Unknown directive / ignoring: %s\n", line)
		}
	}

	return election, nil
}

func (e *Election) Include(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)

	election := Election{}
	election.Candidates = map[string]string{}

	for scanner.Scan() {
		ballotText := scanner.Text()
		if len(strings.TrimSpace(ballotText)) == 0 {
			continue
		}
		ballotMeta := strings.Split(ballotText, e.BallotSeperator)
		if len(ballotMeta) < 2 {
			fmt.Fprintf(os.Stderr, "Invalid ballot / ignoring: %v\n", ballotText)
			continue
		}
		rankingsText := strings.Split(strings.TrimSpace(ballotMeta[1]), e.RankSeperator)
		ballot := Ballot{ID: ballotMeta[0], Rankings: []string{}}

		for _, rank := range rankingsText {
			if len(rank)-3 == -3 {
				fmt.Fprintf(os.Stderr, "Invalid ballot / ignoring: %v\n", ballotText)
				continue
			}
			ballot.Rankings = append(ballot.Rankings, rank[:len(rank)-3])
		}

		e.Ballots = append(e.Ballots, ballot)
	}

	return nil
}

func (e *Election) Tally() []Candidate {
	candidates := []Candidate{}

	for candidate, candidateName := range e.Candidates {
		c := Candidate{}
		c.Name = candidateName
		c.ID = candidate
		c.Rankings = make([]int, len(e.Candidates))

		for _, ballot := range e.Ballots {
			rank := find(ballot.Rankings, candidate)
			if rank == -1 {
				continue
			}
			c.Rankings[rank]++
		}

		candidates = append(candidates, c)
	}

	return candidates
}

// Find returns the smallest index i at which x == a[i],
// or -1 if there is no such index.
func find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return -1
}
