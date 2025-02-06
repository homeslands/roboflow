package node

import (
	"encoding/json"
)

type Data struct {
	union json.RawMessage
}

func (d *Data) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &d.union)
}

func (d Data) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.union)
}

func (d Data) AsEmptyData() (EmptyData, error) {
	var ret EmptyData
	err := json.Unmarshal(d.union, &ret)
	return ret, err
}

func (d *Data) FromEmptyData(e EmptyData) error {
	ret, err := json.Marshal(e)
	d.union = ret
	return err
}

func (d Data) AsTriggerData() (TriggerData, error) {
	var ret TriggerData
	err := json.Unmarshal(d.union, &ret)
	return ret, err
}

func (d *Data) FromTriggerData(t TriggerData) error {
	ret, err := json.Marshal(t)
	d.union = ret
	return err
}

func (d Data) AsControlRaybotInput() (ControlRaybotInput, error) {
	var ret ControlRaybotInput
	err := json.Unmarshal(d.union, &ret)
	return ret, err
}

func (d *Data) FromControlRaybotInput(c ControlRaybotInput) error {
	ret, err := json.Marshal(c)
	d.union = ret
	return err
}
