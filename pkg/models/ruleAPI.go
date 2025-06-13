package models

type RuleAPI interface {
	IsPlayValid(play []Card) bool
	IsCounterPlayValid(play []Card, counterPlay []Card) bool
}
