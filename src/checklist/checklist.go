package checklist

import (
	"github.com/adlio/trello"
)

type trelloChecklistHandler interface {
	CreateChecklist(card *trello.Card, name string, extraArgs trello.Arguments) (*trello.Checklist, error)
	CreateCheckItem(checklist *trello.Checklist, name string, extraArgs trello.Arguments) (*trello.CheckItem, error)
}

// Creator -
type Creator interface {
	Create(card *trello.Card, items []string) (*trello.Checklist, error)
}

// TrelloCreator -
type TrelloCreator struct {
	handler trelloChecklistHandler
}

// Create -
func (c *TrelloCreator) Create(card *trello.Card, items []string) (*trello.Checklist, error) {
	checklist, err := c.handler.CreateChecklist(card, "Checklist", trello.Defaults())

	if err != nil {
		return nil, err
	}

	for _, name := range items {
		_, err = c.handler.CreateCheckItem(checklist, name, trello.Defaults())

		if err != nil {
			return nil, err
		}
	}

	return checklist, nil
}

type (
	// MockCreator -
	MockCreator struct {
		Calls             []*MockCreateCall
		CreateReturnValue MockCreateReturnValue
	}

	// MockCreateCall -
	MockCreateCall struct {
		Card  *trello.Card
		Items []string
	}

	// MockCreateReturnValue -
	MockCreateReturnValue struct {
		checklist *trello.Checklist
		err       error
	}
)

// Create -
func (c *MockCreator) Create(card *trello.Card, items []string) (*trello.Checklist, error) {
	c.Calls = append(c.Calls, &MockCreateCall{Card: card, Items: items})
	return c.CreateReturnValue.checklist, c.CreateReturnValue.err
}

// New -
func New(handler trelloChecklistHandler) *TrelloCreator {
	return &TrelloCreator{handler: handler}
}
