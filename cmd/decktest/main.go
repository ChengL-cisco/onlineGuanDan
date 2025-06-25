package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/ChengL-cisco/onlineGuanDan/pkg/models"
)

func main() {
	// Define command-line flags
	numDecks := flag.Int("decks", 1, "Number of decks to create")
	flag.Parse()

	// Create a new deck with the specified number of decks
	deck := models.NewDeck(*numDecks)

	// Print deck information
	fmt.Printf("Created deck with %d decks (%d cards total)\n", *numDecks, deck.Count())

	cards := deck.Split(2)
	// first player's card
	fmt.Println("First player's cards:")
	fmt.Println(cards[0].String())
	cards[0].Sort(models.Ten)
	fmt.Println("First player's sorted cards:")
	//fmt.Println(models.CardsString(cards[0].GetCards()))
	fmt.Println(cards[0].String())
	for {
		fmt.Println("Type 'y' to indicate you are ready to start or select the card index or index range:")
		var input string
		fmt.Scan(&input)
		if input == "y" {
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
			cards[0].MoveNDCards(srcIndexes, dest)
		} else {
			cards[0].MoveNCards(start, end, dest)
		}
		fmt.Println("First player's sorted cards:")
		fmt.Println(cards[0].String())
	}
}
