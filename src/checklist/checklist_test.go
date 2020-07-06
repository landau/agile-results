package checklist

import (
	"errors"
	"reflect"
	"testing"

	"github.com/adlio/trello"
)

type (
	MockChecklistHandler struct {
		checklistCalls       []MockChecklistCall
		checkItemCalls       []MockCheckItemCall
		checkListReturnValue MockCheckListReturnValue
	}

	MockChecklistCall struct {
		card      *trello.Card
		name      string
		extraArgs trello.Arguments
	}

	MockCheckListReturnValue struct {
		checklist *trello.Checklist
		err       error
	}

	MockCheckItemCall struct {
		checklist *trello.Checklist
		name      string
		extraArgs trello.Arguments
	}
)

func (m *MockChecklistHandler) CreateChecklist(card *trello.Card, name string, extraArgs trello.Arguments) (*trello.Checklist, error) {
	m.checklistCalls = append(m.checklistCalls, MockChecklistCall{card: card, name: name, extraArgs: extraArgs})
	return m.checkListReturnValue.checklist, m.checkListReturnValue.err
}

func (m *MockChecklistHandler) CreateCheckItem(checklist *trello.Checklist, name string, extraArgs trello.Arguments) (*trello.CheckItem, error) {
	m.checkItemCalls = append(m.checkItemCalls, MockCheckItemCall{checklist: checklist, name: name, extraArgs: extraArgs})
	// I don't care about the return value here because it doesn't need assertion
	return &trello.CheckItem{}, nil
}

func TestTrelloCreator_Create(t *testing.T) {
	type fields struct {
		handler trelloChecklistHandler
	}

	type args struct {
		card  *trello.Card
		items []string
	}

	checklist := &trello.Checklist{}
	card := &trello.Card{Name: "Test Card"}
	items := []string{"hi", "bye"}

	t.Run("Successfully creates a checklist with items", func(t *testing.T) {
		handler := &MockChecklistHandler{
			checkListReturnValue: MockCheckListReturnValue{checklist: checklist, err: nil},
		}
		c := &TrelloCreator{handler: handler}
		got, _ := c.Create(card, items)

		if !reflect.DeepEqual(got, checklist) {
			t.Errorf("TrelloCreator.Create() = %v, want %v", got, checklist)
		}
	})

	t.Run("Calls handler.CreateChecklist once", func(t *testing.T) {
		handler := &MockChecklistHandler{
			checkListReturnValue: MockCheckListReturnValue{checklist: checklist, err: nil},
		}
		c := &TrelloCreator{handler: handler}
		c.Create(card, items)

		if len(handler.checklistCalls) != 1 {
			t.Errorf("Got %d CreateChecklist calls, want %d", len(handler.checklistCalls), 1)
		}
	})

	t.Run("Calls handler.CreateCheckItem for each item", func(t *testing.T) {
		handler := &MockChecklistHandler{
			checkListReturnValue: MockCheckListReturnValue{checklist: checklist, err: nil},
		}
		c := &TrelloCreator{handler: handler}
		c.Create(card, items)

		if len(handler.checkItemCalls) != len(items) {
			t.Errorf("Got %d CreateCheckItem calls, want %d", len(handler.checkItemCalls), len(items))
		}

		for i, item := range items {
			call := handler.checkItemCalls[i]

			if call.checklist != checklist {
				t.Errorf("CreateCheckItem() called with = %v, want %v", call.checklist, checklist)
			}

			if call.name != item {
				t.Errorf("CreateCheckItem() called with %v, want %v", call.name, item)
			}
		}
	})

	t.Run("Returns an error if handler.CreateChecklist fails", func(t *testing.T) {
		expectedErr := errors.New("test error")

		handler := &MockChecklistHandler{
			checkListReturnValue: MockCheckListReturnValue{checklist: nil, err: expectedErr},
		}
		c := &TrelloCreator{handler: handler}
		_, err := c.Create(card, items)

		if err != expectedErr {
			t.Errorf("TrelloCreator.Create() error = %v, wantErr %v", err, expectedErr)
			return
		}
	})
}
