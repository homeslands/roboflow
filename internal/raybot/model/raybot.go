package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Raybot struct {
	ID        uuid.UUID
	Name      string
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Status    RaybotStatus
}

func NewRaybot(name string) (*Raybot, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	now := time.Now()
	return &Raybot{
		ID:        uuid.New(),
		Name:      name,
		Token:     uuid.New().String(),
		Status:    RaybotStatusOffline,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (r *Raybot) Activate() {
	if r.Status == RaybotStatusOffline {
		r.Status = RaybotStatusIdle
		r.UpdatedAt = time.Now()
	}
}

func (r *Raybot) Deactivate() {
	if r.Status != RaybotStatusOffline {
		r.Status = RaybotStatusOffline
		r.UpdatedAt = time.Now()
	}
}

func (r *Raybot) StartWorking() {
	if r.Status == RaybotStatusIdle {
		r.Status = RaybotStatusBusy
		r.UpdatedAt = time.Now()
	}
}

func (r *Raybot) StopWorking() {
	if r.Status == RaybotStatusBusy {
		r.Status = RaybotStatusIdle
		r.UpdatedAt = time.Now()
	}
}
