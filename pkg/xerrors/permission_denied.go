package xerrors

var (
	_ PermissionDenied = (*PermissionDeniedError)(nil)
	_ Error            = (*PermissionDeniedError)(nil)
)

type PermissionDenied interface {
	error
	IsPermissionDenied()
}

type PermissionDeniedError struct {
	*XError
}

// ThrowPermissionDenied creates a new PermissionDeniedError
// Default error code is "permission_denied".
func ThrowPermissionDenied(parent error, message string, options ...Option) error {
	e := PermissionDeniedError{CreateXError(parent, "permission_denied", message)}
	for _, opt := range options {
		opt(e.XError)
	}
	return &e
}

func (err *PermissionDeniedError) IsPermissionDenied() {}

func IsPermissionDenied(err error) bool {
	_, ok := err.(PermissionDenied)
	return ok
}

func (err *PermissionDeniedError) Is(target error) bool {
	t, ok := target.(*PermissionDeniedError)
	if !ok {
		return false
	}
	return err.XError.Is(t.XError)
}

func (err *PermissionDeniedError) Unwrap() error {
	return err.XError
}
