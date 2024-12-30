package raybotclient

type InboundResponseErrorData struct {
	Reason string `json:"reason"`
}

type InboundResponseScanLocationData struct {
	Locations []string `json:"locations"`
}
