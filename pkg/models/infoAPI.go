package models

// InfoAPI defines the interface for accessing and modifying game information
type InfoAPI interface {
	// Player counts and indexes
	GetNumPlayers() int
	SetNumPlayers(numPlayers int)

	// Group names
	GetGrp1Name() string
	SetGrp1Name(name string)
	GetGrp2Name() string
	SetGrp2Name(name string)

	// Ready states
	GetReadyToStartIndexes() []int
	SetReadyToStartIndexes(indexes []int)
	ResetReadyToStartIndexes()
	GetReadyToPlayIndexes() []int
	SetReadyToPlayIndexes(indexes []int)
	ResetReadyToPlayIndexes()

	// Round information
	GetIsFirstRound() bool
	SetIsFirstRound(isFirstRound bool)
	GetIsRoundInSession() bool
	SetIsRoundInSession(isRoundInSession bool)

	// Current player
	GetCurrentPlayerIndex() int
	SetCurrentPlayerIndex(index int)

	// Trump rank
	GetTrumpRank() Rank
	SetTrumpRank(rank Rank)

	// Group scores
	GetGrpScores() [2]int
	SetGrpScores(scores [2]int)

	// Finished player indexes
	GetFirstFinishedIndex() int
	SetFirstFinishedIndex(index int)
	GetSecondFinishedIndex() int
	SetSecondFinishedIndex(index int)
	GetLastFinishedIndex() int
	SetLastFinishedIndex(index int)
	GetSecondToLastFinishedIndex() int
	SetSecondToLastFinishedIndex(index int)

	// Last played cards
	GetLastPlayedCards() []Card
	SetLastPlayedCards(cards []Card)
}

// Verify at compile time that *Info implements InfoAPI
var _ InfoAPI = (*Info)(nil)
