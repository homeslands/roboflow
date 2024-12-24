package service

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/internal/command/event"
	"github.com/tuanvumaihuynh/roboflow/internal/command/model"
	"github.com/tuanvumaihuynh/roboflow/internal/command/port"
)

var _ port.Service = (*CommandService)(nil)

type CommandRepository interface {
	GetCommand(ctx context.Context, id uuid.UUID) (*model.Command, error)
	ListCommands(ctx context.Context, raybotId uuid.UUID) ([]*model.Command, error)
	CreateCommand(ctx context.Context, cmd *model.Command) error
	UpdateCommand(ctx context.Context, cmd *model.Command) error
	DeleteCommand(ctx context.Context, id uuid.UUID) error
}

type EventPublisher interface {
	Publish(topic string, payload []byte) error
}

type CommandService struct {
	commandRepo    CommandRepository
	eventPublisher EventPublisher
}

func NewCommandService(commandRepo CommandRepository, publisher EventPublisher) *CommandService {
	if commandRepo == nil {
		panic("nil commandRepo")
	}
	if publisher == nil {
		panic("nil publisher")
	}
	return &CommandService{
		commandRepo:    commandRepo,
		eventPublisher: publisher,
	}
}

func (s CommandService) GetCommand(ctx context.Context, id uuid.UUID) (*model.Command, error) {
	return s.commandRepo.GetCommand(ctx, id)
}

func (s CommandService) ListCommands(ctx context.Context, raybotId uuid.UUID) ([]*model.Command, error) {
	return s.commandRepo.ListCommands(ctx, raybotId)
}

func (s CommandService) CreateCommand(ctx context.Context, cmd *model.Command) (*model.Command, error) {
	if err := s.commandRepo.CreateCommand(ctx, cmd); err != nil {
		return nil, err
	}

	ev := event.CommandCreated(*cmd)
	msgJSON, err := json.Marshal(ev)
	if err != nil {
		return nil, err
	}

	if err := s.eventPublisher.Publish(event.TopicCommandCreated, msgJSON); err != nil {
		return nil, err
	}

	return cmd, nil
}

func (s CommandService) UpdateCommand(ctx context.Context, cmd *model.Command) (*model.Command, error) {
	if err := s.commandRepo.UpdateCommand(ctx, cmd); err != nil {
		return nil, err
	}

	// msgJSON, err := json.Marshal(cmd)
	// if err != nil {
	// 	return nil, err
	// }
	// if err := s.eventPublisher.Publish("command.updated", msgJSON); err != nil {
	// 	return nil, err
	// }
	return cmd, nil
}

func (s CommandService) DeleteCommand(ctx context.Context, id uuid.UUID) error {
	if err := s.commandRepo.DeleteCommand(ctx, id); err != nil {
		return err
	}

	// msgJSON, err := json.Marshal(id)
	// if err != nil {
	// 	return err
	// }
	// if err := s.eventPublisher.Publish("command.deleted", msgJSON); err != nil {
	// 	return err
	// }
	return nil
}
