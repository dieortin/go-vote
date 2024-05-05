package main

import (
	"fmt"

	"github.com/google/uuid"
)

type Poll struct {
	Title   string
	Options map[string]*PollOption
}

func NewPoll(title string, options []string) *Poll {
	poll := Poll{
		Title:   title,
		Options: make(map[string]*PollOption),
	}

	for _, v := range options {
		poll.Options[v] = NewPollOption(v)
	}

	return &poll
}

func (p Poll) HasVoteFromUser(voter uuid.UUID) bool {
	for _, option := range p.Options {
		if option.HasVoteFromUser(voter) {
			return true
		}
	}
	return false
}

func (p Poll) NumVotesFor(optionName string) (int, error) {
	option, ok := p.Options[optionName]
	if !ok {
		return 0, fmt.Errorf("option %q does not exist in this poll", optionName)
	}
	return option.NumVotes(), nil
}

func (p Poll) AddVote(voter uuid.UUID, optionName string) error {
	if p.HasVoteFromUser(voter) {
		return fmt.Errorf("user with uuid %v already voted in this poll", voter)
	}
	option, ok := p.Options[optionName]
	if !ok {
		return fmt.Errorf("option %s does not exist", optionName)
	}

	option.AddVote(voter)
	return nil
}

type PollOption struct {
	Name  string
	Votes map[uuid.UUID]bool
}

func NewPollOption(name string) *PollOption {
	return &PollOption{
		Name:  name,
		Votes: make(map[uuid.UUID]bool),
	}
}

func (po PollOption) HasVoteFromUser(voter uuid.UUID) bool {
	_, ok := po.Votes[voter]
	return ok
}

func (po PollOption) AddVote(voter uuid.UUID) {
	po.Votes[voter] = true
}

func (po PollOption) NumVotes() int {
	return len(po.Votes)
}

type PollStorage map[uuid.UUID]*Poll

func NewPollStorage() PollStorage {
	return make(PollStorage)
}

func (ps PollStorage) AddPoll(poll *Poll) (uuid.UUID, error) {
	newUuid, err := uuid.NewUUID()
	if err != nil {
		return uuid.Nil, err
	}

	ps[newUuid] = poll
	return newUuid, nil
}

func (ps PollStorage) GetPoll(id uuid.UUID) *Poll {
	poll, ok := ps[id]
	if !ok {
		return nil
	}
	return poll
}
