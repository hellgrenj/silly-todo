package validation

// ErrMissingField is a custom validation error
type ErrMissingField string

func (e ErrMissingField) Error() string {
	return string(e) + " is required"
}
