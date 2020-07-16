package app

import (
	"fmt"
	"landau/agile-results/src/ollert"
	"landau/agile-results/src/prompt"
	"landau/agile-results/src/utils"
	"strings"
)

func filter(labels []ollert.ILabel, f func(l ollert.ILabel) bool) []ollert.ILabel {
	filtered := make([]ollert.ILabel, 0)

	for _, v := range labels {
		if f(v) {
			filtered = append(filtered, v)
		}
	}

	return filtered
}

func mapLabelsToString(labels []ollert.ILabel, f func(l ollert.ILabel) string) []string {
	mapped := make([]string, len(labels))

	for i, label := range labels {
		mapped[i] = f(label)
	}

	return mapped
}

func filterLabelsFromLabelNames(names []string, labels []ollert.ILabel) []ollert.ILabel {
	return filter(labels, func(l ollert.ILabel) bool {
		return utils.IndexOf(len(names), func(i int) bool {
			return names[i] == l.Name()
		}) > -1
	})
}

// SelectLabelIDs -
func SelectLabelIDs(labels []ollert.ILabel, prompter prompt.Prompter) ([]string, error) {
	labelNames := mapLabelsToString(labels, func(l ollert.ILabel) string { return l.Name() })

	selected, err := prompter.Prompt(
		fmt.Sprintf("Selected %d labels (%v): ", len(labelNames), strings.Join(labelNames, ", ")),
	)

	if err != nil {
		return make([]string, 0), err
	}

	return mapLabelsToString(
		filterLabelsFromLabelNames(strings.Split(selected, ","), labels),
		func(l ollert.ILabel) string { return l.ID() },
	), nil
}
