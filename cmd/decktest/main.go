package main

import (
	"fmt"

	"github.com/ChengL-cisco/onlineGuanDan/pkg/models"
)

func main() {
	attemptDeck, numCardsLeft, equivalentDeck, err := models.ParseClientPlayMessage("K-H K-C K-D A-H A-C A-D;21;")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(attemptDeck.String())
	fmt.Println(equivalentDeck.String())
	fmt.Println(numCardsLeft)
}
