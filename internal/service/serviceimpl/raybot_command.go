package serviceimpl

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/internal/db/sqldb"
	raybotcommand "github.com/tuanvumaihuynh/roboflow/internal/model/raybot_command"
	"github.com/tuanvumaihuynh/roboflow/internal/pubsub"
	"github.com/tuanvumaihuynh/roboflow/internal/repository"
	"github.com/tuanvumaihuynh/roboflow/internal/service"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/validator"
)

var _ service.RaybotCommandService = (*raybotCommandService)(nil)

type raybotCommandService struct {
	raybotCommandRepo repository.RaybotCommandRepository
	sqlDBProvider     sqldb.Provider
	publisher         message.Publisher
	validator         validator.Validator
}

func newRaybotCommandService(
	raybotCommandRepo repository.RaybotCommandRepository,
	sqlDBProvider sqldb.Provider,
	publisher message.Publisher,
	validator validator.Validator,
) *raybotCommandService {
	return &raybotCommandService{
		raybotCommandRepo: raybotCommandRepo,
		sqlDBProvider:     sqlDBProvider,
		publisher:         publisher,
		validator:         validator,
	}
}

func (s raybotCommandService) GetRaybotCommand(ctx context.Context, params service.GetRaybotCommandParams) (raybotcommand.RaybotCommand, error) {
	if err := s.validator.Validate(params); err != nil {
		return raybotcommand.RaybotCommand{}, fmt.Errorf("validate params: %w", err)
	}

	rbc, err := s.raybotCommandRepo.GetRaybotCommand(ctx, s.sqlDBProvider.DB(), params.ID)
	if err != nil {
		return raybotcommand.RaybotCommand{}, fmt.Errorf("repo get raybot command: %w", err)
	}

	return rbc, nil
}

func (s raybotCommandService) ListRaybotCommandsByRaybotID(ctx context.Context, params service.ListRaybotCommandsByRaybotIDParams) (paging.List[raybotcommand.RaybotCommand], error) {
	if err := s.validator.Validate(params); err != nil {
		return paging.List[raybotcommand.RaybotCommand]{}, fmt.Errorf("validate params: %w", err)
	}

	rbcs, err := s.raybotCommandRepo.ListRaybotCommandsByRaybotID(ctx, s.sqlDBProvider.DB(), params.RaybotID, params.PagingParams, params.Sorts)
	if err != nil {
		return paging.List[raybotcommand.RaybotCommand]{}, fmt.Errorf("repo list raybot commands by raybot ID: %w", err)
	}

	return rbcs, nil
}

func (s raybotCommandService) CreateRaybotCommand(ctx context.Context, params service.CreateRaybotCommandParams) (raybotcommand.RaybotCommand, error) {
	if err := s.validator.Validate(params); err != nil {
		return raybotcommand.RaybotCommand{}, fmt.Errorf("validate params: %w", err)
	}

	rbc := raybotcommand.NewRaybotCommand(params.RaybotID, params.Type, params.Inputs)
	err := s.raybotCommandRepo.CreateRaybotCommand(ctx, s.sqlDBProvider.DB(), rbc)
	if err != nil {
		return raybotcommand.RaybotCommand{}, fmt.Errorf("repo create raybot command: %w", err)
	}

	// Publish event
	ev := pubsub.RaybotCommandCreated{
		RaybotID:  params.RaybotID,
		CommandID: rbc.ID,
		Type:      rbc.Type,
		Inputs:    rbc.Inputs,
	}
	payload, err := json.Marshal(ev)
	if err != nil {
		return raybotcommand.RaybotCommand{}, fmt.Errorf("marshal event: %w", err)
	}

	msg := message.NewMessage(uuid.NewString(), payload)
	if err := s.publisher.Publish(pubsub.RaybotCommandCreatedTopic, msg); err != nil {
		return raybotcommand.RaybotCommand{}, fmt.Errorf("publisher publish event: %w", err)
	}

	return rbc, nil
}

func (s raybotCommandService) UpdateRaybotCommand(ctx context.Context, params service.UpdateRaybotCommandParams) (raybotcommand.RaybotCommand, error) {
	if err := s.validator.Validate(params); err != nil {
		return raybotcommand.RaybotCommand{}, fmt.Errorf("validate params: %w", err)
	}

	rbc, err := s.raybotCommandRepo.UpdateRaybotCommand(
		ctx,
		s.sqlDBProvider.DB(),
		repository.UpdateRaybotCommandParams{
			ID:             params.ID,
			Status:         params.Status,
			SetStatus:      params.SetStatus,
			Outputs:        params.Outputs,
			SetOutputs:     params.SetOutputs,
			Error:          params.Error,
			SetError:       params.SetError,
			CompletedAt:    params.CompletedAt,
			SetCompletedAt: params.SetCompletedAt,
		},
	)
	if err != nil {
		return raybotcommand.RaybotCommand{}, fmt.Errorf("repo update raybot command: %w", err)
	}

	return rbc, nil
}

func (s raybotCommandService) DeleteRaybotCommand(ctx context.Context, params service.DeleteRaybotCommandParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	if err := s.raybotCommandRepo.DeleteRaybotCommandsByRaybotID(ctx, s.sqlDBProvider.DB(), params.ID); err != nil {
		return fmt.Errorf("repo delete raybot command: %w", err)
	}

	return nil
}
