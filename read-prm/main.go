package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/willbarkoff/ranked-choice-voting/read-prm/election"
)

func main() {
	election, err := election.ElectionFromConfig("../election-2009-burlington-mayoral/2009 Burlington.in")
	if err != nil {
		panic(err)
	}
	spew.Dump(election.Tally())
}
