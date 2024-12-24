package event

import "github.com/tuanvumaihuynh/roboflow/internal/model"

const (
	TopicRaybotCommandCreated = "raybot.command.created"
)

type RaybotCommandCreated = model.RaybotCommand
