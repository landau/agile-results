package app

import (
	"fmt"
	"landau/agile-results/src/prompt"
	"strings"

	"github.com/adlio/trello"
	"github.com/sirupsen/logrus"
)

func createCurlDelCmd(cardID string) string {
	return fmt.Sprintf(
		"curl -sXDELETE \"https://api.trello.com/1/cards/%s?key=$TRELLO_API_KEY&token=$TRELLO_API_TOKEN\"",
		cardID,
	)
}

func mapLabelsToLabelNames(labels []*trello.Label) []string {
	labelNames := make([]string, len(labels))

	for i, l := range labels {
		labelNames[i] = l.Name
	}

	return labelNames
}

func indexOf(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

func selectLabelsFromLabelNames(labelNames []string, labels []*trello.Label) []string {
	selectedLabels := make([]string, 0)

	for _, labelName := range labelNames {
		i := indexOf(len(labels), func(i int) bool { return labels[i].Name == labelName })
		selectedLabels = append(selectedLabels, labels[i].ID)
	}

	return selectedLabels
}

func selectLabels(labels []*trello.Label, prompter prompt.Prompter) ([]string, error) {
	labelNames := mapLabelsToLabelNames(labels)
	selected, err := prompter.Prompt(
		fmt.Sprintf("Select labels (%v): ", strings.Join(labelNames, ", ")),
	)

	if err != nil {
		return make([]string, 0), err
	}

	return selectLabelsFromLabelNames(strings.Split(selected, ","), labels), nil
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
	Labels      []*trello.Label
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

	selectedLabels, err := selectLabels(config.Labels, prompter)
	logrus.Debugf("Selected %d labels: %v", len(selectedLabels), selectedLabels)

	// FIXME: This should put the card at the end of  the list.
	card := &trello.Card{
		IDList:   config.ListID,
		Name:     cardName,
		IDLabels: selectedLabels,
	}

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
