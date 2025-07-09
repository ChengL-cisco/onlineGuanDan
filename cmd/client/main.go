package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/ChengL-cisco/onlineGuanDan/pkg/models"
	"github.com/gorilla/websocket"
)

var (
	serverAddr        = flag.String("server", "localhost:8080", "WebSocket server address")
	name              = flag.String("name", "Player", "Player name")
	reader            = bufio.NewReader(os.Stdin)
	index             = 0
	playerDeck        *models.Deck
	playAttempt       []models.Card
	equivalentAttempt []models.Card
	trumpRank         models.Rank
	finishedIndexes   []int
)

func organizeCards(conn *websocket.Conn) {
	for {
		fmt.Println("Type 'y' to indicate you are ready to start or select the card index or index range:")
		var input string
		fmt.Scan(&input)
		if input == "y" {
			if err := conn.WriteMessage(websocket.TextMessage, models.BuildClientMessage(index, "start", *name)); err != nil {
				log.Printf("Error sending start message: %v", err)
				return
			}
			break
		}

		fmt.Println("Select the destination index:")
		var destStr string
		fmt.Scan(&destStr)

		// Parse source (can be a single index or a range)
		var start, end int
		var srcIndexes []int
		if strings.Contains(input, ",") {
			// Handle multiple index input (e.g., "2,4,6")
			parts := strings.Split(input, ",")
			srcIndexes = make([]int, len(parts))
			for i, part := range parts {
				var err error
				srcIndexes[i], err = strconv.Atoi(part)
				if err != nil || srcIndexes[i] < 0 {
					fmt.Println("Invalid index. Please use positive numbers")
					return
				}
			}
		} else if strings.Contains(input, "-") {
			// Handle range input (e.g., "2-4")
			parts := strings.Split(input, "-")
			if len(parts) != 2 {
				fmt.Println("Invalid range format. Use 'start-end' (e.g., '2-4')")
				return
			}
			var err1, err2 error
			start, err1 = strconv.Atoi(parts[0])
			end, err2 = strconv.Atoi(parts[1])
			if err1 != nil || err2 != nil || start < 0 || end < start {
				fmt.Println("Invalid range. Please use positive numbers in format 'start-end'")
				return
			}
		} else {
			// Handle single index input
			var err error
			start, err = strconv.Atoi(input)
			if err != nil || start < 0 {
				fmt.Println("Invalid index. Please use a positive number")
				return
			}
			end = start
		}

		// Parse destination index
		dest, err := strconv.Atoi(destStr)
		if err != nil || dest < 0 {
			fmt.Println("Invalid destination index. Please use a positive number")
			return
		}

		fmt.Printf("Moving cards from index %d to %d to position %d\n", start, end, dest)
		if strings.Contains(input, ",") {
			playerDeck.MoveNDCards(srcIndexes, dest)
		} else {
			playerDeck.MoveNCards(start, end, dest)
		}
		fmt.Println("First player's sorted cards:")
		fmt.Println(playerDeck.String())
	}
}

func disconnect(conn *websocket.Conn) error {
	// Send a close message to the server
	err := conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "User requested disconnect"),
		time.Now().Add(time.Minute))
	if err != nil {
		log.Printf("Warning: error sending close message: %v", err)
		return fmt.Errorf("error sending close message: %w", err)
	}

	// Wait for the server to acknowledge the close message
	conn.SetReadDeadline(time.Now().Add(5 * time.Second)) // Set a read deadline for the acknowledgement
	for {
		_, _, err = conn.ReadMessage()
		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			break // Server acknowledged the closure
		}
		if err != nil {
			break // An error occurred
		}
	}

	// Close the connection
	return conn.Close()
}

// handleAllJoined handles the allJoined message from the server
func handleAllJoined(conn *websocket.Conn) error {
	log.Printf("Everybody joined, ready to start")
	fmt.Print("Type Y to indicate ready to start or N to quit the game ")
	readyToStart, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}
	readyToStart = strings.ToUpper(strings.TrimSpace(readyToStart))

	if readyToStart == "N" {
		log.Println("User chose to quit the game")
		if err := disconnect(conn); err != nil {
			return fmt.Errorf("error disconnecting: %w", err)
		}
		return fmt.Errorf("user chose to quit the game")
	}

	if err := conn.WriteMessage(websocket.TextMessage, models.BuildClientMessage(index, "ready", *name)); err != nil {
		return fmt.Errorf("error sending ready message: %w", err)
	}

	return nil
}

func getCardsFromIndexes() []models.Card {
	for {
		fmt.Println(playerDeck.String())
		fmt.Println("It's your turn to play!")
		fmt.Println("pick the card indexes to play or type 'p' to pass:")
		var input string
		fmt.Scan(&input)
		fmt.Println(input)
		fmt.Println(input == "p")
		if input == "p" {
			return nil
		}
		var start, end int
		var sourceIndexes []int
		if strings.Contains(input, ",") {
			// Handle multiple index input (e.g., "2,4,6")
			parts := strings.Split(input, ",")
			sourceIndexes = make([]int, len(parts))
			for i, part := range parts {
				var err error
				sourceIndexes[i], err = strconv.Atoi(part)
				if err != nil || sourceIndexes[i] < 0 {
					fmt.Println("Invalid index. Please use positive numbers")
					continue
				}
			}
		} else if strings.Contains(input, "-") {
			// Handle range input (e.g., "2-4")
			parts := strings.Split(input, "-")
			if len(parts) != 2 {
				fmt.Println("Invalid range format. Use 'start-end' (e.g., '2-4')")
				continue
			}
			var err1, err2 error
			start, err1 = strconv.Atoi(parts[0])
			end, err2 = strconv.Atoi(parts[1])
			if err1 != nil || err2 != nil || start < 0 || end < start {
				fmt.Println("Invalid range. Please use positive numbers in format 'start-end'")
				continue
			}
			sourceIndexes = make([]int, end-start+1)
			for i := 0; i <= end-start; i++ {
				sourceIndexes[i] = start + i
			}
		} else {
			// Handle single index input
			var err error
			start, err = strconv.Atoi(input)
			if err != nil || start < 0 {
				fmt.Println("Invalid index. Please use a positive number")
				continue
			}
			sourceIndexes = []int{start}
		}

		cards := make([]models.Card, 0, len(sourceIndexes))
		for _, idx := range sourceIndexes {
			if idx < len(playerDeck.GetCards()) {
				cards = append(cards, playerDeck.GetCards()[idx])
			}
		}
		return cards
	}
}

// selectAndJoinSlot handles the slot selection and join process
func selectAndJoinSlot(conn *websocket.Conn, slotsData string) error {
	slots := strings.Fields(slotsData)
	if len(slots) == 0 {
		return fmt.Errorf("no available slots")
	}

	fmt.Printf("Available slots: %s, pick one to join\n", slotsData)
	fmt.Print("Enter slot number to join: ")

	// Read input with proper error handling
	selectedSlot, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}
	selectedSlot = strings.TrimSpace(selectedSlot)

	// Validate input is not empty
	if selectedSlot == "" {
		return fmt.Errorf("no slot selected")
	}

	log.Printf("Selected slot: %s", selectedSlot)

	index, err = strconv.Atoi(selectedSlot)
	if err != nil {
		return fmt.Errorf("invalid slot number: %w", err)
	}

	log.Printf("Selected index: %d", index)

	// Send the selected slot back to the server
	if err := conn.WriteMessage(websocket.TextMessage, models.BuildClientMessage(index, "join", *name)); err != nil {
		return fmt.Errorf("error sending join message: %w", err)
	}

	return nil
}

// handleWildCards checks for wild cards in the playAttempt and prompts the user for their equivalent
// Returns the equivalent attempt with wild cards replaced, or an error if input is invalid
func handleWildCards(playAttempt []models.Card, trumpRank models.Rank) []models.Card {
	equivalent := make([]models.Card, len(playAttempt))
	copy(equivalent, playAttempt)
	replaced := false

	for i, card := range playAttempt {
		// Check if card is a heart or the current trump rank
		if card.Suit == models.Heart && card.Rank == trumpRank {
			replaced = false
			for !replaced {
				fmt.Printf("Wild card detected: %s\n", card.String())
				fmt.Println("Enter the card this should be equivalent to (e.g., 'A-S' for Ace of Spades):")

				reader := bufio.NewReader(os.Stdin)
				input, err := reader.ReadString('\n')
				if err != nil {
					log.Printf("error reading input: %v", err)
					continue
				}
				input = strings.TrimSpace(input)

				// Parse the input card
				if input == "Jr" || input == "BJr" {
					log.Printf("cannot use jokers as equivalent cards")
					continue
				}

				parts := strings.Split(input, "-")
				if len(parts) != 2 {
					log.Printf("invalid card format, expected format like 'A-S' or '10-H'")
					continue
				}

				rankStr, suitStr := parts[0], parts[1]
				rank, err := models.StringToRank(rankStr)
				if err != nil {
					log.Printf("invalid rank: %v", err)
					continue
				}

				var suit models.Suit
				switch strings.ToUpper(suitStr) {
				case "S":
					suit = models.Spade
				case "H":
					suit = models.Heart
				case "D":
					suit = models.Diamond
				case "C":
					suit = models.Club
				default:
					log.Printf("invalid suit: %s, must be S, H, D, or C", suitStr)
					continue
				}

				if suit == models.Heart && rank == trumpRank {
					log.Printf("wild card not allowed: %s", card.String())
					continue
				}

				equivalent[i] = models.NewCard(suit, rank)
				replaced = true
			}
		}
	}

	return equivalent
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *serverAddr, Path: "/ws"}
	log.Printf("Connecting to %s", u.String())

	// Connect to WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	// Start reading messages from server
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
			}
			log.Printf("Received text message: %s", string(message))
			msg, err := models.ParseServerMessage(message)
			if err != nil {
				log.Printf("Failed to parse message: %v", err)
				continue
			}
			log.Printf("Received message: %v", msg)
			switch msg.Action {
			case "availableSlots":
				if err := selectAndJoinSlot(conn, msg.Data); err != nil {
					log.Printf("Error joining slot: %v", err)
					return
				}
			case "joinConfirm":
				log.Printf("Joined successfully")
			case "allJoined":
				if err := handleAllJoined(conn); err != nil {
					log.Printf("Error handling all joined: %v", err)
					// Continue to wait for another input if user didn't confirm
					if err.Error() == "user not ready to start" {
						continue
					}
					return
				}
			case "startRound":
				deck, tRank, fIndexes, err := models.ParseStartRoundServerMessage(msg.Data)
				if err != nil {
					log.Printf("Failed to parse start round message: %v", err)
					return
				}

				// Store the deck for future use
				playerDeck = deck
				trumpRank = tRank
				finishedIndexes = fIndexes
				fmt.Println(deck.String())
				fmt.Printf("Trump rank: %s\n", models.RankToString(trumpRank))
				fmt.Printf("Finished indexes: %v\n", finishedIndexes)
				organizeCards(conn)

			case "play":
				playerIndex, err := strconv.Atoi(msg.Data)
				if err != nil {
					log.Printf("Failed to parse index: %v", err)
					return
				}
				fmt.Printf("Player %d's turn\n", playerIndex)
				if index == playerIndex {
					cards := getCardsFromIndexes()
					if cards == nil {
						conn.WriteMessage(websocket.TextMessage, models.BuildClientMessage(index, "pass", ""))
					} else {
						playAttempt = cards
						equivalentAttempt = handleWildCards(cards, trumpRank)
						conn.WriteMessage(websocket.TextMessage, models.BuildClientMessage(index, "playAttempt", models.ConstructClientPlayMessage(cards, playerDeck.Count()-len(playAttempt), equivalentAttempt)))
					}
				}
			case "invalidPlay":
				fmt.Println("Invalid play, trying again")
				cards := getCardsFromIndexes()
				playAttempt = cards

				equivalentAttempt = handleWildCards(cards, trumpRank)
				conn.WriteMessage(websocket.TextMessage, models.BuildClientMessage(index, "playAttempt", models.ConstructClientPlayMessage(cards, playerDeck.Count()-len(playAttempt), equivalentAttempt)))

			case "validPlay":
				fmt.Println("Valid play")
				playerDeck.PlayN(playAttempt)
				fmt.Println(playerDeck.String())
				msg := models.ConstructClientPlayMessage(playAttempt, playerDeck.Count(), equivalentAttempt)
				conn.WriteMessage(websocket.TextMessage, models.BuildClientMessage(index, "play", msg))
			case "lastPlay":
				fmt.Println("Last play:")
				playerIndex, numCardsLeft, attemptDeck, equivalentDeck, err := models.ParseLastPlayServerMessage(msg.Data)
				if err != nil {
					log.Printf("Failed to parse last play message: %v", err)
					return
				}
				fmt.Printf("Player %d's last play:\n", playerIndex)
				fmt.Println(attemptDeck.String())
				fmt.Printf("Number of cards left: %d\n", numCardsLeft)
				fmt.Printf("Equivalent play:\n")
				fmt.Println(equivalentDeck.String())
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("Interrupt received, closing connection...")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}

	}
}
