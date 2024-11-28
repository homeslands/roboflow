package model

type RaybotStatus string

func (s RaybotStatus) String() string {
	return string(s)
}

const (
	RaybotStatusOffline RaybotStatus = "OFFLINE"
	RaybotStatusIdle    RaybotStatus = "IDLE"
	RaybotStatusBusy    RaybotStatus = "BUSY"
)
