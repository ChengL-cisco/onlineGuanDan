package models

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// Deck represents a collection of playing cards
type Deck struct {
	cards []Card
}

// NewDeck creates and returns a new, shuffled deck of cards
func NewDeck(numDecks int) *Deck {
	d := &Deck{}
	d.Initialize(numDecks)
	return d
}

// Initialize creates and returns a new shuffled deck with the specified number of card sets
// Each set contains 54 cards (52 standard + 2 jokers)
func (d *Deck) Initialize(numDecks int) []Card {
	d.cards = make([]Card, 0, numDecks*54)

	for i := 0; i < numDecks; i++ {
		// Add standard cards (2-Ace of each suit)
		for _, suit := range []Suit{Spade, Heart, Diamond, Club} {
			for rank := Two; rank <= Ace; rank++ {
				d.cards = append(d.cards, NewCard(suit, rank))
			}
		}

		// Add jokers
		d.cards = append(d.cards, NewCard("", Joker))    // Small Joker
		d.cards = append(d.cards, NewCard("", BigJoker)) // Big Joker
	}

	// Shuffle the deck
	d.Shuffle()

	return d.cards
}

// Split divides the deck into numPlayers equal parts.
// Returns a slice of card slices, where each inner slice represents one player's cards.
// If the deck can't be evenly divided, some players may receive one more card than others.
// Returns nil if numPlayers is less than 1.
func (d *Deck) Split(numPlayers int) []*Deck {
	if numPlayers <= 0 || len(d.cards) == 0 {
		return nil
	}

	// Calculate number of cards per player
	cardsPerPlayer := len(d.cards) / numPlayers
	extraCards := len(d.cards) % numPlayers

	result := make([]*Deck, 0, numPlayers)
	start := 0

	for i := 0; i < numPlayers; i++ {
		// Calculate how many cards this player gets
		count := cardsPerPlayer
		if i < extraCards {
			count++
		}

		// If we've reached the end of the deck, break early
		if start >= len(d.cards) {
			break
		}

		// Calculate end index, ensuring we don't go past the end of the slice
		end := start + count
		if end > len(d.cards) {
			end = len(d.cards)
		}

		// Add this player's cards to the result
		playerCards := make([]Card, end-start)
		copy(playerCards, d.cards[start:end])
		playerDeck := &Deck{cards: playerCards}
		result = append(result, playerDeck)

		start = end
	}

	return result
}

// Shuffle randomizes the order of cards in the deck
func (d *Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

// Draw removes and returns the top card from the deck
func (d *Deck) Draw() (Card, bool) {
	if len(d.cards) == 0 {
		return Card{}, false
	}
	card := d.cards[0]
	d.cards = d.cards[1:]
	return card, true
}

// DrawN draws n cards from the top of the deck
func (d *Deck) DrawN(n int) []Card {
	if n <= 0 {
		return []Card{}
	}
	if n > len(d.cards) {
		n = len(d.cards)
	}
	cards := d.cards[:n]
	d.cards = d.cards[n:]
	return cards
}

// Count returns the number of cards remaining in the deck
func (d *Deck) Count() int {
	return len(d.cards)
}

// String returns a string representation of all cards in the deck
func (d *Deck) String() string {
	if d.IsEmpty() {
		return "[empty deck]"
	}

	result := ""
	// Construct two lines, first line is index, second line is the card
	indexLine := ""
	cardLine := ""
	for i, card := range d.cards {
		// Calculate padding for index to center it above the card
		indexStr := fmt.Sprintf("%2d", i)

		// Add index with padding
		indexLine += fmt.Sprintf("%-*s", 5, indexStr)
		cardLine += fmt.Sprintf("%-*s", 5, card.String())
	}
	result = indexLine + "\n" + cardLine
	return result
}

// IsEmpty returns true if the deck has no cards
func (d *Deck) IsEmpty() bool {
	return len(d.cards) == 0
}

// Add adds a card to the bottom of the deck
func (d *Deck) Add(card Card) {
	d.cards = append(d.cards, card)
}

// AddToTop adds a card to the top of the deck
func (d *Deck) AddToTop(card Card) {
	d.cards = append([]Card{card}, d.cards...)
}

// PlayIndex removes and returns the card at the specified index.
// The second return value indicates whether the card was successfully removed.
// Returns (zero-value Card, false) if the index is out of bounds.
func (d *Deck) PlayIndex(index int) (Card, bool) {
	if index < 0 || index >= len(d.cards) {
		return Card{}, false
	}

	card := d.cards[index]
	// Remove the card by slicing it out
	d.cards = append(d.cards[:index], d.cards[index+1:]...)
	return card, true
}

// PlayIndexN removes cards at the specified indices from the deck.
// Indices are 0-based and must be valid for the current deck.
// Returns the removed cards in the order of their indices, or nil if any index is invalid.
// The operation is atomic - either all cards are removed or none are.
func (d *Deck) PlayIndexN(indices []int) []Card {
	if len(indices) == 0 {
		return []Card{}
	}

	// Create a map to track which indices to remove
	toRemove := make(map[int]bool)
	maxIndex := len(d.cards) - 1

	// First, validate all indices
	for _, idx := range indices {
		if idx < 0 || idx > maxIndex {
			return nil // Invalid index
		}
		toRemove[idx] = true
	}

	// If we're removing all cards
	if len(toRemove) == len(d.cards) {
		removed := make([]Card, len(d.cards))
		copy(removed, d.cards)
		d.cards = d.cards[:0] // Clear the deck
		return removed
	}

	// Create a new slice for the remaining cards
	newCards := make([]Card, 0, len(d.cards)-len(toRemove))
	removed := make([]Card, 0, len(toRemove))

	// Preserve the order of the indices in the input
	removedMap := make(map[int]Card)
	for idx := range toRemove {
		removedMap[idx] = d.cards[idx]
	}

	// Build the new slice and collect removed cards in order
	for i, card := range d.cards {
		if _, exists := toRemove[i]; exists {
			removed = append(removed, removedMap[i])
		} else {
			newCards = append(newCards, card)
		}
	}

	d.cards = newCards
	return removed
}

// Play removes the specified card from the deck if it exists.
// Returns true if the card was found and removed, false otherwise.
// Note: This performs a linear search through the deck.
func (d *Deck) Play(card Card) bool {
	for i, c := range d.cards {
		if c.Rank == card.Rank && c.Suit == card.Suit {
			// Remove the card by slicing it out
			d.cards = append(d.cards[:i], d.cards[i+1:]...)
			return true
		}
	}
	return false
}

// PlayN removes the specified cards from the deck.
// Returns true if all cards were found and removed, false otherwise.
// If any card is not found, no cards are removed (atomic operation).
func (d *Deck) PlayN(cards []Card) bool {
	if len(cards) == 0 {
		return true
	}

	// First, check if all cards exist in the deck
	cardCounts := make(map[Card]int)
	for _, card := range cards {
		cardCounts[card]++
	}

	// Verify all cards exist with sufficient quantity
	tempCounts := make(map[Card]int)
	for card, needed := range cardCounts {
		found := 0
		for _, c := range d.cards {
			if c.Rank == card.Rank && c.Suit == card.Suit {
				found++
				if found == needed {
					break
				}
			}
		}
		if found < needed {
			return false
		}
		tempCounts[card] = needed
	}

	// If we got here, all cards exist - now remove them
	newCards := make([]Card, 0, len(d.cards)-len(cards))
	removeCounts := make(map[Card]int)
	for _, card := range cards {
		removeCounts[card]++
	}

	for _, card := range d.cards {
		if count, exists := removeCounts[card]; exists && count > 0 {
			removeCounts[card]--
		} else {
			newCards = append(newCards, card)
		}
	}

	d.cards = newCards
	return true
}

// MoveNDCards moves multiple cards specified by their indices to the destination index.
// Returns true if the move was successful, false if any index is out of bounds.
// The moved cards will maintain their relative order.
// If the destination is within the source range, returns false as this would create an invalid state.
func (d *Deck) MoveNDCards(srcIndexes []int, dest int) bool {
	if len(srcIndexes) == 0 {
		return false
	}

	// Validate indices and check for duplicates
	seen := make(map[int]bool)
	minIndex := len(d.cards)
	maxIndex := 0
	for _, idx := range srcIndexes {
		if idx < 0 || idx >= len(d.cards) || seen[idx] {
			return false
		}
		seen[idx] = true
		if idx < minIndex {
			minIndex = idx
		}
		if idx > maxIndex {
			maxIndex = idx
		}
	}

	// Check if destination is within the source range
	if dest >= minIndex && dest <= maxIndex+1 {
		return false
	}

	// Sort the source indices to maintain relative order
	sortedIndexes := make([]int, len(srcIndexes))
	copy(sortedIndexes, srcIndexes)
	sort.Ints(sortedIndexes)

	// Extract the cards to move
	cardsToMove := make([]Card, len(sortedIndexes))
	for i, idx := range sortedIndexes {
		cardsToMove[i] = d.cards[idx]
	}

	// Remove the cards from the source (from highest to lowest to maintain correct indices)
	for i := len(sortedIndexes) - 1; i >= 0; i-- {
		idx := sortedIndexes[i]
		d.cards = append(d.cards[:idx], d.cards[idx+1:]...)
	}

	// Adjust destination index if it was after the source range
	if dest > maxIndex {
		dest -= len(sortedIndexes)
	}

	// Insert the cards at the destination
	d.cards = append(d.cards[:dest], append(cardsToMove, d.cards[dest:]...)...)
	return true
}

// MoveNCards moves a range of cards from start to end (inclusive) to the destination index.
// Returns true if the move was successful, false if any index is out of bounds.
// The moved cards will maintain their relative order.
// If the destination is within the source range, returns false as this would create an invalid state.
func (d *Deck) MoveNCards(start, end, dest int) bool {
	// Validate indices
	if start < 0 || end >= len(d.cards) || dest > len(d.cards) || start > end {
		return false
	}

	// Check if destination is within the source range
	if dest >= start && dest <= end+1 {
		return false
	}

	// Extract the cards to move
	cardsToMove := make([]Card, end-start+1)
	copy(cardsToMove, d.cards[start:end+1])

	// Remove the cards from the source
	d.cards = append(d.cards[:start], d.cards[end+1:]...)

	// Adjust destination index if it was after the source range
	if dest > end {
		dest -= (end - start + 1)
	}

	// Insert the cards at the destination
	d.cards = append(d.cards[:dest], append(cardsToMove, d.cards[dest:]...)...)

	return true
}

// MoveCard moves a card from src index to dest index in the deck.
// Returns true if the move was successful, false if either index is out of bounds.
// If src and dest are the same, returns true without modifying the deck.
// The card at src index is removed and inserted at dest index, shifting other cards as needed.
func (d *Deck) MoveCard(src, dest int) bool {
	if src < 0 || src >= len(d.cards) || dest < 0 || dest > len(d.cards) {
		return false
	}

	if src == dest {
		return true // No-op if source and destination are the same
	}

	// Get the card to move
	card := d.cards[src]

	// Remove the card from the source position
	if src < dest {
		// When moving forward, we need to adjust the destination index
		// because we're removing an element before it
		dest--
	}

	// Create a new slice with the card removed
	d.cards = append(d.cards[:src], d.cards[src+1:]...)

	// Insert the card at the destination
	d.cards = append(d.cards[:dest], append([]Card{card}, d.cards[dest:]...)...)

	return true
}

// GetCards returns a copy of the cards in the deck.
// This ensures the internal slice cannot be modified from outside the package.
func (d *Deck) GetCards() []Card {
	if len(d.cards) == 0 {
		return nil
	}
	cards := make([]Card, len(d.cards))
	copy(cards, d.cards)
	return cards
}

// Sort sorts the cards in the deck with the following order:
// 1. Jokers (Big Joker > Small Joker)
// 2. Cards matching the trump rank (by suit: Spade > Heart > Club > Diamond)
// 3. Other cards (by rank, then by suit)
func (d *Deck) Sort(trumpRank Rank) {
	if len(d.cards) <= 1 {
		return
	}

	// Sort the cards
	sort.Slice(d.cards, func(i, j int) bool {
		a, b := d.cards[i], d.cards[j]

		// Jokers first (Big Joker > Small Joker)
		if a.Rank == BigJoker && b.Rank != BigJoker {
			return true
		}
		if a.Rank != BigJoker && b.Rank == BigJoker {
			return false
		}
		if a.Rank == Joker && b.Rank != Joker {
			return true
		}
		if a.Rank != Joker && b.Rank == Joker {
			return false
		}

		// Then cards matching the trump rank
		aIsTrump := a.Rank == trumpRank
		bIsTrump := b.Rank == trumpRank

		if aIsTrump && !bIsTrump {
			return true
		}
		if !aIsTrump && bIsTrump {
			return false
		}
		if aIsTrump && bIsTrump {
			// Both are trump rank, sort by suit
			return a.Suit > b.Suit
		}

		// For non-trump cards, sort by rank then by suit
		if a.Rank != b.Rank {
			return a.Rank < b.Rank
		}
		suitOrder := map[Suit]int{Spade: 4, Heart: 3, Club: 2, Diamond: 1}
		return suitOrder[a.Suit] > suitOrder[b.Suit]
	})
}
