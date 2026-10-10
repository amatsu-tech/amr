package output

type Hinter interface {
	Hint() string
}

type hintedError struct {
	err  error
	hint string
}

var _ Hinter = (*hintedError)(nil)

func WithHint(err error, hint string) error {
	return &hintedError{err: err, hint: hint}
}

func (e *hintedError) Error() string { return e.err.Error() }
func (e *hintedError) Unwrap() error { return e.err }
func (e *hintedError) Hint() string  { return e.hint }
