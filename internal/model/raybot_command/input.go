package raybotcommand

import (
	"encoding/json"
	"fmt"
)

type Inputs struct {
	union json.RawMessage
}

func NewInputs(inputs json.RawMessage) Inputs {
	return Inputs{union: inputs}
}

func (i Inputs) Raw() json.RawMessage {
	return i.union
}

func (i *Inputs) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &i.union)
}

func (i Inputs) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.union)
}

func (i Inputs) AsMoveToLocationInput() (MoveToLocationInput, error) {
	var ret MoveToLocationInput
	err := json.Unmarshal(i.union, &ret)
	return ret, err
}

func (i *Inputs) FromMoveToLocationInput(m MoveToLocationInput) error {
	ret, err := json.Marshal(m)
	i.union = ret
	return err
}

func (i Inputs) AsLiftBoxInput() (LiftBoxInput, error) {
	var ret LiftBoxInput
	err := json.Unmarshal(i.union, &ret)
	return ret, err
}

func (i *Inputs) FromLiftBoxInput(m LiftBoxInput) error {
	ret, err := json.Marshal(m)
	i.union = ret
	return err
}

func (i Inputs) AsDropBoxInput() (DropBoxInput, error) {
	var ret DropBoxInput
	err := json.Unmarshal(i.union, &ret)
	return ret, err
}

func (i *Inputs) FromDropBoxInput(m DropBoxInput) error {
	ret, err := json.Marshal(m)
	i.union = ret
	return err
}

func (i Inputs) AsCheckQRCodeInput() (CheckQRCodeInput, error) {
	var ret CheckQRCodeInput
	err := json.Unmarshal(i.union, &ret)
	return ret, err
}

func (i *Inputs) FromCheckQRCodeInput(m CheckQRCodeInput) error {
	ret, err := json.Marshal(m)
	i.union = ret
	return err
}

func (i Inputs) AsSpeakInput() (SpeakInput, error) {
	var ret SpeakInput
	err := json.Unmarshal(i.union, &ret)
	return ret, err
}

func (i *Inputs) FromSpeakInput(m SpeakInput) error {
	ret, err := json.Marshal(m)
	i.union = ret
	return err
}

type MoveDirection string

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (m *MoveDirection) UnmarshalText(text []byte) error {
	switch MoveDirection(text) {
	case MoveDirectionForward, MoveDirectionBackward:
		*m = MoveDirection(text)
	default:
		return fmt.Errorf("invalid MoveDirection: %s", text)
	}
	return nil
}

const (
	MoveDirectionForward  MoveDirection = "FORWARD"
	MoveDirectionBackward MoveDirection = "BACKWARD"
)

type MoveToLocationInput struct {
	Location  string        `json:"location"`
	Direction MoveDirection `json:"direction"`
}

type LiftBoxInput struct {
	Distance int32 `json:"distance"`
}

type DropBoxInput struct {
	Distance int32 `json:"distance"`
}

type CheckQRCodeInput struct {
	QRCode string `json:"qr_code"`
}

type SpeakInput struct {
	Text string `json:"text"`
}
