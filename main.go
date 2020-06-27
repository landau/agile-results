package main

import (
	"fmt"
	"log"
	"os"

	"github.com/adlio/trello"
)

func main() {
	// TODO: assert that these values are set
	apiKey := os.Getenv("TRELLO_API_KEY")
	token := os.Getenv("TRELLO_TOKEN")
	listID := os.Getenv("TRELLO_LIST_ID")

	client := trello.NewClient(apiKey, token)

	// TODO: This should print in a verbose mode
	fmt.Printf("Creating card  on list %s\n", listID)

	// FIXME: This should go to end of list.
	card := &trello.Card{Name: "Test", IDList: listID}
	err := client.CreateCard(card, trello.Defaults())

	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	deleteCmd := fmt.Sprintf(
		"curl -s -XDELETE \"https://api.trello.com/1/cards/%s?key=$TRELLO_API_KEY&token=$TRELLO_API_TOKEN\"",
		card.ID,
	)
	fmt.Printf(
		"\nCard %s created successfully!\nID: %s\nURL: %s\nDelete: %s\n",
		card.Name, card.ID, card.ShortURL, deleteCmd,
	)
}
