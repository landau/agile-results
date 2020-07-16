package app

import (
	"landau/agile-results/src/ollert"
	"landau/agile-results/src/prompt"
	"reflect"
	"testing"

	"github.com/adlio/trello"
	"github.com/pkg/errors"
)

func Test_filter(t *testing.T) {
	labels := []ollert.ILabel{
		ollert.NewLabel(&trello.Label{Name: "foo"}),
		ollert.NewLabel(&trello.Label{Name: "bar"}),
	}

	type args struct {
		labels []ollert.ILabel
		f      func(l ollert.ILabel) bool
	}

	tests := []struct {
		name string
		args args
		want []ollert.ILabel
	}{
		{
			"Returns matching labels",
			args{labels: labels, f: func(l ollert.ILabel) bool { return l.Name() == labels[0].Name() }},
			[]ollert.ILabel{labels[0]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filter(tt.args.labels, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mapLabelsToString(t *testing.T) {
	ids := []string{"foo", "bar"}
	labels := []ollert.ILabel{
		ollert.NewLabel(&trello.Label{ID: ids[0]}),
		ollert.NewLabel(&trello.Label{ID: ids[1]}),
	}

	type args struct {
		labels []ollert.ILabel
		f      func(l ollert.ILabel) string
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"Maps Label.ID to a []string",
			args{labels: labels, f: func(l ollert.ILabel) string { return l.ID() }},
			ids,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapLabelsToString(tt.args.labels, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapLabelsToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterLabelsFromLabelNames(t *testing.T) {
	names := []string{"foo"}
	labels := []ollert.ILabel{
		ollert.NewLabel(&trello.Label{Name: names[0]}),
		ollert.NewLabel(&trello.Label{Name: "bar"}),
	}

	type args struct {
		names  []string
		labels []ollert.ILabel
	}

	tests := []struct {
		name string
		args args
		want []ollert.ILabel
	}{
		{
			"Filters Labels based on a []string of Label.Name",
			args{names: names, labels: labels},
			[]ollert.ILabel{labels[0]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterLabelsFromLabelNames(tt.args.names, tt.args.labels)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterLabelsFromLabelNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelectLabelIDs(t *testing.T) {
	labels := []ollert.ILabel{
		ollert.NewLabel(&trello.Label{ID: "id", Name: "name"}),
	}

	t.Run("Returns a list of selected labels", func(t *testing.T) {
		prompter := &prompt.MockPrompter{
			PromptReturnValue: prompt.MockPrompterPromptReturnValue{
				S: labels[0].Name(), Err: nil,
			},
		}

		got, err := SelectLabelIDs(labels, prompter)

		if err != nil {
			t.Errorf("SelectLabelIDs() error = %v, wantErr %v", err, false)
			return
		}

		want := []string{labels[0].ID()}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("SelectLabelIDs() = %v, want %v", got, want)
		}
	})

	t.Run("Returns a list of selected labels", func(t *testing.T) {
		expectedErr := errors.Errorf("Test error")
		prompter := &prompt.MockPrompter{
			PromptReturnValue: prompt.MockPrompterPromptReturnValue{
				S: "", Err: expectedErr,
			},
		}
		_, err := SelectLabelIDs(labels, prompter)

		if err != expectedErr {
			t.Errorf("SelectLabelIDs() error = %v, wantErr %v", err, expectedErr)
			return
		}
	})
}
