package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/adlio/trello"
	"github.com/sirupsen/logrus"
)

// Flags - Supported CLI flags
type Flags struct {
	verbose bool
}

// parseFlags - Parses CLI flags
func parseFlags() Flags {
	verbose := flag.Bool("v", false, "Verbose")
	flag.Parse()
	return Flags{verbose: *verbose}
}

func main() {

	flags := parseFlags()

	if flags.verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// TODO: assert that these values are set
	apiKey := os.Getenv("TRELLO_API_KEY")
	token := os.Getenv("TRELLO_TOKEN")
	listID := os.Getenv("TRELLO_LIST_ID")

	client := trello.NewClient(apiKey, token)

	logrus.Debugf("Creating card  on list %s\n", listID)

	// FIXME: This should go to end of list.
	card := &trello.Card{Name: "Test", IDList: listID}
	err := client.CreateCard(card, trello.Defaults())

	if err != nil {
		logrus.Fatalf("Failed: %v", err)
	}

	deleteCmd := fmt.Sprintf(
		"curl -sXDELETE \"https://api.trello.com/1/cards/%s?key=$TRELLO_API_KEY&token=$TRELLO_API_TOKEN\"",
		card.ID,
	)
	logrus.Infof(
		"Card %s created successfully!\nID: %s\nURL: %s\nDelete: %s\n",
		card.Name, card.ID, card.ShortURL, deleteCmd,
	)
}
