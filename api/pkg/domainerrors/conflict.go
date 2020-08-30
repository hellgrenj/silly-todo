package domainerrors

import "fmt"

// ConflictError is a domain error for when there is a conflict hindering the current operation from completing
type ConflictError struct {
	expl string
	err  error
}

func (c ConflictError) Error() string {
	return fmt.Sprintf("%s, %v", c.expl, c.err)
}

// ExternalError prints an error safe for external use (in an API response for example)
func (c ConflictError) ExternalError() string {
	return fmt.Sprintf("%s", c.expl)
}

//NewConflictError constructs a new ConflictError with an explanation
func NewConflictError(expl string, e error) ConflictError {
	return ConflictError{expl, e}
}
