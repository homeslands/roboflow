package dynamicvalue

import (
	"encoding/json"
	"fmt"
)

type SourceType string

// UnmarshalText implements the encoding.TextUnmarshaler interface
func (v *SourceType) UnmarshalText(text []byte) error {
	sourceType := SourceType(text)
	if _, ok := SourceTypeMap[sourceType]; !ok {
		return fmt.Errorf("invalid SourceType: %s", text)
	}
	*v = sourceType
	return nil
}

const (
	SourceTypeStatic    SourceType = "STATIC"
	SourceTypeReference SourceType = "REFERENCE"
)

var SourceTypeMap = map[SourceType]struct{}{
	SourceTypeStatic:    {},
	SourceTypeReference: {},
}

// NodeReference represents a reference to another node's output
type NodeReference struct {
	// NodeID is the unique identifier of the referenced node
	NodeID string `json:"node_id"`
	// Key is the specific output key we want to access from the referenced node
	Key string `json:"key"`
}

// DynamicValue[T] represents a value that can be either a static value
// of type T or a reference to another node's output
type DynamicValue[T any] struct {
	// Type indicates whether this is a STATIC or REFERENCE value
	Type SourceType `json:"type"`
	// StaticValue holds the actual value if this is a static value
	StaticValue *T `json:"static_value,omitempty"`
	// Reference holds the node reference if this is a dynamic value
	Reference *NodeReference `json:"reference,omitempty"`
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (d *DynamicValue[T]) UnmarshalJSON(data []byte) error {
	type Alias DynamicValue[T]
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(d),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("unmarshal dynamic value: %w", err)
	}

	// Validate the type
	switch d.Type {
	case SourceTypeStatic:
		if d.Reference != nil {
			return fmt.Errorf("static type cannot have reference")
		}
		if d.StaticValue == nil {
			return fmt.Errorf("static type must have static_value")
		}
	case SourceTypeReference:
		if d.StaticValue != nil {
			return fmt.Errorf("reference type cannot have static value")
		}
		if d.Reference == nil {
			return fmt.Errorf("reference type must have reference")
		}
		if d.Reference.NodeID == "" {
			return fmt.Errorf("reference must have node_id")
		}
		if d.Reference.Key == "" {
			return fmt.Errorf("reference must have key")
		}
	default:
		return fmt.Errorf("invalid dynamic value type: %s", d.Type)
	}

	return nil
}

// NewStaticValue creates a new DynamicValue with a static value
func NewStaticValue[T any](value T) *DynamicValue[T] {
	return &DynamicValue[T]{
		Type:        SourceTypeStatic,
		StaticValue: &value,
	}
}

// NewReferenceValue creates a new DynamicValue with a reference
func NewReferenceValue[T any](nodeID, key string) *DynamicValue[T] {
	return &DynamicValue[T]{
		Type: SourceTypeReference,
		Reference: &NodeReference{
			NodeID: nodeID,
			Key:    key,
		},
	}
}

// GetStaticValue returns the static value if available
func (d *DynamicValue[T]) GetStaticValue() (T, error) {
	var zero T
	if d.Type != SourceTypeStatic {
		return zero, fmt.Errorf("not a static value")
	}
	if d.StaticValue == nil {
		return zero, fmt.Errorf("static value is nil")
	}
	return *d.StaticValue, nil
}

// GetNodeReference returns the node reference if available
func (d *DynamicValue[T]) GetNodeReference() (*NodeReference, error) {
	if d.Type != SourceTypeReference {
		return nil, fmt.Errorf("not a reference value")
	}
	return d.Reference, nil
}
