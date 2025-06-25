package models

import (
	"fmt"
	"strings"
)

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

// SuitToInitial returns the single-letter string representation of a suit
// S for Spade, H for Heart, D for Diamond, and C for Club
func SuitToInitial(s Suit) string {
	switch s {
	case Spade:
		return "S"
	case Heart:
		return "H"
	case Diamond:
		return "D"
	case Club:
		return "C"
	default:
		return "?"
	}
}

// String returns a string representation of the card
func (c Card) String() string {
	if c.Rank == Joker {
		return "Jr"
	}
	if c.Rank == BigJoker {
		return "BJr"
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

func (c Card) CardString() string {
	if c.Rank == Joker {
		return "Jr"
	}
	if c.Rank == BigJoker {
		return "BJr"
	}

	switch c.Rank {
	case Jack:
		return fmt.Sprintf("J-%s", SuitToInitial(c.Suit))
	case Queen:
		return fmt.Sprintf("Q-%s", SuitToInitial(c.Suit))
	case King:
		return fmt.Sprintf("K-%s", SuitToInitial(c.Suit))
	case Ace:
		return fmt.Sprintf("A-%s", SuitToInitial(c.Suit))
	default:
		return fmt.Sprintf("%d-%s", c.Rank, SuitToInitial(c.Suit))
	}
}

// CardsString returns a formatted string representation of a slice of cards
// The output shows each card's string representation separated by spaces
// Example output: "2-S 3-H K-D A-C BJr"
func CardsString(cards []Card) string {
	if len(cards) == 0 {
		return ""
	}

	var result strings.Builder
	for i, card := range cards {
		if i > 0 {
			result.WriteString(" ") // One space between cards
		}
		result.WriteString(card.CardString())
	}

	return result.String()
}

// parseCard parses a card string into a Card struct
// Supported formats: "2-S", "J-H", "Q-D", "K-C", "A-S", "Jr", "BJr"
func parseCard(cardStr string) (Card, error) {
	// Handle jokers
	switch cardStr {
	case "Jr":
		return Card{Rank: Joker}, nil
	case "BJr":
		return Card{Rank: BigJoker}, nil
	}

	// Parse rank and suit
	parts := strings.Split(cardStr, "-")
	if len(parts) != 2 {
		return Card{}, fmt.Errorf("invalid card format: %s", cardStr)
	}

	rankStr, suitStr := parts[0], parts[1]

	// Parse rank
	var rank Rank
	switch rankStr {
	case "J":
		rank = Jack
	case "Q":
		rank = Queen
	case "K":
		rank = King
	case "A":
		rank = Ace
	default:
		_, err := fmt.Sscanf(rankStr, "%d", &rank)
		if err != nil || rank < Two || rank > Ten {
			return Card{}, fmt.Errorf("invalid rank: %s", rankStr)
		}
	}

	// Parse suit
	var suit Suit
	switch suitStr {
	case "S":
		suit = Spade
	case "H":
		suit = Heart
	case "D":
		suit = Diamond
	case "C":
		suit = Club
	default:
		return Card{}, fmt.Errorf("invalid suit: %s", suitStr)
	}

	return NewCard(suit, rank), nil
}

// NewDeckFromString creates a new deck from a space-separated string of cards
// Example: "2-S 3-H K-D A-C BJr"
func NewDeckFromString(cardsStr string) (*Deck, error) {
	if cardsStr == "" {
		return &Deck{cards: []Card{}}, nil
	}

	cardStrs := strings.Fields(cardsStr)
	cards := make([]Card, 0, len(cardStrs))

	for _, cardStr := range cardStrs {
		card, err := parseCard(cardStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse card '%s': %w", cardStr, err)
		}
		cards = append(cards, card)
	}

	return &Deck{cards: cards}, nil
}

// NewCard creates a new card with the given suit and rank
func NewCard(suit Suit, rank Rank) Card {
	return Card{
		Suit: suit,
		Rank: rank,
	}
}
