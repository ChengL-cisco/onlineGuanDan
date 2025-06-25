package models

import "sort"

// NumOfDecks calculates the number of decks needed based on the number of players.
func NumOfDecks(numOfPlayers int) int {
	if numOfPlayers <= 0 {
		return 0
	}
	return int(float64(numOfPlayers) * 0.5)
}

// Rule represents the game rules and logic
// This struct will contain methods that define how the game is played
// and how game state transitions occur
type Rule struct {
	info InfoAPI
}

// isPlayValid validates a card play according to the game rules
// It returns true for the following cases:
// 1. Single card
// 2. Pair of cards with the same rank
// 3. Three of a kind
// 4. Four of a kind
// 5. Five cards: all same rank, or full house (3+2), or straight (5 consecutive ranks)
// 6. Six cards: all same rank, or three pairs with consecutive ranks, or two triplets with consecutive ranks
// 7. Seven or more cards of the same rank
func (r *Rule) IsPlayValid(play []Card) bool {
	switch len(play) {
	case 0:
		return false
	case 1:
		return true // Single card is always valid
	case 2:
		// Must be a pair (same rank)
		return play[0].Rank == play[1].Rank
	case 3:
		// Must be three of a kind
		return play[0].Rank == play[1].Rank && play[1].Rank == play[2].Rank
	case 4:
		// Must be four of a kind
		return play[0].Rank == play[1].Rank &&
			play[1].Rank == play[2].Rank &&
			play[2].Rank == play[3].Rank
	case 5:
		return r.isValidFiveCardPlay(play)
	case 6:
		return r.isValidSixCardPlay(play)
	default: // 7 or more cards
		// Must all be the same rank
		return r.allSameRank(play)
	}
}

func (r *Rule) IsCounterPlayValid(play []Card, counterPlay []Card) bool {
	if len(play) == 0 {
		return false // No play to counter
	} else if len(play) == 1 {
		// 1. counterPlay is a single card, and its rank is greater than play
		if len(counterPlay) == 1 && r.IsRankGreater(counterPlay[0].Rank, play[0].Rank) {
			return true
		}

		// 2. counterPlay has 4 or more cards and all the cards are of the same rank
		if len(counterPlay) >= 4 && r.allSameRank(counterPlay) {
			return true
		}

		// 3. counterPlay is a straight flush
		if r.isStraightFlush(counterPlay) {
			return true
		}

		// Otherwise, not a valid counter play
		return false
	} else if len(play) == 2 {
		// 1. counterPlay is a pair with a higher rank
		if len(counterPlay) == 2 && counterPlay[0].Rank == counterPlay[1].Rank && r.IsRankGreater(counterPlay[0].Rank, play[0].Rank) {
			return true
		}

		// 2. counterPlay has 4 or more cards and all the cards are of the same rank
		if len(counterPlay) >= 4 && r.allSameRank(counterPlay) {
			return true
		}

		// 3. counterPlay is a straight flush
		if r.isStraightFlush(counterPlay) {
			return true
		}

		// Otherwise, not a valid counter play
		return false
	} else if len(play) == 3 {
		// 1. counterPlay is a three of a kind with a higher rank
		if len(counterPlay) == 3 && counterPlay[0].Rank == counterPlay[1].Rank && counterPlay[1].Rank == counterPlay[2].Rank && r.IsRankGreater(counterPlay[0].Rank, play[0].Rank) {
			return true
		}

		// 2. counterPlay has 4 or more cards and all the cards are of the same rank
		if len(counterPlay) >= 4 && r.allSameRank(counterPlay) {
			return true
		}

		// 3. counterPlay is a straight flush
		if r.isStraightFlush(counterPlay) {
			return true
		}

		// Otherwise, not a valid counter play
		return false
	} else if len(play) == 4 {
		// 1. counterPlay is a four of a kind with a higher rank
		if len(counterPlay) == 4 && counterPlay[0].Rank == counterPlay[1].Rank &&
			counterPlay[1].Rank == counterPlay[2].Rank &&
			counterPlay[2].Rank == counterPlay[3].Rank &&
			r.IsRankGreater(counterPlay[0].Rank, play[0].Rank) {
			return true
		}

		// 2. counterPlay has 5 or more cards and all the cards are of the same rank
		if len(counterPlay) >= 5 && r.allSameRank(counterPlay) {
			return true
		}

		// 3. counterPlay is a straight flush
		if r.isStraightFlush(counterPlay) {
			return true
		}

		// Otherwise, not a valid counter play
		return false
	} else if len(play) == 5 {
		// If play is a five of a kind, return true if counterPlay is a straight flush or counter play is also a five of a kind but with a higher rank
		if r.allSameRank(play) {
			if r.isStraightFlush(counterPlay) {
				return true
			}
			if r.allSameRank(counterPlay) && r.IsRankGreater(counterPlay[0].Rank, play[0].Rank) {
				return true
			}
		} else if r.isStraightFlush(play) {
			// If play is a straight flush, return true if counterPlay is also a straight flush but the ending rank is higher
			if r.isStraightFlush(counterPlay) && r.IsRankGreater(counterPlay[len(counterPlay)-1].Rank, play[len(play)-1].Rank) {
				return true
			}
		}

		// Otherwise, not a valid counter play
		return false
	} else if len(play) == 6 {
		// If play is six of a kind
		if r.allSameRank(play) {
			// 1. counterPlay is also a six of a kind but with a higher rank
			if r.allSameRank(counterPlay) && len(counterPlay) == 6 && r.IsRankGreater(counterPlay[0].Rank, play[0].Rank) {
				return true
			}
			// 2. counterPlay is a 7 or more of a kind
			if r.allSameRank(counterPlay) && len(counterPlay) >= 7 {
				return true
			}
		} else if rankCount := r.countRanks(play); len(rankCount) == 3 {
			// play is three pairs with consecutive ranks
			isThreePairs := true
			for _, count := range rankCount {
				if count != 2 {
					isThreePairs = false
					break
				}
			}
			if isThreePairs && r.areRanksConsecutive(getSortedRanks(rankCount)) {
				// counterPlay must also be three pairs with consecutive ranks and higher ending rank
				if len(counterPlay) == 6 {
					counterRankCount := r.countRanks(counterPlay)
					isCounterThreePairs := true
					for _, count := range counterRankCount {
						if count != 2 {
							isCounterThreePairs = false
							break
						}
					}
					if isCounterThreePairs && r.areRanksConsecutive(getSortedRanks(counterRankCount)) {
						// Compare the highest rank in both plays
						playRanks := getSortedRanks(rankCount)
						counterRanks := getSortedRanks(counterRankCount)
						if r.IsRankGreater(counterRanks[2], playRanks[2]) {
							return true
						}
					}
				}
				// counterPlay is a bomb
				if r.allSameRank(counterPlay) && len(counterPlay) >= 4 {
					return true
				}
			}
		} else if len(rankCount) == 2 {
			// play is two triplets with consecutive ranks
			isTwoTriplets := true
			for _, count := range rankCount {
				if count != 3 {
					isTwoTriplets = false
					break
				}
			}
			if isTwoTriplets && r.areRanksConsecutive(getSortedRanks(rankCount)) {
				// counterPlay must also be two triplets with consecutive ranks and higher ending rank
				if len(counterPlay) == 6 {
					counterRankCount := r.countRanks(counterPlay)
					isCounterTwoTriplets := true
					for _, count := range counterRankCount {
						if count != 3 {
							isCounterTwoTriplets = false
							break
						}
					}
					if isCounterTwoTriplets && r.areRanksConsecutive(getSortedRanks(counterRankCount)) {
						playRanks := getSortedRanks(rankCount)
						counterRanks := getSortedRanks(counterRankCount)
						if r.IsRankGreater(counterRanks[1], playRanks[1]) {
							return true
						}
					}
				}
				// counterPlay is a bomb
				if r.allSameRank(counterPlay) && len(counterPlay) >= 4 {
					return true
				}
			}
		}
		return false
	} else {
		if r.allSameRank(counterPlay) {
			if len(counterPlay) == len(play) && r.IsRankGreater(counterPlay[0].Rank, play[0].Rank) {
				return true
			} else if len(counterPlay) > len(play) {
				return true
			}
		}
		return false
	}
}

// IsRankGreater checks if rank1 is greater than rank2
// Returns true if rank1 is greater than rank2
func (r *Rule) IsRankGreater(rank1 Rank, rank2 Rank) bool {
	if rank1 == rank2 {
		return false
	}
	trump := r.info.GetTrumpRank()
	if rank1 != trump && rank2 != trump {
		return rank1 > rank2
	}

	if rank1 == trump {
		if rank2 == Joker || rank2 == BigJoker {
			return false
		}
		return true
	} else {
		if rank1 == Joker || rank1 == BigJoker {
			return true
		}
		return false
	}
}

// isValidFiveCardPlay checks if a 5-card play is valid
func (r *Rule) isValidFiveCardPlay(cards []Card) bool {
	// Check for five of a kind
	if r.allSameRank(cards) {
		return true
	}

	// Check for full house (3+2)
	rankCount := r.countRanks(cards)
	if len(rankCount) == 2 {
		for _, count := range rankCount {
			if count == 2 || count == 3 {
				continue
			}
			return false
		}
		return true
	}

	// Check for straight (5 consecutive ranks)
	return r.isStraight(cards)
}

// isValidSixCardPlay checks if a 6-card play is valid
func (r *Rule) isValidSixCardPlay(cards []Card) bool {
	// Check for six of a kind
	if r.allSameRank(cards) {
		return true
	}

	rankCount := r.countRanks(cards)

	// Check for three pairs with consecutive ranks
	if len(rankCount) == 3 {
		for _, count := range rankCount {
			if count != 2 {
				return false
			}
		}
		return r.areRanksConsecutive(getSortedRanks(rankCount))
	}

	// Check for two triplets with consecutive ranks
	if len(rankCount) == 2 {
		for _, count := range rankCount {
			if count != 3 {
				return false
			}
		}
		return r.areRanksConsecutive(getSortedRanks(rankCount))
	}

	return false
}

// allSameRank checks if all cards have the same rank
func (r *Rule) allSameRank(cards []Card) bool {
	if len(cards) == 0 {
		return false
	}
	firstRank := cards[0].Rank
	for _, card := range cards[1:] {
		if card.Rank != firstRank {
			return false
		}
	}
	return true
}

// countRanks returns a map of rank to count
func (r *Rule) countRanks(cards []Card) map[Rank]int {
	rankCount := make(map[Rank]int)
	for _, card := range cards {
		rankCount[card.Rank]++
	}
	return rankCount
}

// isStraight checks if cards form a straight (consecutive ranks)
func (r *Rule) isStraight(cards []Card) bool {
	if len(cards) < 2 {
		return false
	}

	sorted := make([]Card, len(cards))
	copy(sorted, cards)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Rank < sorted[j].Rank
	})

	// Special case for Ace-5 straight (A-2-3-4-5)
	if len(sorted) == 5 && sorted[0].Rank == Two && sorted[1].Rank == Three &&
		sorted[2].Rank == Four && sorted[3].Rank == Five && sorted[4].Rank == Ace {
		return true
	}

	// Check normal straight
	for i := 1; i < len(sorted); i++ {
		if int(sorted[i].Rank)-int(sorted[i-1].Rank) != 1 {
			return false
		}
	}
	return true
}

// isStraightFlush checks if the cards form a straight flush (5 consecutive ranks of the same suit)
func (r *Rule) isStraightFlush(cards []Card) bool {
	if len(cards) != 5 {
		return false
	}

	// Check if all cards are of the same suit
	firstSuit := cards[0].Suit
	for _, card := range cards[1:] {
		if card.Suit != firstSuit {
			return false
		}
	}

	// Check if they form a straight
	return r.isStraight(cards)
}

// areRanksConsecutive checks if the ranks are consecutive
func (r *Rule) areRanksConsecutive(ranks []Rank) bool {
	if len(ranks) < 2 {
		return true
	}

	sort.Slice(ranks, func(i, j int) bool {
		return ranks[i] < ranks[j]
	})

	// Special case for Ace-2-3 straight
	if ranks[0] == Two && ranks[1] == Three && ranks[2] == Four &&
		ranks[3] == Five && ranks[4] == Six && ranks[5] == Ace {
		return true
	}

	for i := 1; i < len(ranks); i++ {
		if int(ranks[i])-int(ranks[i-1]) != 1 {
			return false
		}
	}
	return true
}

// getSortedRanks returns a sorted slice of unique ranks from the rank count map
func getSortedRanks(rankCount map[Rank]int) []Rank {
	var ranks []Rank
	for rank := range rankCount {
		ranks = append(ranks, rank)
	}
	sort.Slice(ranks, func(i, j int) bool {
		return ranks[i] < ranks[j]
	})
	return ranks
}
