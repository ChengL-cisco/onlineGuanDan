package models

import "fmt"

// Suit represents the suit of a playing card
type Suit string

// Constants for card suits
const (
	Spade   Suit = "♠"
	Heart   Suit = "♥"
	Diamond Suit = "♦"
	Club    Suit = "♣"
)

// Rank represents the rank of a playing card
type Rank int

// Constants for card ranks
const (
	Two Rank = iota + 2
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace

	Joker    // Small Joker
	BigJoker // Big Joker
)

// Card represents a playing card with a suit and rank
type Card struct {
	Suit Suit
	Rank Rank
}

// String returns a string representation of the card
func (c Card) String() string {
	if c.Rank == Joker {
		return "Joker"
	}
	if c.Rank == BigJoker {
		return "Big Joker"
	}

	switch c.Rank {
	case Jack:
		return fmt.Sprintf("J%s", c.Suit)
	case Queen:
		return fmt.Sprintf("Q%s", c.Suit)
	case King:
		return fmt.Sprintf("K%s", c.Suit)
	case Ace:
		return fmt.Sprintf("A%s", c.Suit)
	default:
		return fmt.Sprintf("%d%s", c.Rank, c.Suit)
	}
}

// NewCard creates a new card with the given suit and rank
func NewCard(suit Suit, rank Rank) Card {
	return Card{
		Suit: suit,
		Rank: rank,
	}
}
