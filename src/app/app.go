package app

import (
	"fmt"
	"landau/agile-results/src/ollert"
	"landau/agile-results/src/prompt"

	"github.com/adlio/trello"
	"github.com/sirupsen/logrus"
)

func createCurlDelCmd(cardID string) string {
	return fmt.Sprintf(
		"curl -sXDELETE \"https://api.trello.com/1/cards/%s?key=$TRELLO_API_KEY&token=$TRELLO_API_TOKEN\"",
		cardID,
	)
}

// Config -
type Config struct {
	BoardID      string
	ListID       string
	HasChecklist bool

	Client   ollert.IClient
	Logrus   *logrus.Logger
	Prompter prompt.Prompter
}

// RunApp -
func RunApp(config *Config) (ollert.ICard, error) {
	client := config.Client
	logrus := config.Logrus
	prompter := config.Prompter

	logrus.Debugf("Creating card  on list %s\n", config.ListID)
	logrus.Trace("Awaiting user input")

	cardName, err := prompter.Prompt("Card Name: ")

	if err != nil {
		return nil, err
	}

	board, err := client.GetBoard(config.BoardID, ollert.Defaults())

	if err != nil {
		return nil, err
	}

	labels, err := board.GetLabels(ollert.Defaults())

	if err != nil {
		return nil, err
	}

	// TODO: I think it may be cleaner/faster to store things locally.
	// Through a dot file with json inside?
	// TODO: This method should not be concerned with `prompter`
	selectedLabels, err := SelectLabelIDs(labels, prompter)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("Selected %d labels: %v", len(selectedLabels), selectedLabels)

	// FIXME: This should put the card at the end of  the list.
	// err = card.MoveToBottomOfList() // TODO: add test using ICard

	card, err := client.CreateCard(
		&trello.Card{
			IDList:   config.ListID,
			Name:     cardName,
			IDLabels: selectedLabels,
		},
		ollert.Defaults(),
	)

	if err != nil {
		return nil, err
	}

	if config.HasChecklist {
		items, err := prompter.PromptList("Checklist Items: ")

		if err != nil {
			return nil, err
		}

		_, err = client.CreateChecklist(card, "Checklist", items, ollert.Defaults())

		if err != nil {
			// If it fails here, do we delete the card?
			return nil, err
		}
	}

	// TODO: Refactor this to be an expected response like { Card, Curl: {} }?
	logrus.Infof(
		"Card \"%s\" created successfully!\nID: %s\nURL: %s\nDelete: %s\n",
		card.Name(), card.ID(), card.ShortURL(), createCurlDelCmd(card.ID()),
	)

	return card, nil
}
