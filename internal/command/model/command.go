package model

import (
	"time"

	"github.com/google/uuid"
)

type Command struct {
	RaybotID    uuid.UUID
	ID          uuid.UUID
	Type        CommandType
	Status      CommandStatus
	CreatedAt   time.Time
	CompletedAt *time.Time
}

func NewCommand(raybotID uuid.UUID, commandType CommandType) *Command {
	return &Command{
		RaybotID:  raybotID,
		ID:        uuid.New(),
		Type:      commandType,
		Status:    CommandStatusPending,
		CreatedAt: time.Now(),
	}
}
