package app

import "github.com/adlio/trello"

// BoardFetcher -
type BoardFetcher interface {
	Fetch(args trello.Arguments) (board *trello.Board, err error)
}

// TrelloBoardFetcher -
type TrelloBoardFetcher struct {
	boardID string
	client  *trello.Client
}

// Fetch -
func (t *TrelloBoardFetcher) Fetch(args trello.Arguments) (*trello.Board, error) {
	return t.client.GetBoard(t.boardID, trello.Defaults())
}

// NewBoardFetcher -
func NewBoardFetcher(boardID string, client *trello.Client) *TrelloBoardFetcher {
	return &TrelloBoardFetcher{boardID: boardID, client: client}
}
