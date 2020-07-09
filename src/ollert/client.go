package ollert

import (
	"github.com/adlio/trello"
	"github.com/sirupsen/logrus"
)

// IClient - interface for client
type IClient interface {
	CreateCard(card *trello.Card, args trello.Arguments) (ICard, error)
	GetBoard(boardID string, args trello.Arguments) (IBoard, error)
	CreateChecklist(card ICard, name string, items []string, args trello.Arguments) (IChecklist, error)
}

// --- Client Implementations

// Client - A trello wrapper client, exposing only what we need
type Client struct {
	client ITrelloClient
}

// GetBoard -
func (c *Client) GetBoard(boardID string, args trello.Arguments) (IBoard, error) {
	board, err := c.client.GetBoard(boardID, args)

	if err != nil {
		return nil, err
	}

	return NewBoard(board), nil
}

// CreateCard -
func (c *Client) CreateCard(card *trello.Card, args trello.Arguments) (ICard, error) {
	err := c.client.CreateCard(card, args)

	if err != nil {
		return nil, err
	}

	return NewCard(card), nil
}

// CreateChecklist -
func (c *Client) CreateChecklist(card ICard, name string, items []string, args trello.Arguments) (IChecklist, error) {
	checklist, err := c.client.CreateChecklist(card.Card(), "Checklist", args)

	if err != nil {
		return nil, err
	}

	for _, name := range items {
		_, err = c.client.CreateCheckItem(checklist, name, args)

		if err != nil {
			return nil, err
		}
	}

	return NewCheckList(checklist), nil

}

// NewClient - Use this to create a real trello client
func NewClient(key string, token string, logger *logrus.Logger) IClient {
	trelloClient := trello.NewClient(key, token)
	trelloClient.Logger = logger
	c := &Client{client: trelloClient}
	return c
}

// --- Mock Client Implementations

type (
	// MockClientConfig -
	MockClientConfig struct {
		GetBoardReturn        GetBoardReturn
		CreateCardReturn      CreateCardReturn
		CreateChecklistReturn CreateChecklistReturn
	}

	// MockClient - A trello wrapper client, exposing only what we need
	MockClient struct {
		GetBoardCalls  []GetBoardCall
		getBoardReturn GetBoardReturn

		CreateCardCalls  []CreateCardCall
		createCardReturn CreateCardReturn

		CreateChecklistCalls  []CreateChecklistCall
		createChecklistReturn CreateChecklistReturn
	}

	// GetBoardCall -
	GetBoardCall struct {
		BoardID string
		Args    trello.Arguments
	}

	// GetBoardReturn -
	GetBoardReturn struct {
		Board IBoard
		Err   error
	}

	// CreateCardCall -
	CreateCardCall struct {
		Card *trello.Card
		Args trello.Arguments
	}

	// CreateCardReturn -
	CreateCardReturn struct {
		Card ICard
		Err  error
	}

	// CreateChecklistCall -
	CreateChecklistCall struct {
		Card  ICard
		Name  string
		Items []string
		Args  trello.Arguments
	}

	// CreateChecklistReturn -
	CreateChecklistReturn struct {
		Checklist IChecklist
		Err       error
	}
)

// GetBoard -
func (c *MockClient) GetBoard(boardID string, args trello.Arguments) (IBoard, error) {
	c.GetBoardCalls = append(c.GetBoardCalls, GetBoardCall{BoardID: boardID, Args: args})
	return c.getBoardReturn.Board, c.getBoardReturn.Err
}

// CreateCard -
func (c *MockClient) CreateCard(card *trello.Card, args trello.Arguments) (ICard, error) {
	c.CreateCardCalls = append(c.CreateCardCalls, CreateCardCall{Card: card, Args: args})

	// If a user needs to specify a card returned, then use this pathway
	if c.createCardReturn.Card != nil {
		return c.createCardReturn.Card, c.createCardReturn.Err
	}

	return NewCard(card), c.createCardReturn.Err
}

// CreateChecklist -
func (c *MockClient) CreateChecklist(card ICard, name string, items []string, args trello.Arguments) (IChecklist, error) {
	c.CreateChecklistCalls = append(c.CreateChecklistCalls, CreateChecklistCall{
		Card:  card,
		Name:  name,
		Items: items,
		Args:  args,
	})
	return c.createChecklistReturn.Checklist, c.createChecklistReturn.Err
}

// NewMockClient - Use this to create a mock client
func NewMockClient(c *MockClientConfig) *MockClient {
	return &MockClient{
		getBoardReturn:        c.GetBoardReturn,
		createCardReturn:      c.CreateCardReturn,
		createChecklistReturn: c.CreateChecklistReturn,
	}
}
