package model

type CommandType string

func (c CommandType) String() string {
	return string(c)
}

const (
	CommandTypeMoveForward    CommandType = "MOVE_FORWARD"
	CommandTypeMoveBackward   CommandType = "MOVE_BACKWARD"
	CommandTypeMoveToLocation CommandType = "MOVE_TO_LOCATION"
	// CommandTypeDropItem     CommandType = "drop_item"
	// CommandTypePickItem     CommandType = "pick_item"
	// CommandTypeBlinkLed     CommandType = "blink_led"
)
