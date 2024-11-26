package xerrors

var (
	_ InvalidArgument = (*InvalidArgumentError)(nil)
	_ Error           = (*InvalidArgumentError)(nil)
)

type InvalidArgument interface {
	error
	IsInvalidArgument()
}

type InvalidArgumentError struct {
	*XError
}

// ThrowInvalidArgument creates a new InvalidArgumentError.
// Default error code is "invalid_argument".
func ThrowInvalidArgument(parent error, message string, options ...Option) error {
	e := InvalidArgumentError{CreateXError(parent, "invalid_argument", message)}
	for _, opt := range options {
		opt(e.XError)
	}
	return &e
}

func (err *InvalidArgumentError) IsInvalidArgument() {}

func IsErrorInvalidArgument(err error) bool {
	_, ok := err.(InvalidArgument)
	return ok
}

func (err *InvalidArgumentError) Is(target error) bool {
	t, ok := target.(*InvalidArgumentError)
	if !ok {
		return false
	}
	return err.XError.Is(t.XError)
}

func (err *InvalidArgumentError) As(target any) bool {
	targetErr, ok := target.(*InvalidArgumentError)
	if !ok {
		return false
	}
	*targetErr = *err
	return true
}

func (err *InvalidArgumentError) Unwrap() error {
	return err.XError
}
