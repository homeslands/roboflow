package xerrors

var (
	_ AlreadyExists = (*AlreadyExistsError)(nil)
	_ Error         = (*AlreadyExistsError)(nil)
)

type AlreadyExists interface {
	error
	IsAlreadyExists()
}

type AlreadyExistsError struct {
	*XError
}

// ThrowAlreadyExists creates a new AlreadyExistsError.
// Default error code is "already_exists".
func ThrowAlreadyExists(parent error, message string, options ...Option) error {
	e := AlreadyExistsError{CreateXError(parent, "already_exists", message)}
	for _, opt := range options {
		opt(e.XError)
	}
	return &e
}

func (err *AlreadyExistsError) IsAlreadyExists() {}

func (err *AlreadyExistsError) Is(target error) bool {
	t, ok := target.(*AlreadyExistsError)
	if !ok {
		return false
	}
	return err.XError.Is(t.XError)
}

func IsErrorAlreadyExists(err error) bool {
	_, ok := err.(AlreadyExists)
	return ok
}

func (err *AlreadyExistsError) Unwrap() error {
	return err.XError
}
