package entity

import "time"

type difficulty int

const (
	Eazy difficulty = iota
	Medium
	Hard
)

func (d difficulty) IsValid() bool {
	if d <= Eazy && d >= Hard {
		return true
	}

	return false
}

type PossibleAnswerChoice uint8

const (
	PossibleAnswerA PossibleAnswerChoice = iota + 1
	PossibleAnswerB
	PossibleAnswerC
	PossibleAnswerD
)

func (c PossibleAnswerChoice) IsValid() bool {
	if c >= PossibleAnswerA && c <= PossibleAnswerD {
		return true
	}

	return false
}

type Quesiton struct {
	ID                uint
	Content           string
	CategoryID        uint
	Difficulty        difficulty
	PossibleAnswerIDs []uint
	CorrectAnswerID   uint
	TimeToAnswer      time.Duration
}

type PossibleAnswer struct {
	ID      uint
	Choice  PossibleAnswerChoice
	Content string
}
