package app

import (
	"reflect"
	"testing"

	"github.com/adlio/trello"
	"github.com/pkg/errors"
)

func Test_filter(t *testing.T) {
	labels := []*trello.Label{{Name: "foo"}, {Name: "bar"}}

	type args struct {
		labels []*trello.Label
		f      func(l *trello.Label) bool
	}

	tests := []struct {
		name string
		args args
		want []*trello.Label
	}{
		{
			"Returns matching labels",
			args{labels: labels, f: func(l *trello.Label) bool { return l.Name == labels[0].Name }},
			[]*trello.Label{labels[0]},
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
	labels := []*trello.Label{{ID: ids[0]}, {ID: ids[1]}}

	type args struct {
		labels []*trello.Label
		f      func(l *trello.Label) string
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"Maps Label.ID to a []string",
			args{labels: labels, f: func(l *trello.Label) string { return l.ID }},
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
	labels := []*trello.Label{{Name: names[0]}, {Name: "bar"}}

	type args struct {
		names  []string
		labels []*trello.Label
	}

	tests := []struct {
		name string
		args args
		want []*trello.Label
	}{
		{
			"Filters Labels based on a []string of Label.Name",
			args{names: names, labels: labels},
			[]*trello.Label{labels[0]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterLabelsFromLabelNames(tt.args.names, tt.args.labels); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterLabelsFromLabelNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelectLabelIDs(t *testing.T) {
	labels := []*trello.Label{{ID: "id", Name: "name"}}

	t.Run("Returns a list of selected labels", func(t *testing.T) {
		prompter := &MockPrompter{ReturnValue: PrompterReturnValue{s: labels[0].Name, err: nil}}
		got, err := SelectLabelIDs(labels, prompter)

		if err != nil {
			t.Errorf("SelectLabelIDs() error = %v, wantErr %v", err, false)
			return
		}

		want := []string{labels[0].ID}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("SelectLabelIDs() = %v, want %v", got, want)
		}
	})

	t.Run("Returns a list of selected labels", func(t *testing.T) {
		expectedErr := errors.Errorf("Test error")
		prompter := &MockPrompter{ReturnValue: PrompterReturnValue{s: "", err: expectedErr}}
		_, err := SelectLabelIDs(labels, prompter)

		if err != expectedErr {
			t.Errorf("SelectLabelIDs() error = %v, wantErr %v", err, expectedErr)
			return
		}
	})
}
