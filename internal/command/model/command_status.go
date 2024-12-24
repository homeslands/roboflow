package model

type CommandStatus string

func (s CommandStatus) String() string {
	return string(s)
}

const (
	CommandStatusPending    CommandStatus = "PENDING"
	CommandStatusInProgress CommandStatus = "IN_PROGRESS"
	CommandStatusSuccess    CommandStatus = "SUCCESS"
	CommandStatusFailed     CommandStatus = "FAILED"
)
