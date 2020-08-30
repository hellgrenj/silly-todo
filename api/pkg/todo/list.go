package todo

import "github.com/hellgrenj/silly-todo/pkg/validation"

// List is the struct for a todo list
type List struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Items []Item `json:"items"`
}

// OK is the validation function for the struct List
func (l *List) OK() error {
	if len(l.Name) == 0 {
		return validation.ErrMissingField("Name")
	}
	return nil
}
