package ollert

import "github.com/adlio/trello"

// ITrelloClient - This interface eases testing due to some of the abstractions
// inside the Trello client. For example, a trello.Board is capabable of fetching
// Labels. This makes sense from a domain perspective, but complicates testing
// because of underlying network calls. You could mock the network calls, but
// that violates privacy principle. Another issue is public field access. For
// Example, Card.MoveToBottom() is a useful method to test, but interfaces are
// limited to methods and therefore accessing Card.Name is not possible from
// an interface's POV. This makes me question programming patterns in Go in general.
// Should fields always be private and only expose Set/Get-ters so that testing
// is easier to Mock? TLDR; this interface seems mostly useful for testing
// client.go
type ITrelloClient interface {
	CreateCard(card *trello.Card, args trello.Arguments) error
	GetBoard(boardID string, args trello.Arguments) (*trello.Board, error)
	CreateChecklist(card *trello.Card, name string, args trello.Arguments) (*trello.Checklist, error)
	CreateCheckItem(checklist *trello.Checklist, name string, args trello.Arguments) (*trello.CheckItem, error)
}

// --- Mock Trello Client Implementations

type (
	// MockTrelloClientConfig -
	MockTrelloClientConfig struct {
		GetBoardReturn        TrelloGetBoardReturn
		CreateCardReturn      TrelloCreateCardReturn
		CreateChecklistReturn TrelloCreateChecklistReturn
		CreateCheckItemReturn TrelloCreateCheckItemReturn
	}

	// MockTrelloClient - A trello wrapper client
	MockTrelloClient struct {
		GetBoardCalls  []TrelloGetBoardCall
		getBoardReturn TrelloGetBoardReturn

		CreateCardCalls  []TrelloCreateCardCall
		createCardReturn TrelloCreateCardReturn

		CreateChecklistCalls  []TrelloCreateChecklistCall
		createChecklistReturn TrelloCreateChecklistReturn

		CreateCheckItemCalls  []TrelloCreateCheckItemCall
		createCheckItemReturn TrelloCreateCheckItemReturn
	}

	// TrelloGetBoardCall -
	TrelloGetBoardCall struct {
		BoardID string
		Args    trello.Arguments
	}

	// TrelloGetBoardReturn -
	TrelloGetBoardReturn struct {
		Board *trello.Board
		Err   error
	}

	// TrelloCreateCardCall -
	TrelloCreateCardCall struct {
		Card *trello.Card
		Args trello.Arguments
	}

	// TrelloCreateCardReturn -
	TrelloCreateCardReturn error

	// TrelloCreateChecklistCall -
	TrelloCreateChecklistCall struct {
		Card *trello.Card
		Name string
		Args trello.Arguments
	}

	// TrelloCreateChecklistReturn -
	TrelloCreateChecklistReturn struct {
		Checklist *trello.Checklist
		Err       error
	}

	// TrelloCreateCheckItemCall -
	TrelloCreateCheckItemCall struct {
		Checklist *trello.Checklist
		Name      string
		Args      trello.Arguments
	}

	// TrelloCreateCheckItemReturn -
	TrelloCreateCheckItemReturn struct {
		CheckItem *trello.CheckItem
		Err       error
	}
)

// GetBoard -
func (c *MockTrelloClient) GetBoard(boardID string, args trello.Arguments) (*trello.Board, error) {
	c.GetBoardCalls = append(c.GetBoardCalls, TrelloGetBoardCall{BoardID: boardID, Args: args})
	return c.getBoardReturn.Board, c.getBoardReturn.Err
}

// CreateCard -
func (c *MockTrelloClient) CreateCard(card *trello.Card, args trello.Arguments) error {
	c.CreateCardCalls = append(c.CreateCardCalls, TrelloCreateCardCall{Card: card, Args: args})
	return c.createCardReturn
}

// CreateChecklist -
func (c *MockTrelloClient) CreateChecklist(card *trello.Card, name string, args trello.Arguments) (*trello.Checklist, error) {
	c.CreateChecklistCalls = append(c.CreateChecklistCalls, TrelloCreateChecklistCall{
		Card: card,
		Name: name,
		Args: args,
	})
	return c.createChecklistReturn.Checklist, c.createChecklistReturn.Err
}

// CreateCheckItem -
func (c *MockTrelloClient) CreateCheckItem(checklist *trello.Checklist, name string, args trello.Arguments) (*trello.CheckItem, error) {
	c.CreateCheckItemCalls = append(c.CreateCheckItemCalls, TrelloCreateCheckItemCall{
		Checklist: checklist,
		Name:      name,
		Args:      args,
	})
	return c.createCheckItemReturn.CheckItem, c.createCheckItemReturn.Err
}

// NewMockTrelloClient - Use this to create a mock trello client
func NewMockTrelloClient(c *MockTrelloClientConfig) ITrelloClient {
	return &MockTrelloClient{
		getBoardReturn:        c.GetBoardReturn,
		createCardReturn:      c.CreateCardReturn,
		createChecklistReturn: c.CreateChecklistReturn,
		createCheckItemReturn: c.CreateCheckItemReturn,
	}
}
