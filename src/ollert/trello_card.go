package ollert

import (
	"github.com/adlio/trello"
)

// TrelloCardConfig -
type TrelloCardConfig struct {
	ID       string
	Name     string
	IDLabels []string
	IDList   string
}

// NewTrelloCard -
func NewTrelloCard(c *TrelloCardConfig) *trello.Card {
	return &trello.Card{
		ID:       c.ID,
		Name:     c.Name,
		IDLabels: c.IDLabels,
		IDList:   c.IDList,
	}
}
