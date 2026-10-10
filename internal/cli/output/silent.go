package output

type silentError struct {
	err error
}

func Silent(err error) error {
	return &silentError{err: err}
}

func (e *silentError) Error() string { return e.err.Error() }
func (e *silentError) Unwrap() error { return e.err }
