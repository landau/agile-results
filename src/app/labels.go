package app

import (
	"fmt"
	"landau/agile-results/src/prompt"
	"landau/agile-results/src/utils"
	"strings"

	"github.com/adlio/trello"
)

func filter(labels []*trello.Label, f func(l *trello.Label) bool) []*trello.Label {
	filtered := make([]*trello.Label, 0)

	for _, v := range labels {
		if f(v) {
			filtered = append(filtered, v)
		}
	}

	return filtered
}

func mapLabelsToString(labels []*trello.Label, f func(l *trello.Label) string) []string {
	mapped := make([]string, len(labels))

	for i, label := range labels {
		mapped[i] = f(label)
	}

	return mapped
}

func filterLabelsFromLabelNames(names []string, labels []*trello.Label) []*trello.Label {
	return filter(labels, func(l *trello.Label) bool {
		return utils.IndexOf(len(names), func(i int) bool {
			return names[i] == l.Name
		}) > -1
	})
}

// SelectLabelIDs -
func SelectLabelIDs(labels []*trello.Label, prompter prompt.Prompter) ([]string, error) {
	labelNames := mapLabelsToString(labels, func(l *trello.Label) string { return l.Name })

	selected, err := prompter.Prompt(
		fmt.Sprintf("Selected %d labels (%v): ", len(labelNames), strings.Join(labelNames, ", ")),
	)

	if err != nil {
		return make([]string, 0), err
	}

	return mapLabelsToString(
		filterLabelsFromLabelNames(strings.Split(selected, ","), labels),
		func(l *trello.Label) string { return l.ID },
	), nil
}
