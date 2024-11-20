package xerrors

var (
	_ Internal = (*InternalError)(nil)
	_ Error    = (*InternalError)(nil)
)

type Internal interface {
	error
	IsInternal()
}

type InternalError struct {
	*XError
}

// ThrowInternal creates a new InternalError.
// Default error code is "internal_server_error".
func ThrowInternal(parent error, message string, options ...Option) error {
	e := InternalError{CreateXError(parent, "internal_server_error", message)}
	for _, opt := range options {
		opt(e.XError)
	}
	return &e
}

func (err *InternalError) IsInternal() {}

func IsInternal(err error) bool {
	_, ok := err.(Internal)
	return ok
}

func (err *InternalError) Is(target error) bool {
	t, ok := target.(*InternalError)
	if !ok {
		return false
	}
	return err.XError.Is(t.XError)
}

func (err *InternalError) Unwrap() error {
	return err.XError
}
