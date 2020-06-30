package app

import (
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

type MockPrompter struct {
	calls       []string
	ReturnValue string
}

func (p *MockPrompter) Prompt(s string) (string, error) {
	p.calls = append(p.calls, s)
	return p.ReturnValue, nil
}

type MockCardCreatorCall struct {
	card *trello.Card
	args trello.Arguments
}

type MockCardCreator struct {
	calls       []MockCardCreatorCall
	ReturnValue error
}

func (c *MockCardCreator) CreateCard(card *trello.Card, extraArgs trello.Arguments) error {
	c.calls = append(c.calls, MockCardCreatorCall{card: card, args: extraArgs})
	return nil
}

func TestRunApp(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	listID := "1234"
	prompter := &MockPrompter{ReturnValue: "Foo"}
	cardCreator := &MockCardCreator{}

	config := &Config{
		Logrus:      logger,
		Prompter:    prompter,
		ListID:      listID,
		CardCreator: cardCreator,
	}

	t.Run("Successfully creates a card", func(t *testing.T) {
		RunApp(config)

		assertCallCount(
			t,
			len(cardCreator.calls),
			1,
			"CardCreatore.CreateCard should only be called once",
		)

		assertCallCount(
			t,
			len(prompter.calls),
			1,
			"Prompter.prompt should only be called once",
		)

		if prompter.calls[0] != "Card Name: " {
			t.Error("Expected Prompter.Prompt to recevie 'Card Name: '")
		}

		if prompter.ReturnValue != cardCreator.calls[0].card.Name {
			t.Errorf(
				"s = %v, want %v",
				prompter.ReturnValue,
				cardCreator.calls[0].card.Name,
			)
		}
	})
}

func assertCallCount(t *testing.T, got int, want int, s string) {
	if got != want {
		t.Errorf("%s = %v, want %v", s, got, want)
	}
}
