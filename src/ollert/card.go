package ollert

import "github.com/adlio/trello"

// ICard -
type ICard interface {
	Card() *trello.Card
	ID() string
	Name() string
	ShortURL() string
	MoveToBottomOfList() error
	IDLabels() []string
}

// Card -
type Card struct {
	card *trello.Card
}

// MoveToBottomOfList -
func (c *Card) MoveToBottomOfList() error {
	return c.card.MoveToBottomOfList()
}

// Card -
func (c *Card) Card() *trello.Card {
	return c.card
}

// ID -
func (c *Card) ID() string {
	return c.card.ID
}

// Name -
func (c *Card) Name() string {
	return c.card.Name
}

// ShortURL -
func (c *Card) ShortURL() string {
	return c.card.ShortURL
}

// IDLabels -
func (c *Card) IDLabels() []string {
	return c.card.IDLabels
}

// NewCard - Use this to create a new Trello Card
func NewCard(c *trello.Card) ICard {
	return &Card{card: c}
}
