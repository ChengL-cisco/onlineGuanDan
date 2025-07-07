package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// action can be ["join", "ready", "start", "tribute", "return", "playAttempt", "play", "pass", "leave"]
type ClientMessage struct {
	Index  int    `json:"index"`
	Action string `json:"action"`
	Data   string `json:"data"`
}

// ServerMessage represents a message sent from server to client
// action can be ["availableSlots", "joinConfirm", "allJoined", "startRound", "play", "validPlay", "invalidPlay", "lastPlay"]
type ServerMessage struct {
	Action string `json:"action"`
	Data   string `json:"data"`
}

// ConstructStartRoundServerMessage constructs a start round message string from the given deck and info
// deck is the deck of cards
// info is the game info
func ConstructStartRoundServerMessage(deck DeckAPI, info InfoAPI) string {
	return CardsString(deck.GetCards()) + ";" + RankToString(info.GetTrumpRank()) + ";" + strings.Trim(fmt.Sprint(info.GetFinishedIndexes()), "[]")
}

func ParseStartRoundServerMessage(msg string) (DeckAPI, Rank, []int, error) {
	// Split the message by semicolon
	parts := strings.SplitN(msg, ";", 3)
	if len(parts) != 3 {
		return nil, 0, nil, fmt.Errorf("invalid message format: expected 3 parts separated by ';'")
	}

	deckStr := parts[0]
	trumpRankStr := parts[1]
	finishedIndexesStr := parts[2]

	// Parse deck
	deck, err := NewDeckFromString(deckStr)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("failed to parse deck: %v", err)
	}

	// Parse trump rank
	trumpRank, err := StringToRank(trumpRankStr)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("failed to parse trump rank: %v", err)
	}

	// Parse finished indexes (comma-separated integers)
	var finishedIndexes []int
	if finishedIndexesStr != "" {
		indexStrs := strings.Split(finishedIndexesStr, ",")
		for _, idxStr := range indexStrs {
			idx, err := strconv.Atoi(strings.TrimSpace(idxStr))
			if err != nil {
				return nil, 0, nil, fmt.Errorf("failed to parse finished index '%s': %v", idxStr, err)
			}
			finishedIndexes = append(finishedIndexes, idx)
		}
	}

	return deck, trumpRank, finishedIndexes, nil
}

// ConstructClientPlayMessage constructs a play message string from the given attempt and equivalent cards
// equivalent can be nil
func ConstructClientPlayMessage(attempt []Card, numCardsLeft int, equivalent []Card) string {
	if equivalent == nil {
		equivalent = []Card{}
	}
	return fmt.Sprintf("%s;%d;%s", CardsString(attempt), numCardsLeft, CardsString(equivalent))
}

// ParseClientPlayMessage parses a play message string into its components
func ParseClientPlayMessage(msg string) (DeckAPI, int, DeckAPI, error) {
	// Split the message by semicolon
	parts := strings.SplitN(msg, ";", 3)
	if len(parts) != 3 {
		return nil, 0, nil, fmt.Errorf("invalid message format: expected 3 parts separated by ';'")
	}

	attemptStr := parts[0]
	numCardsLeft, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, 0, nil, fmt.Errorf("invalid number of cards left: %v", err)
	}
	equivalentStr := parts[2]

	// Parse attempt cards
	attemptDeck, err := NewDeckFromString(attemptStr)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("failed to parse attempt cards: %v", err)
	}

	// Parse equivalent cards
	equivalentDeck, err := NewDeckFromString(equivalentStr)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("failed to parse equivalent cards: %v", err)
	}

	return attemptDeck, numCardsLeft, equivalentDeck, nil
}

// BuildClientMessage is a helper function to build a structured client message
func BuildClientMessage(index int, action string, data interface{}) []byte {
	msg := ClientMessage{
		Index:  index,
		Action: action,
	}

	// Convert data to string if it's not already
	if str, ok := data.(string); ok {
		msg.Data = str
	} else if data != nil {
		msg.Data = fmt.Sprintf("%v", data)
	}

	message, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling client message: %v", err)
		return nil
	}
	return message
}

// ParseClientMessage parses a JSON-encoded client message into a ClientMessage struct.
// It returns the parsed message and any error encountered.
func ParseClientMessage(data []byte) (*ClientMessage, error) {
	var msg ClientMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse client message: %w", err)
	}

	// Validate required fields
	if msg.Action == "" {
		return nil, fmt.Errorf("missing required field: action")
	}

	return &msg, nil
}

// ParseServerMessage parses a JSON-encoded server message into a ServerMessage struct.
// It returns the parsed message and any error encountered.
func ParseServerMessage(data []byte) (*ServerMessage, error) {
	var msg ServerMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse server message: %w", err)
	}

	// Validate required fields
	if msg.Action == "" {
		return nil, fmt.Errorf("missing required field: action")
	}

	return &msg, nil
}

// BuildServerMessage is a helper function to build a structured server message
func BuildServerMessage(action string, data interface{}) []byte {
	msg := ServerMessage{
		Action: action,
	}

	// Convert data to string if it's not already
	if str, ok := data.(string); ok {
		msg.Data = str
	} else if data != nil {
		msg.Data = fmt.Sprintf("%v", data)
	}

	message, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling server message: %v", err)
		return nil
	}
	return message
}
