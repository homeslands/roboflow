package node

import (
	"fmt"
)

// Type is the type of the node.
type Type string

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *Type) UnmarshalText(text []byte) error {
	nodeType := Type(text)
	if _, ok := TypeMap[nodeType]; !ok {
		return fmt.Errorf("invalid NodeType: %s", text)
	}
	*t = nodeType
	return nil
}

const (
	TypeEmpty         Type = "EMPTY"
	TypeTrigger       Type = "TRIGGER"
	TypeControlRaybot Type = "CONTROL_RAYBOT"
)

var TypeMap = map[Type]struct{}{
	TypeEmpty:         {},
	TypeTrigger:       {},
	TypeControlRaybot: {},
}

type Position struct {
	X float32 `json:"x" validate:"required"`
	Y float32 `json:"y" validate:"required"`
}

type Node struct {
	ID          string   `json:"id" validate:"required,uuid"`
	Type        Type     `json:"type" validate:"required,enum"`
	Initialized bool     `json:"initialized" validate:"required"`
	Position    Position `json:"position" validate:"required"`
	Label       string   `json:"label" validate:"required,alphanumspace,min=1,max=100"`
	Data        Data     `json:"data" validate:"required"`
}
