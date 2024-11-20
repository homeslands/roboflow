package xerrors

type NotFound interface {
	error
	IsNotFound()
}

type NotFoundError struct {
	*XError
}

// ThrowNotFound creates a new NotFoundError.
// Default error code is "not_found".
func ThrowNotFound(parent error, message string, options ...Option) error {
	e := NotFoundError{CreateXError(parent, "not_found", message)}
	for _, opt := range options {
		opt(e.XError)
	}
	return &e
}

func (err *NotFoundError) IsNotFound() {}

func IsNotFound(err error) bool {
	_, ok := err.(NotFound)
	return ok
}

func (err *NotFoundError) Is(target error) bool {
	t, ok := target.(*NotFoundError)
	if !ok {
		return false
	}
	return err.XError.Is(t.XError)
}

func (err *NotFoundError) Unwrap() error {
	return err.XError
}
