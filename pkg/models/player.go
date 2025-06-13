package models

// Player represents a player in the card game
type Player struct {
	// index is the index of the player
	index int

	// name is the display name of the player
	name string

	// hand contains the player's current cards as a Deck
	hand Deck
	// finishedRank is the rank of the player when the round ends
	finishedRank int
	// infoAPI is the interface for game information
	info InfoAPI
}

// GetIndex returns the player's index
func (p *Player) GetIndex() int {
	return p.index
}

// Sit sets the player's index
func (p *Player) Sit(index int) {
	p.index = index
}

// ReadyToStart marks the player as ready to start by updating the readyToStartIndexes in info
func (p *Player) ReadyToStart() {
	if p.info == nil {
		return
	}
	indexes := p.info.GetReadyToStartIndexes()
	if indexes == nil {
		indexes = make([]int, p.info.GetNumPlayers())
	}

	indexes[p.index] = 1
	p.info.SetReadyToStartIndexes(indexes)

}

// ReadyToPlay marks the player as ready to play by updating the readyToPlayIndexes in info
func (p *Player) ReadyToPlay() {
	if p.info == nil {
		return
	}
	indexes := p.info.GetReadyToPlayIndexes()
	if indexes == nil {
		indexes = make([]int, p.info.GetNumPlayers())
	}

	indexes[p.index] = 1
	p.info.SetReadyToPlayIndexes(indexes)
}

// Pass advances the current player index to the next player, wrapping around to 0 if needed
func (p *Player) Pass() {
	if p.info == nil {
		return
	}
	current := p.info.GetCurrentPlayerIndex()
	next := (current + 1) % p.info.GetNumPlayers()
	p.info.SetCurrentPlayerIndex(next)
}

// LeaveGame resets the readyToStartIndexes in info to 0 at the player's index
func (p *Player) LeaveGame() {
	if p.info == nil {
		return
	}
	indexes := p.info.GetReadyToStartIndexes()
	indexes[p.index] = 0
	p.info.SetReadyToStartIndexes(indexes)
}

// GetName returns the player's name
func (p *Player) GetName() string {
	return p.name
}

// SetName sets the player's name
func (p *Player) SetName(name string) {
	p.name = name
}

// GetHand returns the player's hand
func (p *Player) GetHand() Deck {
	return p.hand
}

// SetHand sets the player's hand
func (p *Player) SetHand(hand Deck) {
	p.hand = hand
}

// GetFinishedRank returns the player's finished rank
func (p *Player) GetFinishedRank() int {
	return p.finishedRank
}

// SetFinishedRank sets the player's finished rank
func (p *Player) SetFinishedRank(rank int) {
	p.finishedRank = rank
}
