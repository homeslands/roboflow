package node

import (
	"encoding/json"

	dynamicvalue "github.com/tuanvumaihuynh/roboflow/internal/model/workflow/dynamic_value"
)

type ControlRaybotInput struct {
	union json.RawMessage
}

func (i *ControlRaybotInput) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &i.union)
}

func (i ControlRaybotInput) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.union)
}

func (i ControlRaybotInput) AsMoveToLocationInput() (MoveToLocationInput, error) {
	var ret MoveToLocationInput
	err := json.Unmarshal(i.union, &ret)
	return ret, err
}

func (i *ControlRaybotInput) FromMoveToLocationInput(m MoveToLocationInput) error {
	ret, err := json.Marshal(m)
	i.union = ret
	return err
}

func (i ControlRaybotInput) AsLiftBoxInput() (LiftBoxInput, error) {
	var ret LiftBoxInput
	err := json.Unmarshal(i.union, &ret)
	return ret, err
}

func (i *ControlRaybotInput) FromLiftBoxInput(m LiftBoxInput) error {
	ret, err := json.Marshal(m)
	i.union = ret
	return err
}

func (i ControlRaybotInput) AsDropBoxInput() (DropBoxInput, error) {
	var ret DropBoxInput
	err := json.Unmarshal(i.union, &ret)
	return ret, err
}

func (i *ControlRaybotInput) FromDropBoxInput(m DropBoxInput) error {
	ret, err := json.Marshal(m)
	i.union = ret
	return err
}

func (i ControlRaybotInput) AsCheckQRCodeInput() (CheckQRCodeInput, error) {
	var ret CheckQRCodeInput
	err := json.Unmarshal(i.union, &ret)
	return ret, err
}

func (i *ControlRaybotInput) FromCheckQRCodeInput(m CheckQRCodeInput) error {
	ret, err := json.Marshal(m)
	i.union = ret
	return err
}

type MoveToLocationInput struct {
	Location  dynamicvalue.DynamicValue[string] `json:"location"`
	Direction dynamicvalue.DynamicValue[string] `json:"direction"`
}

type LiftBoxInput struct {
	Distance dynamicvalue.DynamicValue[*int32] `json:"distance"`
}

type DropBoxInput struct {
	Distance dynamicvalue.DynamicValue[int32] `json:"distance"`
}

type CheckQRCodeInput struct {
	QRCode dynamicvalue.DynamicValue[string] `json:"qr_code"`
}
