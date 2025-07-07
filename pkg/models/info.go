package models

// Info is a placeholder struct for game information
type Info struct {
	// numPlayers is the number of players in the game
	numPlayers int
	// grp1Name is the name of group 1
	grp1Name string
	// grp2Name is the name of group 2
	grp2Name string
	// readyToStart is a map of player indexes to their ready status
	readyToStart map[int]bool
	// readyToPlay is a map of player indexes to their ready to play status
	readyToPlay map[int]bool
	// isFirstRound is true if it is the first round of the game
	isFirstRound bool
	// isRoundInSession is true if a round is currently in session
	isRoundInSession bool
	// currentPlayerIndex is the index of the current player
	currentPlayerIndex int
	// trumpRank is the trump rank for the current round
	trumpRank Rank
	// grpScores is the score of each group
	grpScores [2]int
	// firstFinishedIndex is the index of the first player to finish a round
	firstFinishedIndex int
	// secondFinishedIndex is the index of the second player to finish a round
	secondFinishedIndex int
	// lastFinishedIndex is the index of the last player to finish a round
	lastFinishedIndex int
	// secondToLastFinishedIndex is the index of the second to last player to finish a round
	secondToLastFinishedIndex int
	// lastPlayedCards is the list of cards played by the last player
	lastPlayedCards []Card
	// lastPlayedIndex is the index of the last player to play
	lastPlayedIndex int
	// available slots for players
	availableSlots map[int]bool
	// names of players
	names map[int]string
	// finishedIndexes is a list of indexes of players who have finished a round
	finishedIndexes []int
}

// GetLastPlayedIndex returns the index of the last player to play
func (i *Info) GetLastPlayedIndex() int {
	return i.lastPlayedIndex
}

// SetLastPlayedIndex sets the index of the last player to play
func (i *Info) SetLastPlayedIndex(index int) {
	i.lastPlayedIndex = index
}

// GetFinishedIndexes returns the list of player indexes who have finished the round
func (i *Info) GetFinishedIndexes() []int {
	if i.finishedIndexes == nil {
		i.finishedIndexes = make([]int, 0)
	}
	return i.finishedIndexes
}

// SetFinishedIndexes sets the list of player indexes who have finished the round
func (i *Info) SetFinishedIndexes(indexes []int) {
	i.finishedIndexes = indexes
}

// ResetFinishedIndexes clears the list of finished player indexes
func (i *Info) ResetFinishedIndexes() {
	i.finishedIndexes = nil
}

// GetReadyToStartMap returns a reference to the readyToStart map
func (i *Info) GetReadyToStartMap() map[int]bool {
	if i.readyToStart == nil {
		i.readyToStart = make(map[int]bool)
	}
	return i.readyToStart
}

// SetReadyToStartMap sets the readyToStart map reference
func (i *Info) SetReadyToStartMap(ready map[int]bool) {
	i.readyToStart = ready
}

// GetAvailableSlots returns a reference to the availableSlots map
func (i *Info) GetAvailableSlots() map[int]bool {
	if i.availableSlots == nil {
		i.availableSlots = make(map[int]bool)
	}
	return i.availableSlots
}

// SetAvailableSlots sets the availableSlots map reference
func (i *Info) SetAvailableSlots(slots map[int]bool) {
	i.availableSlots = slots
}

// GetNames returns a reference to the names map
func (i *Info) GetNames() map[int]string {
	if i.names == nil {
		i.names = make(map[int]string)
	}
	return i.names
}

// SetNames sets the names map reference
func (i *Info) SetNames(names map[int]string) {
	i.names = names
}

// GetNumPlayers returns the number of players in the game
func (i *Info) GetNumPlayers() int {
	return i.numPlayers
}

// SetNumPlayers sets the number of players in the game
func (i *Info) SetNumPlayers(numPlayers int) {
	i.numPlayers = numPlayers
}

// GetGrp1Name returns the name of group 1
// Returns "Group1" if no name has been set
func (i *Info) GetGrp1Name() string {
	if i.grp1Name == "" {
		return "Group1"
	}
	return i.grp1Name
}

// SetGrp1Name sets the name of group 1
func (i *Info) SetGrp1Name(name string) {
	i.grp1Name = name
}

// GetGrp2Name returns the name of group 2
// Returns "Group2" if no name has been set
func (i *Info) GetGrp2Name() string {
	if i.grp2Name == "" {
		return "Group2"
	}
	return i.grp2Name
}

// SetGrp2Name sets the name of group 2
func (i *Info) SetGrp2Name(name string) {
	i.grp2Name = name
}

// SetReadyToPlay sets the map of player indexes to their ready to play status
func (i *Info) SetReadyToPlay(readyMap map[int]bool) {
	i.readyToPlay = readyMap
}

// GetReadyToPlay returns the map of player indexes to their ready to play status
func (i *Info) GetReadyToPlay() map[int]bool {
	if i.readyToPlay == nil {
		i.readyToPlay = make(map[int]bool)
	}
	return i.readyToPlay
}

// AddReadyToPlay sets a player's ready to play status to true
func (i *Info) AddReadyToPlay(index int) {
	if i.readyToPlay == nil {
		i.readyToPlay = make(map[int]bool)
	}
	i.readyToPlay[index] = true
}

// RemoveReadyToPlay removes a player's ready to play status
func (i *Info) RemoveReadyToPlay(index int) {
	delete(i.readyToPlay, index)
}

// IsReadyToPlay checks if a player is ready to play
func (i *Info) IsReadyToPlay(index int) bool {
	return i.readyToPlay != nil && i.readyToPlay[index]
}

// AllPlayersReadyToPlay checks if all players are ready to play
func (i *Info) AllPlayersReadyToPlay() bool {
	if i.readyToPlay == nil {
		return false
	}
	return len(i.readyToPlay) == i.numPlayers
}

// GetIsFirstRound returns whether it's the first round of the game
func (i *Info) GetIsFirstRound() bool {
	return i.isFirstRound
}

// SetIsFirstRound sets whether it's the first round of the game
func (i *Info) SetIsFirstRound(isFirstRound bool) {
	i.isFirstRound = isFirstRound
}

// GetIsRoundInSession returns whether a round is currently in session
func (i *Info) GetIsRoundInSession() bool {
	return i.isRoundInSession
}

// SetIsRoundInSession sets whether a round is currently in session
func (i *Info) SetIsRoundInSession(isRoundInSession bool) {
	i.isRoundInSession = isRoundInSession
}

// GetCurrentPlayerIndex returns the index of the current player
func (i *Info) GetCurrentPlayerIndex() int {
	return i.currentPlayerIndex
}

// SetCurrentPlayerIndex sets the index of the current player
func (i *Info) SetCurrentPlayerIndex(index int) {
	i.currentPlayerIndex = index
}

// GetTrumpRank returns the trump rank for the current round
func (i *Info) GetTrumpRank() Rank {
	return i.trumpRank
}

// SetTrumpRank sets the trump rank for the current round
func (i *Info) SetTrumpRank(rank Rank) {
	i.trumpRank = rank
}

// GetGrpScores returns the scores of both groups
func (i *Info) GetGrpScores() [2]int {
	return i.grpScores
}

// SetGrpScores sets the scores of both groups
func (i *Info) SetGrpScores(scores [2]int) {
	i.grpScores = scores
}

// GetFirstFinishedIndex returns the index of the first player to finish a round
func (i *Info) GetFirstFinishedIndex() int {
	return i.firstFinishedIndex
}

// SetFirstFinishedIndex sets the index of the first player to finish a round
func (i *Info) SetFirstFinishedIndex(index int) {
	i.firstFinishedIndex = index
}

// GetSecondFinishedIndex returns the index of the second player to finish a round
func (i *Info) GetSecondFinishedIndex() int {
	return i.secondFinishedIndex
}

// SetSecondFinishedIndex sets the index of the second player to finish a round
func (i *Info) SetSecondFinishedIndex(index int) {
	i.secondFinishedIndex = index
}

// GetLastFinishedIndex returns the index of the last player to finish a round
func (i *Info) GetLastFinishedIndex() int {
	return i.lastFinishedIndex
}

// SetLastFinishedIndex sets the index of the last player to finish a round
func (i *Info) SetLastFinishedIndex(index int) {
	i.lastFinishedIndex = index
}

// GetSecondToLastFinishedIndex returns the index of the second to last player to finish a round
func (i *Info) GetSecondToLastFinishedIndex() int {
	return i.secondToLastFinishedIndex
}

// SetSecondToLastFinishedIndex sets the index of the second to last player to finish a round
func (i *Info) SetSecondToLastFinishedIndex(index int) {
	i.secondToLastFinishedIndex = index
}

// GetLastPlayedCards returns the last played cards
// Note: Returns the actual slice, not a copy
func (i *Info) GetLastPlayedCards() []Card {
	return i.lastPlayedCards
}

// SetLastPlayedCards sets the last played cards
// Note: Stores the provided slice directly, no copying is done
func (i *Info) SetLastPlayedCards(cards []Card) {
	i.lastPlayedCards = cards
}
