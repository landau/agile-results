package app

import (
	"fmt"
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

// CardCreator -
type CardCreator interface {
	CreateCard(card *trello.Card, extraArgs trello.Arguments) error
}

// Config -
type Config struct {
	Logrus      *logrus.Logger
	CardCreator CardCreator
	Prompter    prompt.Prompter
	ListID      string
}

// RunApp -
func RunApp(config *Config) error {
	cardCreator := config.CardCreator
	logrus := config.Logrus
	prompter := config.Prompter

	logrus.Debugf("Creating card  on list %s\n", config.ListID)
	logrus.Trace("Awaiting user input")

	cardName, err := prompter.Prompt("Card Name: ")

	if err != nil {
		return err
	}

	// FIXME: This should put the card at the end of  the list.
	card := &trello.Card{Name: cardName, IDList: config.ListID}
	err = cardCreator.CreateCard(card, trello.Defaults())

	if err != nil {
		return err
	}

	// TODO: Refactor this to be an expected response like { Card, Curl: {} }?
	logrus.Infof(
		"Card \"%s\" created successfully!\nID: %s\nURL: %s\nDelete: %s\n",
		card.Name, card.ID, card.ShortURL, createCurlDelCmd(card.ID),
	)

	return nil
}
