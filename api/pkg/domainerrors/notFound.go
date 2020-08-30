package domainerrors

import "fmt"

// NotFoundError is a domain error for when an item can not be found
type NotFoundError struct {
	itemName string
	err      error
}

func (n NotFoundError) Error() string {
	return fmt.Sprintf("could not find %s, %v", n.itemName, n.err)
}

// ExternalError prints an error safe for external use (in an API response for example)
func (n NotFoundError) ExternalError() string {
	return fmt.Sprintf("could not find %s", n.itemName)
}

//NewNotFoundError constructs a NotFoundError with an itemName and the original error
func NewNotFoundError(itemName string, e error) NotFoundError {
	return NotFoundError{itemName, e}
}
