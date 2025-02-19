package converter

import (
	"github.com/tuanvumaihuynh/roboflow/internal/controller/http/oas/gen"
	"github.com/tuanvumaihuynh/roboflow/internal/model/raybot"
)

func ToRaybotResponse(m raybot.Raybot) gen.RaybotResponse {
	return gen.RaybotResponse{
		Id:              m.ID,
		Name:            m.Name,
		ControlMode:     string(m.ControlMode),
		IsOnline:        m.IsOnline,
		IpAddress:       m.IPAddress,
		LastConnectedAt: m.LastConnectedAt,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}
