package ollert

import "github.com/adlio/trello"

// IChecklist -
type IChecklist interface {
	ID() string
	CheckItems() []ICheckItem
}

// Checklist -
type Checklist struct {
	checklist  *trello.Checklist
	checkItems []ICheckItem
}

// ID -
func (c *Checklist) ID() string {
	return c.checklist.ID
}

// CheckItems -
func (c *Checklist) CheckItems() []ICheckItem {
	return c.checkItems
}

// NewCheckList -
func NewCheckList(checklist *trello.Checklist) IChecklist {
	c := &Checklist{checklist: checklist}

	for _, item := range checklist.CheckItems {
		c.checkItems = append(c.checkItems, NewCheckItem(&item))
	}

	return c
}

// ICheckItem -
type ICheckItem interface {
	ID() string
}

// CheckItem -
type CheckItem struct {
	checkItem *trello.CheckItem
}

// ID -
func (c *CheckItem) ID() string {
	return c.checkItem.ID
}

// NewCheckItem -
func NewCheckItem(checkItem *trello.CheckItem) ICheckItem {
	return &CheckItem{checkItem: checkItem}
}
