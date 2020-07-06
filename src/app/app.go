package app

import (
	"fmt"
	"landau/agile-results/src/checklist"
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

// TODO: Refactor the creators and fetchers to be a trello client wrapper in which
// I wrap the functionality and provide a mockable interface as well. It is
// quiet burdensome to use the trello api and contously provide interdependent
// interfaces for all the actions I want to perform.

// Config -
type Config struct {
	CardCreator      CardCreator
	ChecklistCreator checklist.Creator
	HasChecklist     bool
	LabelFetcher     LabelFetcher
	ListID           string
	Logrus           *logrus.Logger
	Prompter         prompt.Prompter
}

// RunApp -
func RunApp(config *Config) (*trello.Card, error) {
	cardCreator := config.CardCreator
	logrus := config.Logrus
	prompter := config.Prompter

	logrus.Debugf("Creating card  on list %s\n", config.ListID)
	logrus.Trace("Awaiting user input")

	cardName, err := prompter.Prompt("Card Name: ")

	if err != nil {
		return nil, err
	}

	labels, err := config.LabelFetcher.Fetch(trello.Defaults())

	if err != nil {
		return nil, err
	}

	// I think it may be cleaner, def faster, to store things locally.
	// Through a dot file with json inside? Works nicely as a backup space as well.
	// Not sure what it would be for yet. I reduce API usage and have data more
	// readily available. Basically, have everything on hand except for card
	// creation. In fact, I could use that mechanism as way to build my own
	// data model. Perhaps a review of the Trello API to determine some capablities
	// that I am unaware of.
	//
	// The domains would then be relevant locally stored Trello data,
	// the AR domain, which is creating a card in a certain column, making it
	// easier to add labels, and create linked subtasks from a checklist.
	// By removing the Trello API, aside from CardCreator

	// Refactor: This method should not be concerned with `prompter`
	selectedLabels, err := SelectLabelIDs(labels, prompter)
	logrus.Debugf("Selected %d labels: %v", len(selectedLabels), selectedLabels)

	// FIXME: This should put the card at the end of  the list.
	// err = card.MoveToBottomOfList() // How to test this?
	card := &trello.Card{
		IDList:   config.ListID,
		Name:     cardName,
		IDLabels: selectedLabels,
	}

	err = cardCreator.CreateCard(card, trello.Defaults())

	if err != nil {
		return nil, err
	}

	if config.HasChecklist {
		items, err := prompter.PromptList("Checklist Items: ")

		if err != nil {
			return nil, err
		}

		_, err = config.ChecklistCreator.Create(card, items)

		if err != nil {
			// If it fails here, do we delete the card? Then we need a CardDeleter :(
			// See note above about a TrelloWrapper which limits access to the API
			return nil, err
		}
	}

	// TODO: Refactor this to be an expected response like { Card, Curl: {} }?
	logrus.Infof(
		"Card \"%s\" created successfully!\nID: %s\nURL: %s\nDelete: %s\n",
		card.Name, card.ID, card.ShortURL, createCurlDelCmd(card.ID),
	)

	return card, nil
}
