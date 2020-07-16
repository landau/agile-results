package app

import (
	"fmt"
	"landau/agile-results/src/ollert"
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

func TestRunApp(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	listID := "1234"
	labels := []ollert.ILabel{
		ollert.NewLabel(&trello.Label{ID: "id1", Name: "name1"}),
		ollert.NewLabel(&trello.Label{ID: "id2", Name: "name2"}),
	}

	t.Run("Successfully creates a card with labels and checklist", func(t *testing.T) {
		prompter := &prompt.MockPrompter{
			PromptReturnValue: prompt.MockPrompterPromptReturnValue{S: "name1,name2", Err: nil},
		}

		mockBoard := ollert.NewMockBoard(&ollert.MockBoardConfig{
			ID:              "id",
			Name:            "name",
			GetLabelsReturn: ollert.GetLabelsReturn{Labels: labels},
		})

		clientConfig := &ollert.MockClientConfig{
			GetBoardReturn: ollert.GetBoardReturn{
				Board: mockBoard,
			},
			CreateChecklistReturn: ollert.CreateChecklistReturn{
				Checklist: ollert.NewCheckList(&trello.Checklist{}),
			},
		}
		client := ollert.NewMockClient(clientConfig)

		appConfig := &Config{
			Client:   client,
			ListID:   listID,
			Logrus:   logger,
			Prompter: prompter,
		}

		card, err := RunApp(appConfig)

		if err != nil {
			t.Errorf("An unexpected error occured: %v", err)
			return
		}

		assertCallCount(t, len(client.CreateCardCalls), 1, "client.CreateCardCalls call count")

		if card.ID() != client.CreateCardCalls[0].Card.ID {
			t.Errorf(
				"Expected card ID %s, but got %s",
				client.CreateCardCalls[0].Card.ID,
				card.ID(),
			)
		}

		// Validate prompter ---
		assertCallCount(t, len(prompter.PromptCalls), 2, "Prompter.Prompt call count")
		assertCallCount(t, len(prompter.PromptListCalls), 0, "Prompter.PromptList call count")

		if prompter.PromptCalls[0] != "Card Name: " {
			t.Error("Expected Prompter.Prompt to recevie 'Card Name: '")
		}

		prompt2 := fmt.Sprintf(
			"Selected 2 labels (%v, %v): ", labels[0].Name(), labels[1].Name(),
		)
		if prompter.PromptCalls[1] != prompt2 {
			t.Errorf("Expected Prompter.Prompt to recevie '%s', got '%s", prompt2, prompter.PromptCalls[1])
		}

		if prompter.PromptReturnValue.S != card.Name() {
			t.Errorf(
				"Card name should be = %s, but got %s",
				prompter.PromptReturnValue.S,
				card.Name(),
			)
		}

		// Validate labels ---
		assertCallCount(t, len(mockBoard.GetLabelsCalls), 1, "board.GetLabelCalls call count")

		if len(card.IDLabels()) < 1 {
			t.Error("Expected card to have > 0 IDLabels")
		}

		if card.IDLabels()[0] != labels[0].ID() {
			t.Errorf(
				"Card should have label = %s, but got %s",
				labels[0].ID(),
				card.IDLabels()[0],
			)
		}

		if card.IDLabels()[1] != labels[1].ID() {
			t.Errorf(
				"Card should have label = %s, but got %s",
				labels[1].ID(),
				card.IDLabels()[1],
			)
		}

		// Validate Checklist ---
		assertCallCount(
			t, len(client.CreateChecklistCalls), 0, "client.CreateChecklistCalls call count",
		)
	})

	t.Run("Creates a card with a checklist", func(t *testing.T) {
		prompter := &prompt.MockPrompter{
			PromptReturnValue:     prompt.MockPrompterPromptReturnValue{S: "name1,name2", Err: nil},
			PromptListReturnValue: prompt.MockPrompterPromptListReturnValue{Items: []string{"list"}},
		}

		mockBoard := ollert.NewMockBoard(&ollert.MockBoardConfig{
			ID:              "id",
			Name:            "name",
			GetLabelsReturn: ollert.GetLabelsReturn{Labels: labels},
		})

		clientConfig := &ollert.MockClientConfig{
			GetBoardReturn: ollert.GetBoardReturn{
				Board: mockBoard,
			},
			CreateCardReturn: ollert.CreateCardReturn{
				Card: ollert.NewCard(&trello.Card{ID: "id"}),
			},
			CreateChecklistReturn: ollert.CreateChecklistReturn{
				Checklist: ollert.NewCheckList(&trello.Checklist{}),
			},
		}
		client := ollert.NewMockClient(clientConfig)

		appConfig := &Config{
			Client:       client,
			HasChecklist: true,
			ListID:       listID,
			Logrus:       logger,
			Prompter:     prompter,
		}

		card, err := RunApp(appConfig)
		if err != nil {
			t.Errorf("Unexpected Error: %v", err)
		}

		assertCallCount(t, len(prompter.PromptListCalls), 1, "Prompter.PromptList call count")
		assertCallCount(
			t, len(client.CreateChecklistCalls), 1, "client.CreateChecklistCalls call count",
		)

		if client.CreateChecklistCalls[0].Card.ID() != card.ID() {
			t.Errorf(
				"Checklist should have been called with %v, but got %v",
				card.ID(),
				client.CreateChecklistCalls[0].Card.ID(),
			)
		}

		items := client.CreateChecklistCalls[0].Items
		if !reflect.DeepEqual(items, []string{"list"}) {
			t.Errorf(
				"Checklist should have been called with %v, but got %v",
				[]string{"list"},
				items,
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
