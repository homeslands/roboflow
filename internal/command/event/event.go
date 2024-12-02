package event

import (
	commandModel "github.com/tuanvumaihuynh/roboflow/internal/command/model"
)

const (
	TopicCommandCreated = "command.created"
)

type CommandCreated = commandModel.Command
