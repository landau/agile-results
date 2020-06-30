package app

import (
	"fmt"
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

type PrompterReturnValue struct {
	s   string
	err error
}

type MockPrompter struct {
	calls []string
	// TODO: How do I nest this and use in the code?????
	ReturnValue PrompterReturnValue
}

func (p *MockPrompter) Prompt(s string) (string, error) {
	p.calls = append(p.calls, s)
	return p.ReturnValue.s, p.ReturnValue.err
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
	prompter := &MockPrompter{ReturnValue: PrompterReturnValue{s: "foo,bar", err: nil}}
	cardCreator := &MockCardCreator{}
	labels := []*trello.Label{{ID: "abc", Name: "foo"}, {ID: "efg", Name: "bar"}}

	config := &Config{
		Logrus:      logger,
		Prompter:    prompter,
		ListID:      listID,
		CardCreator: cardCreator,
		Labels:      labels,
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
			2,
			"Prompter.prompt should be called twice",
		)

		if prompter.calls[0] != "Card Name: " {
			t.Error("Expected Prompter.Prompt to recevie 'Card Name: '")
		}

		prompt2 := fmt.Sprintf("Selected 2 labels (%v, %v): ", labels[0].Name, labels[1].Name)
		if prompter.calls[1] != prompt2 {
			t.Errorf("Expected Prompter.Prompt to recevie '%s', got '%s", prompt2, prompter.calls[1])
		}

		if prompter.ReturnValue.s != cardCreator.calls[0].card.Name {
			t.Errorf(
				"Card name should be = %s, but got %s",
				prompter.ReturnValue.s,
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
