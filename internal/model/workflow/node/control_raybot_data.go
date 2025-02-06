package node

import (
	"fmt"
)

type ControlRaybotType string

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *ControlRaybotType) UnmarshalText(text []byte) error {
	controlRaybotType := ControlRaybotType(text)
	if _, ok := ControlRaybotTypeMap[controlRaybotType]; !ok {
		return fmt.Errorf("invalid ControlRaybotType: %s", text)
	}
	*t = controlRaybotType
	return nil
}

const (
	ControlRaybotTypeStop           ControlRaybotType = "STOP"
	ControlRaybotTypeMoveForward    ControlRaybotType = "MOVE_FORWARD"
	ControlRaybotTypeMoveBackward   ControlRaybotType = "MOVE_BACKWARD"
	ControlRaybotTypeMoveToLocation ControlRaybotType = "MOVE_TO_LOCATION"
	ControlRaybotTypeOpenBox        ControlRaybotType = "OPEN_BOX"
	ControlRaybotTypeCloseBox       ControlRaybotType = "CLOSE_BOX"
	ControlRaybotTypeLiftBox        ControlRaybotType = "LIFT_BOX"
	ControlRaybotTypeDropBox        ControlRaybotType = "DROP_BOX"
	ControlRaybotTypeCheckQRCode    ControlRaybotType = "CHECK_QR"
	ControlRaybotTypeWaitGetItem    ControlRaybotType = "WAIT_GET_ITEM"
	ControlRaybotTypeScanLocation   ControlRaybotType = "SCAN_LOCATION"
	ControlRaybotTypeSpeak          ControlRaybotType = "SPEAK"
)

var ControlRaybotTypeMap = map[ControlRaybotType]struct{}{
	ControlRaybotTypeStop:           {},
	ControlRaybotTypeMoveForward:    {},
	ControlRaybotTypeMoveBackward:   {},
	ControlRaybotTypeMoveToLocation: {},
	ControlRaybotTypeOpenBox:        {},
	ControlRaybotTypeCloseBox:       {},
	ControlRaybotTypeLiftBox:        {},
	ControlRaybotTypeDropBox:        {},
	ControlRaybotTypeCheckQRCode:    {},
	ControlRaybotTypeWaitGetItem:    {},
	ControlRaybotTypeScanLocation:   {},
	ControlRaybotTypeSpeak:          {},
}

type ControlRaybotData struct {
	ControlRaybotType ControlRaybotType  `json:"control_raybot_type"`
	TimeoutSec        int                `json:"timeout_sec"`
	Input             ControlRaybotInput `json:"input"`
}
