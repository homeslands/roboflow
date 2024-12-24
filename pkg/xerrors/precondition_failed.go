package xerrors

var (
	_ PreconditionFailed = (*PreconditionFailedError)(nil)
	_ Error              = (*PreconditionFailedError)(nil)
)

type PreconditionFailed interface {
	error
	IsPreconditionFailed()
}

type PreconditionFailedError struct {
	*XError
}

func ThrowPreconditionFailedX(parent error, id, message string) error {
	return &PreconditionFailedError{CreateXError(parent, id, message)}
}

// ThrowPreconditionFailed creates a new PreconditionFailedError.
// Default error code is "precondition_failed".
func ThrowPreconditionFailed(parent error, message string, options ...Option) error {
	e := PreconditionFailedError{CreateXError(parent, "precondition_failed", message)}
	for _, opt := range options {
		opt(e.XError)
	}
	return &e
}

func (err *PreconditionFailedError) IsPreconditionFailed() {}

func IsPreconditionFailed(err error) bool {
	_, ok := err.(PreconditionFailed)
	return ok
}

func (err *PreconditionFailedError) Is(target error) bool {
	t, ok := target.(*PreconditionFailedError)
	if !ok {
		return false
	}
	return err.XError.Is(t.XError)
}

func (err *PreconditionFailedError) Unwrap() error {
	return err.XError
}
