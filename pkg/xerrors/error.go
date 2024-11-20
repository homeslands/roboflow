package xerrors

// ref: https://github.com/zitadel/zitadel/tree/main/internal/zerrors

type Error interface {
	GetParent() error
	GetMessage() string
	SetMessage(string)
	GetCode() string
}
