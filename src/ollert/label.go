package ollert

import (
	"github.com/adlio/trello"
)

// ILabel -
type ILabel interface {
	ID() string
	Name() string
}

// Label -
type Label struct {
	label *trello.Label
}

// ID - Label ID
func (l *Label) ID() string {
	return l.label.ID
}

// Name - Label Name
func (l *Label) Name() string {
	return l.label.Name
}

// NewLabel - Use this to get a new Trello label
func NewLabel(label *trello.Label) ILabel {
	return &Label{label: label}
}
