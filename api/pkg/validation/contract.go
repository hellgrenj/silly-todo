package validation

// Ok represents types capable of validating
// themselves.
type Ok interface {
	OK() error
}
