package main

import (
	"testing"

	"github.com/google/uuid"
)

func TestCreateWithOptions(t *testing.T) {
	poll := NewPoll("mypoll", []string{"a", "b", "c"})
	numOptions := len(poll.Options)
	expectedNum := 3
	if numOptions != expectedNum {
		t.Fatalf(`len(poll.Options) = %q, expected %#q`, numOptions, expectedNum)
	}
}

// Returns a new UUID, panicing if any errors occur
func NewID() uuid.UUID {
	id, err := uuid.NewUUID()
	if err != nil {
		panic("error found while generating UUID")
	}
	return id
}

func AddUniqueVotes(poll Poll, optionName string, numVotes int) {
	for range numVotes {
		id := NewID()
		poll.AddVote(id, optionName)
	}
}

func TestAddVotes(t *testing.T) {
	opt1 := "option 1"
	opt2 := "option 2"
	opt3 := "option 3"

	poll := NewPoll("mypoll", []string{opt1, opt2, opt3})

	for name, option := range poll.Options {
		if len(option.Votes) != 0 {
			t.Fatalf(`len(option.Votes) = %q for option %q, expected %#q`, len(option.Votes), name, 0)
		}
	}

	AddUniqueVotes(*poll, opt1, 2)
	AddUniqueVotes(*poll, opt2, 3)
	AddUniqueVotes(*poll, opt3, 1)

	opt1Votes, _ := poll.NumVotesFor(opt1)
	if opt1Votes != 2 {
		t.Fatalf(`poll.NumVotesFor(opt1) = %v, expected %v`, opt1Votes, 2)
	}

	opt2Votes, _ := poll.NumVotesFor(opt2)
	if opt2Votes != 3 {
		t.Fatalf(`poll.NumVotesFor(opt2) = %v, expected %v`, opt2Votes, 3)
	}

	opt3Votes, _ := poll.NumVotesFor(opt3)
	if opt3Votes != 1 {
		t.Fatalf(`poll.NumVotesFor(opt3) = %v, expected %v`, opt3Votes, 1)
	}
}
