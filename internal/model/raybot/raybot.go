package raybot

import (
	"time"

	"github.com/google/uuid"
)

type ControlMode string

const (
	ControlModeManual ControlMode = "MANUAL"
	ControlModeAuto   ControlMode = "AUTO"
)

type Raybot struct {
	ID              string
	Name            string
	ControlMode     ControlMode
	IsOnline        bool
	IPAddress       *string
	LastConnectedAt *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewRaybot(name string) Raybot {
	now := time.Now()
	return Raybot{
		ID:              uuid.NewString(),
		Name:            name,
		ControlMode:     ControlModeManual,
		IsOnline:        false,
		IPAddress:       nil,
		LastConnectedAt: nil,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
