package models

import (
	"encoding/json"
	"fmt"
	"log"
)

// action can be ["join", "ready", "start", "tribute", "return", "play", "pass", "leave"]
type ClientMessage struct {
	Index  int    `json:"index"`
	Action string `json:"action"`
	Data   string `json:"data"`
}

// ServerMessage represents a message sent from server to client
// action can be ["availableSlots", "joinConfirm", "allJoined", "cards", "play"]
type ServerMessage struct {
	Action string `json:"action"`
	Data   string `json:"data"`
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
