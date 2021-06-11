package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/davecgh/go-spew/spew"
)

// A ballot represents a ranking of candidates A, B, and C
type Ballot []string

var rankings = []string{"a", "b", "c"}
var swapRankings = func(i, j int) { rankings[i], rankings[j] = rankings[j], rankings[i] }
var totalBallots = 10000

func generateBallot() Ballot {
	rand.Shuffle(len(rankings), swapRankings)
	ballot := []string{"", "", ""}
	copy(ballot, rankings)
	return ballot
}

func min(a, b, c int) string {
	if a < b && a < c {
		return "a"
	} else if b < a && b < c {
		return "b"
	} else if c < a && c < b {
		return "c"
	} else {
		return fmt.Sprintf("tie a: %d b: %d c: %d", a, b, c)
	}
}

func count(ballots []Ballot, position int) (int, int, int) {
	a, b, c := 0, 0, 0
	for _, ballot := range ballots {
		if ballot[position] == "a" {
			a++
		} else if ballot[position] == "b" {
			b++
		} else if ballot[position] == "c" {
			c++
		} else {
			panic("invalid ballot: " + spew.Sprint(ballot))
		}
	}
	return a, b, c
}

func main() {
	rand.Seed(time.Now().UnixNano())

	ballots := []Ballot{}

	for i := 0; i < totalBallots; i++ {
		ballot := generateBallot()
		ballots = append(ballots, ballot)
	}

	fmt.Println(traditional(ballots))
	fmt.Println(irv(ballots))
}

func traditional(ballots []Ballot) string {
	a, b, c := count(ballots, 0)

	if a > b && a > c {
		return "a"
	} else if b > a && b > c {
		return "b"
	} else if c > a && c > b {
		return "c"
	} else {
		return fmt.Sprintf("tie a: %d b: %d c: %d", a, b, c)
	}
}

func irv(ballots []Ballot) string {
	a, b, c := count(ballots, 0)

	threshold := len(ballots) / 2

	if a > threshold {
		return "a"
	} else if b > threshold {
		return "b"
	} else if c > threshold {
		return "c"
	}

	// there's no winner, let's throw away old votes
	min := min(a, b, c)

	if len(min) != 1 {
		fmt.Println(min)
		return min
	}

	newBallots := []Ballot{}
	for _, ballot := range ballots {
		if ballot[0] == min {
			ballot = ballot[1:]
		}
		newBallots = append(newBallots, ballot)
	}

	return irv(newBallots)
}
