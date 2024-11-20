package xerrors

var (
	_ Unavailable = (*UnavailableError)(nil)
	_ Error       = (*UnavailableError)(nil)
)

type Unavailable interface {
	error
	IsUnavailable()
}

type UnavailableError struct {
	*XError
}

// ThrowUnavailable creates a new UnavailableError.
// Default error code is "unavailable".
func ThrowUnavailable(parent error, message string, options ...Option) error {
	e := UnavailableError{CreateXError(parent, "unavailable", message)}
	for _, opt := range options {
		opt(e.XError)
	}
	return &e
}

func (err *UnavailableError) IsUnavailable() {}

func IsUnavailable(err error) bool {
	_, ok := err.(Unavailable)
	return ok
}

func (err *UnavailableError) Is(target error) bool {
	t, ok := target.(*UnavailableError)
	if !ok {
		return false
	}
	return err.XError.Is(t.XError)
}

func (err *UnavailableError) Unwrap() error {
	return err.XError
}
