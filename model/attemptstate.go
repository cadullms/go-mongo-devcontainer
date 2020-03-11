package model

const (
	AttemptStateAsk = 0
    AttemptStateCorrect = 1
    AttemptStateIncorrect = 2
    AttemptStatePresent = 3
)

type AttemptState int32