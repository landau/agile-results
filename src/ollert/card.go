package ollert

import (
	"github.com/adlio/trello"
)

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

type (
	// MockCardConfig -
	MockCardConfig struct {
		TrelloCardConfig
		ID       string
		Name     string
		Card     *trello.Card
		ShortURL string
		IDLabels []string
		IDList   string

		// Methods
		MoveToBottomOfListReturnValue error
	}

	// MockCard -
	MockCard struct {
		id       string
		name     string
		card     *trello.Card
		shortURL string
		idLabels []string

		// Mocked methods
		MoveToBottomOfListCallCount   int
		moveToBottomOfListReturnValue error
	}
)

// MoveToBottomOfList -
func (c *MockCard) MoveToBottomOfList() error {
	c.MoveToBottomOfListCallCount++
	return c.moveToBottomOfListReturnValue
}

// Card -
func (c *MockCard) Card() *trello.Card {
	return c.card
}

// ID -
func (c *MockCard) ID() string {
	return c.id
}

// Name -
func (c *MockCard) Name() string {
	return c.name
}

// ShortURL -
func (c *MockCard) ShortURL() string {
	return c.shortURL
}

// IDLabels -
func (c *MockCard) IDLabels() []string {
	return c.idLabels
}

// NewMockCard - Use this to create a new Trello MockCard
func NewMockCard(c *MockCardConfig) *MockCard {
	return &MockCard{
		id:       c.ID,
		name:     c.Name,
		card:     c.Card,
		shortURL: c.ShortURL,
		idLabels: c.IDLabels,

		moveToBottomOfListReturnValue: c.MoveToBottomOfListReturnValue,
	}
}
