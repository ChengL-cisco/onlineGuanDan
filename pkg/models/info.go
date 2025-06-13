package models

// Info is a placeholder struct for game information
type Info struct {
	// numPlayers is the number of players in the game
	numPlayers int
	// grp1Name is the name of group 1
	grp1Name string
	// grp2Name is the name of group 2
	grp2Name string
	// readyToStartIndexes is the list of indexes of players who are ready to start
	readyToStartIndexes []int
	// readyToPlayIndexes is the list of indexes of players who are ready to play
	readyToPlayIndexes []int
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

// GetReadyToStartIndexes returns the list of player indexes who are ready to start
func (i *Info) GetReadyToStartIndexes() []int {
	return i.readyToStartIndexes
}

// SetReadyToStartIndexes sets the list of player indexes who are ready to start
func (i *Info) SetReadyToStartIndexes(indexes []int) {
	i.readyToStartIndexes = make([]int, len(indexes))
	copy(i.readyToStartIndexes, indexes)
}

// GetReadyToPlayIndexes returns the list of player indexes who are ready to play
func (i *Info) GetReadyToPlayIndexes() []int {
	return i.readyToPlayIndexes
}

// SetReadyToPlayIndexes sets the list of player indexes who are ready to play
func (i *Info) SetReadyToPlayIndexes(indexes []int) {
	i.readyToPlayIndexes = make([]int, len(indexes))
	copy(i.readyToPlayIndexes, indexes)
}

// ResetReadyToStartIndexes resets all readyToStartIndexes to 0
func (i *Info) ResetReadyToStartIndexes() {
	for j := range i.readyToStartIndexes {
		i.readyToStartIndexes[j] = 0
	}
}

// ResetReadyToPlayIndexes resets all readyToPlayIndexes to 0
func (i *Info) ResetReadyToPlayIndexes() {
	for j := range i.readyToPlayIndexes {
		i.readyToPlayIndexes[j] = 0
	}
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
