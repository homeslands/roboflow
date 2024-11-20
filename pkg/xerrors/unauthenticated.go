package xerrors

var (
	_ Unauthenticated = (*UnauthenticatedError)(nil)
	_ Error           = (*UnauthenticatedError)(nil)
)

type Unauthenticated interface {
	error
	IsUnauthenticated()
}

type UnauthenticatedError struct {
	*XError
}

// ThrowUnauthenticated creates a new UnauthenticatedError.
// Default error code is "unauthenticated".
func ThrowUnauthenticated(parent error, message string, options ...Option) error {
	e := UnauthenticatedError{CreateXError(parent, "unauthenticated", message)}
	for _, opt := range options {
		opt(e.XError)
	}
	return &e
}

func (err *UnauthenticatedError) IsUnauthenticated() {}

func IsUnauthenticated(err error) bool {
	_, ok := err.(Unauthenticated)
	return ok
}

func (err *UnauthenticatedError) Is(target error) bool {
	t, ok := target.(*UnauthenticatedError)
	if !ok {
		return false
	}
	return err.XError.Is(t.XError)
}

func (err *UnauthenticatedError) Unwrap() error {
	return err.XError
}
