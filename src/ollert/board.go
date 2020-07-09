package ollert

import "github.com/adlio/trello"

// IBoard - Wrapper for trello.Board
type IBoard interface {
	ID() string
	Name() string
	GetLabels(args trello.Arguments) (labels []ILabel, err error)
}

// Board - Use this to take action on a Trello board
type Board struct {
	board *trello.Board
}

// ID -
func (b *Board) ID() string {
	return b.board.ID
}

// Name -
func (b *Board) Name() string {
	return b.board.Name
}

// GetLabels - Fetch labels for this board
func (b *Board) GetLabels(args trello.Arguments) ([]ILabel, error) {
	trelloLabels, err := b.board.GetLabels(args)

	if err != nil {
		return nil, err
	}

	labels := make([]ILabel, len(trelloLabels))

	for i, l := range trelloLabels {
		labels[i] = NewLabel(l)
	}

	return labels, nil
}

// NewBoard - Use this to get a new Trello board
func NewBoard(board *trello.Board) IBoard {
	return &Board{board: board}
}

type (
	// MockBoardConfig -
	MockBoardConfig struct {
		ID              string
		Name            string
		GetLabelsReturn GetLabelsReturn
	}

	// MockBoard -
	MockBoard struct {
		id              string
		name            string
		GetLabelsCalls  []GetLabelsCall
		getLabelsReturn GetLabelsReturn
	}

	// GetLabelsCall -
	GetLabelsCall struct {
		Args trello.Arguments
	}

	// GetLabelsReturn -
	GetLabelsReturn struct {
		Labels []ILabel
		Err    error
	}
)

// ID -
func (b *MockBoard) ID() string {
	return b.id
}

// Name -
func (b *MockBoard) Name() string {
	return b.name
}

// GetLabels -
func (b *MockBoard) GetLabels(args trello.Arguments) ([]ILabel, error) {
	b.GetLabelsCalls = append(b.GetLabelsCalls, GetLabelsCall{Args: args})
	return b.getLabelsReturn.Labels, b.getLabelsReturn.Err
}

// NewMockBoard - Use this to get a mock Board
func NewMockBoard(c *MockBoardConfig) *MockBoard {
	return &MockBoard{
		id:              c.ID,
		name:            c.Name,
		getLabelsReturn: c.GetLabelsReturn,
	}
}
