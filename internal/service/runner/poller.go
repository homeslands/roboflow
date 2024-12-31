package runner

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
)

type RaybotCommandStatusPoller struct {
	raybotCommandRepo model.RaybotCommandRepository
}

func (e RaybotCommandStatusPoller) pollRaybotCommandStatus(ctx context.Context, raybotCommandID uuid.UUID) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			status, err := e.raybotCommandRepo.GetStatus(ctx, raybotCommandID)
			if err != nil {
				return fmt.Errorf("failed to get raybot command by id: %w", err)
			}

			switch status {
			case model.RaybotCommandStatusSuccess:
				return nil
			case model.RaybotCommandStatusFailed:
				return errors.New("raybot command failed")
			}
		}

	}
}
