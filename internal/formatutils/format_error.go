package formatutils

// FormatError is a custom error type for formatting errors
type FormatError struct {
	Msg string
}

func (e *FormatError) Error() string {
	return e.Msg
}
