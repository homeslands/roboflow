package repository

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/internal/command/model"
	"github.com/tuanvumaihuynh/roboflow/internal/command/service"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
)

type CommandRepository struct {
	data map[uuid.UUID]*model.Command
	mu   sync.RWMutex
}

var _ service.CommandRepository = (*CommandRepository)(nil)

func NewMemoryCommandRepository() *CommandRepository {
	return &CommandRepository{
		data: make(map[uuid.UUID]*model.Command),
	}
}

func (r *CommandRepository) GetCommand(_ context.Context, id uuid.UUID) (*model.Command, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cmd, ok := r.data[id]
	if !ok {
		return nil, xerrors.ThrowNotFound(nil, "command not found")
	}

	return cmd, nil
}

func (r *CommandRepository) ListCommands(_ context.Context, raybotId uuid.UUID) ([]*model.Command, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var cmds []*model.Command
	for _, cmd := range r.data {
		if cmd.RaybotID == raybotId {
			cmds = append(cmds, cmd)
		}
	}

	return cmds, nil
}

func (r *CommandRepository) CreateCommand(_ context.Context, cmd *model.Command) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[cmd.ID] = cmd
	return nil
}

func (r *CommandRepository) UpdateCommand(_ context.Context, cmd *model.Command) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[cmd.ID] = cmd
	return nil
}

func (r *CommandRepository) DeleteCommand(_ context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.data, id)
	return nil
}
