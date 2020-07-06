package app

import (
	"fmt"
	"landau/agile-results/src/checklist"
	"landau/agile-results/src/prompt"
	"reflect"
	"testing"

	"github.com/adlio/trello"
	"github.com/sirupsen/logrus"
)

func Test_createCurlDelCmd(t *testing.T) {
	type args struct {
		cardID string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Prints a curl DELETE command with an embedded card id",
			args{cardID: "1234"},
			"curl -sXDELETE \"https://api.trello.com/1/cards/1234?key=$TRELLO_API_KEY&token=$TRELLO_API_TOKEN\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createCurlDelCmd(tt.args.cardID); got != tt.want {
				t.Errorf("createCurlDelCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

type MockCardCreatorCall struct {
	card *trello.Card
	args trello.Arguments
}

type MockCardCreator struct {
	calls       []MockCardCreatorCall
	returnValue error
}

func (c *MockCardCreator) CreateCard(card *trello.Card, extraArgs trello.Arguments) error {
	c.calls = append(c.calls, MockCardCreatorCall{card: card, args: extraArgs})
	return nil
}

type LabelFetcherReturnValue struct {
	labels []*trello.Label
	err    error
}

type MockLabelFetcher struct {
	calls       []trello.Arguments
	returnValue LabelFetcherReturnValue
}

func (l *MockLabelFetcher) Fetch(args trello.Arguments) (labels []*trello.Label, err error) {
	l.calls = append(l.calls, args)
	return l.returnValue.labels, l.returnValue.err
}

func TestRunApp(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	listID := "1234"
	labels := []*trello.Label{{ID: "abc", Name: "foo"}, {ID: "efg", Name: "bar"}}

	t.Run("Successfully creates a card with labels and checklist", func(t *testing.T) {
		prompter := &prompt.MockPrompter{
			PromptReturnValue: prompt.MockPrompterPromptReturnValue{S: "foo,bar", Err: nil},
		}
		cardCreator := &MockCardCreator{}
		labelFetcher := &MockLabelFetcher{returnValue: LabelFetcherReturnValue{labels: labels}}

		checklistCreator := &checklist.MockCreator{
			CreateReturnValue: checklist.MockCreateReturnValue{},
		}

		config := &Config{
			CardCreator:      cardCreator,
			ChecklistCreator: checklistCreator,
			LabelFetcher:     labelFetcher,
			ListID:           listID,
			Logrus:           logger,
			Prompter:         prompter,
		}

		RunApp(config)

		assertCallCount(t, len(labelFetcher.calls), 1, "LabelFetcher.Fetch call count")
		assertCallCount(t, len(cardCreator.calls), 1, "CardCreatore.CreateCard call count")
		assertCallCount(t, len(prompter.PromptCalls), 2, "Prompter.Prompt call count")
		assertCallCount(t, len(prompter.PromptListCalls), 0, "Prompter.PromptList call count")

		if prompter.PromptCalls[0] != "Card Name: " {
			t.Error("Expected Prompter.Prompt to recevie 'Card Name: '")
		}

		prompt2 := fmt.Sprintf("Selected 2 labels (%v, %v): ", labels[0].Name, labels[1].Name)
		if prompter.PromptCalls[1] != prompt2 {
			t.Errorf("Expected Prompter.Prompt to recevie '%s', got '%s", prompt2, prompter.PromptCalls[1])
		}

		if prompter.PromptReturnValue.S != cardCreator.calls[0].card.Name {
			t.Errorf(
				"Card name should be = %s, but got %s",
				prompter.PromptReturnValue.S,
				cardCreator.calls[0].card.Name,
			)
		}

		card := cardCreator.calls[0].card

		// TODO: use deepequal on entire list here instead
		if card.IDLabels[0] != labels[0].ID {
			t.Errorf(
				"Card should have label = %s, but got %s",
				labels[0].ID,
				cardCreator.calls[0].card.IDLabels[0],
			)
		}

		if card.IDLabels[1] != labels[1].ID {
			t.Errorf(
				"Card should have label = %s, but got %s",
				labels[1].ID,
				cardCreator.calls[0].card.IDLabels[1],
			)
		}

		assertCallCount(t, len(checklistCreator.Calls), 0, "CheckListCreator.Creator call count")
	})

	t.Run("Creates a card with a checklist", func(t *testing.T) {
		prompter := &prompt.MockPrompter{
			PromptReturnValue:     prompt.MockPrompterPromptReturnValue{S: "foo,bar", Err: nil},
			PromptListReturnValue: prompt.MockPrompterPromptListReturnValue{Items: []string{"list"}},
		}
		cardCreator := &MockCardCreator{}
		labelFetcher := &MockLabelFetcher{returnValue: LabelFetcherReturnValue{labels: labels}}

		checklistCreator := &checklist.MockCreator{
			CreateReturnValue: checklist.MockCreateReturnValue{},
		}

		config := &Config{
			CardCreator:      cardCreator,
			ChecklistCreator: checklistCreator,
			HasChecklist:     true,
			LabelFetcher:     labelFetcher,
			ListID:           listID,
			Logrus:           logger,
			Prompter:         prompter,
		}

		RunApp(config)

		card := cardCreator.calls[0].card

		assertCallCount(t, len(prompter.PromptListCalls), 1, "Prompter.PromptList call count")
		assertCallCount(t, len(checklistCreator.Calls), 1, "CheckListCreator.Creator call count")

		if checklistCreator.Calls[0].Card != card {
			t.Errorf(
				"Checklist should have been called with %v, but got %v",
				card,
				checklistCreator.Calls[0].Card,
			)
		}

		if !reflect.DeepEqual(checklistCreator.Calls[0].Items, []string{"list"}) {
			t.Errorf(
				"Checklist should have been called with %v, but got %v",
				[]string{"list"},
				checklistCreator.Calls[0].Items,
			)
		}

	})
}

func assertCallCount(t *testing.T, got int, want int, s string) {
	t.Helper()

	if got != want {
		t.Errorf("%s = %v, want %v", s, got, want)
	}
}
