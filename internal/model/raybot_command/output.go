package raybotcommand

import "encoding/json"

type Outputs struct {
	union json.RawMessage
}

func NewOutputs(outputs json.RawMessage) Outputs {
	return Outputs{union: outputs}
}

func (o Outputs) Raw() json.RawMessage {
	return o.union
}

func (o *Outputs) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &o.union)
}

func (o Outputs) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.union)
}

func (o Outputs) AsScanLocationOutputs() (ScanLocationOutputs, error) {
	var ret ScanLocationOutputs
	err := json.Unmarshal(o.union, &ret)
	return ret, err
}

func (o *Outputs) FromScanLocationOutputs(s ScanLocationOutputs) error {
	ret, err := json.Marshal(s)
	o.union = ret
	return err
}

func (o *Outputs) AsEmptyOutputs() (EmptyOutputs, error) {
	var ret EmptyOutputs
	err := json.Unmarshal(o.union, &ret)
	return ret, err
}

func (o *Outputs) FromEmptyOutputs() error {
	ret, err := json.Marshal(EmptyOutputs{})
	o.union = ret
	return err
}

type ScanLocationOutputs struct {
	Locations []string `json:"locations"`
}

type EmptyOutputs struct {
}
