package node

import (
	"encoding/json"
	"fmt"
)

type TriggerData struct {
	TriggerType TriggerType `json:"trigger_type" validate:"required,enum"`
	union       json.RawMessage
}

func (d *TriggerData) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &d.union)
}

func (d TriggerData) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.union)
}

func (d TriggerData) AsOnDemandTriggerData() (OnDemandTriggerData, error) {
	var ret OnDemandTriggerData
	err := json.Unmarshal(d.union, &ret)
	return ret, err
}

func (d *TriggerData) FromOnDemandTriggerData(o OnDemandTriggerData) error {
	ret, err := json.Marshal(o)
	d.union = ret
	return err
}

type TriggerType string

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *TriggerType) UnmarshalText(text []byte) error {
	triggerType := TriggerType(text)
	if _, ok := TriggerTypeMap[triggerType]; !ok {
		return fmt.Errorf("invalid TriggerType: %s", text)
	}
	*t = triggerType
	return nil
}

const (
	TriggerTypeOnDemand TriggerType = "ON_DEMAND"
	// TriggerTypeSchedule TriggerType = "SCHEDULE"
)

var TriggerTypeMap = map[TriggerType]struct{}{
	TriggerTypeOnDemand: {},
	// TriggerTypeSchedule: {},
}

type InputType string

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (i *InputType) UnmarshalText(text []byte) error {
	inputType := InputType(text)
	if _, ok := InputTypeMap[inputType]; !ok {
		return fmt.Errorf("invalid InputType: %s", text)
	}
	*i = inputType
	return nil
}

const (
	InputTypeString InputType = "STRING"
	InputTypeNumber InputType = "NUMBER"
)

var InputTypeMap = map[InputType]struct{}{
	InputTypeString: {},
	InputTypeNumber: {},
}

type RuntimeVariable struct {
	Key          string    `json:"key" validate:"required,alphanumspace,min=1,max=100"`
	InputType    InputType `json:"input_type" validate:"required,enum"`
	Required     bool      `json:"required" validate:"required"`
	DefaultValue any       `json:"default_value"`
}

type OnDemandTriggerData struct {
	RuntimeVariables []RuntimeVariable `json:"runtime_variables" validate:"required,dive"`
}

func (OnDemandTriggerData) TriggerType() TriggerType {
	return TriggerTypeOnDemand
}
