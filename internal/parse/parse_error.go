package parse

// ResultError is a custom error type for parsing club names
type ResultError struct {
	Msg string
}

func (e *ResultError) Error() string {
	return e.Msg
}
