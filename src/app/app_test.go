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
	returnValue PrompterReturnValue
}

func (p *MockPrompter) Prompt(s string) (string, error) {
	p.calls = append(p.calls, s)
	return p.returnValue.s, p.returnValue.err
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
	prompter := &MockPrompter{returnValue: PrompterReturnValue{s: "foo,bar", err: nil}}
	cardCreator := &MockCardCreator{}

	labels := []*trello.Label{{ID: "abc", Name: "foo"}, {ID: "efg", Name: "bar"}}
	labelFetcher := &MockLabelFetcher{returnValue: LabelFetcherReturnValue{labels: labels}}

	config := &Config{
		Logrus:       logger,
		Prompter:     prompter,
		ListID:       listID,
		CardCreator:  cardCreator,
		LabelFetcher: labelFetcher,
	}

	t.Run("Successfully creates a card with labels", func(t *testing.T) {
		RunApp(config)

		assertCallCount(
			t,
			len(labelFetcher.calls),
			1,
			"LabelFetcher.Fetch should only be called once",
		)

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

		if prompter.returnValue.s != cardCreator.calls[0].card.Name {
			t.Errorf(
				"Card name should be = %s, but got %s",
				prompter.returnValue.s,
				cardCreator.calls[0].card.Name,
			)
		}

		// TODO: use deepequal on entire list here instead
		if cardCreator.calls[0].card.IDLabels[0] != labels[0].ID {
			t.Errorf(
				"Card should have label = %s, but got %s",
				labels[0].ID,
				cardCreator.calls[0].card.IDLabels[0],
			)
		}

		if cardCreator.calls[0].card.IDLabels[1] != labels[1].ID {
			t.Errorf(
				"Card should have label = %s, but got %s",
				labels[1].ID,
				cardCreator.calls[0].card.IDLabels[1],
			)
		}
	})
}

func assertCallCount(t *testing.T, got int, want int, s string) {
	if got != want {
		t.Errorf("%s = %v, want %v", s, got, want)
	}
}
