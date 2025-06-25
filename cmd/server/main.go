package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ChengL-cisco/onlineGuanDan/pkg/models"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client represents a connected WebSocket client
type Client struct {
	conn  *websocket.Conn
	send  chan []byte
	Index int // Add Index field to track player index
}

// Hub maintains the set of active clients
var deck models.DeckAPI
var (
	info       = &models.Info{}
	clients    = make(map[int]*Client) // Map of player index to Client
	broadcast  = make(chan []byte)     // Broadcast channel
	mutex      = &sync.Mutex{}         // Mutex to protect clients map
	firstRound = true
)

// handleWebSocket handles WebSocket requests from clients
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// Create new client
	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}

	log.Printf("New client connected.")
	// Get and send available slots to the client
	availableSlots := getAvailableSlots()
	client.conn.WriteMessage(websocket.TextMessage, models.BuildServerMessage("availableSlots", availableSlots))

	// Start goroutine for reading messages
	go client.processClientMsg()
}

// getAvailableSlots returns a space-separated string of available slot numbers
func getAvailableSlots() string {
	mutex.Lock()
	defer mutex.Unlock()

	availableSlots := info.GetAvailableSlots()
	keys := make([]int, 0, len(availableSlots))
	for k := range availableSlots {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// Join the sorted keys into a space-separated string
	msg := ""
	for _, k := range keys {
		msg += fmt.Sprintf("%d ", k)
	}
	return strings.TrimSpace(msg) // Remove trailing space
}

// processClientMsg processes messages from the WebSocket connection
func (c *Client) processClientMsg() {
	defer func() {
		mutex.Lock()
		info.GetAvailableSlots()[c.Index] = true
		delete(info.GetNames(), c.Index)
		delete(clients, c.Index)
		mutex.Unlock()

		c.conn.Close()
		log.Printf("Client disconnected.")
	}()

	for {
		log.Printf("Waiting for message...")
		messageType, message, err := c.conn.ReadMessage()
		log.Printf("Received message: %v", string(message))
		if err != nil {
			log.Printf("WebSocket error: %v", err)
		}

		switch messageType {
		case websocket.TextMessage:
			// Handle text message
			msg, err := models.ParseClientMessage(message)
			if err != nil {
				log.Printf("Failed to parse message: %v", err)
				c.conn.WriteMessage(websocket.TextMessage, models.BuildServerMessage("error", fmt.Sprintf("Failed to parse message: %v", err)))
				continue
			}

			switch msg.Action {
			case "join":
				log.Printf("Client %s wants to join", msg.Data)
				mutex.Lock()
				availableSlots := info.GetAvailableSlots()
				log.Printf("availableSlots: %v", availableSlots)

				if _, exists := availableSlots[msg.Index]; !exists {
					// slot no longer available
					mutex.Unlock()
					c.conn.WriteMessage(websocket.TextMessage, models.BuildServerMessage("availableSlots", getAvailableSlots()))
				} else {
					c.Index = msg.Index
					clients[msg.Index] = c
					delete(availableSlots, msg.Index)
					names := info.GetNames()
					names[msg.Index] = msg.Data
					info.SetNames(names)

					c.conn.WriteMessage(websocket.TextMessage, models.BuildServerMessage("joinConfirm", ""))

					// to do: if everybody joined, broadcast to ready to start
					if len(clients) == info.GetNumPlayers() {
						log.Printf("Everybody joined, getting ready...")
						broadcastMessage(models.BuildServerMessage("allJoined", ""))
					}
					mutex.Unlock()
				}

			case "ready":
				log.Printf("Client %s is ready", msg.Data)
				mutex.Lock()
				info.GetReadyToStartMap()[msg.Index] = true
				// if everybody is ready, send out the cards
				if len(info.GetReadyToStartMap()) == info.GetNumPlayers() {
					log.Printf("Everybody is ready, starting the game...")
					info.SetIsRoundInSession(true)
					if firstRound {
						rand.Seed(time.Now().UnixNano())
						info.SetCurrentPlayerIndex(rand.Intn(info.GetNumPlayers()))
						info.SetTrumpRank(models.Two)
						firstRound = false
					}
					// reset ready to start map
					info.SetReadyToStartMap(make(map[int]bool))

					deck = models.NewDeck(models.NumOfDecks(info.GetNumPlayers()))
					decks := deck.Split(info.GetNumPlayers())
					for index, deck := range decks {
						deck.Sort(info.GetTrumpRank())
						clients[index].conn.WriteMessage(websocket.TextMessage, models.BuildServerMessage("cards", models.CardsString(deck.GetCards())))
					}

				}
				log.Printf("prepare to unlock")
				mutex.Unlock()
				log.Printf("unlocked")
			case "start":
				log.Printf("Client %s started", msg.Data)
				mutex.Lock()
				info.GetReadyToPlay()[msg.Index] = true
				// if everybody is ready, start the round
				if len(info.GetReadyToPlay()) == info.GetNumPlayers() {
					log.Printf("Everybody is ready, starting the round...")
					// reset ready to play map
					info.SetReadyToPlay(make(map[int]bool))
					broadcastMessage(models.BuildServerMessage("play", fmt.Sprintf("%d", info.GetCurrentPlayerIndex())))
				}
				mutex.Unlock()
			case "tribute":
				log.Printf("Client %s tributed", msg.Data)
			case "return":
				log.Printf("Client %d returned", msg.Index)
			case "play":
				log.Printf("Client %d played", msg.Index)
			case "pass":
				log.Printf("Client %d passed", msg.Index)
			case "leave":
				log.Printf("Client %d left", msg.Index)
			default:
				log.Printf("Unknown action: %s", msg.Action)
			}
		case websocket.CloseMessage:
			// to do: release available slots
			return
		default:
			// Handle other message types
			log.Printf("Received message of type %d", messageType)
			return
		}
	}
}

// broadcastMessage sends a message to all connected clients
func broadcastMessage(message []byte) {
	for _, client := range clients {
		if err := client.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Error broadcasting to client: %v", err)
			continue
		}
	}
}

func main() {
	// Define command-line flags
	numPlayers := flag.Int("players", 2, "Number of players in the game")
	port := flag.Int("port", 8080, "Port to run the server on")
	flag.Parse()

	// Initialize game info
	info.SetNumPlayers(*numPlayers)
	clients = make(map[int]*Client)

	// Initialize available slots
	availableSlots := make(map[int]bool)
	for i := 0; i < *numPlayers; i++ {
		availableSlots[i] = true
	}
	info.SetAvailableSlots(availableSlots)

	// Configure WebSocket route
	http.HandleFunc("/ws", handleWebSocket)

	// Start the server
	serverAddr := fmt.Sprintf(":%d", *port)
	log.Printf("Server starting on port %d...\n", *port)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
