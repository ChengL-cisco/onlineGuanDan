package models

// DeckAPI defines the public interface for interacting with a deck of cards
type DeckAPI interface {
	// Initialize returns a slice of a shuffled numDecks cards
	Initialize(numDecks int) []Card

	// Split splits the deck into numPlayers equal parts
	Split(numPlayers int) [][]Card

	// GetCards returns a copy of all cards in the deck
	GetCards() []Card

	// Count returns the number of cards remaining in the deck
	Count() int

	// IsEmpty returns true if the deck has no cards
	IsEmpty() bool

	// Add adds a card to the bottom of the deck
	Add(card Card)

	// AddToTop adds a card to the top of the deck
	AddToTop(card Card)

	// Play removes the specified card from the deck if it exists
	Play(card Card) bool

	// PlayN removes the specified cards from the deck
	PlayN(cards []Card) bool

	// PlayIndex removes and returns the card at the specified index
	PlayIndex(index int) (Card, bool)

	// PlayIndexN removes cards at the specified indices from the deck
	PlayIndexN(indices []int) []Card

	// MoveCard moves a card from src index to dest index in the deck
	MoveCard(src, dest int) bool

	// MoveNCards moves a range of cards from start to end (inclusive) to the destination index
	MoveNCards(start, end, dest int) bool

	// Sort sorts the cards in the deck with jokers first, then cards matching the trump card's rank,
	// then other cards by rank
	// todo: use the sort algorithm provided by Rule
	Sort(trumpCard Card)

	// String returns a string representation of all cards in the deck
	String() string
}

// Verify at compile time that *Deck implements DeckAPI
var _ DeckAPI = (*Deck)(nil)
