package converter

import (
	"github.com/tuanvumaihuynh/roboflow/internal/controller/http/oas/gen"
	raybotcommand "github.com/tuanvumaihuynh/roboflow/internal/model/raybot_command"
)

func ToRaybotCommandResponse(m raybotcommand.RaybotCommand) gen.RaybotCommandResponse {
	return gen.RaybotCommandResponse{
		Id:          m.ID,
		RaybotId:    m.RaybotID,
		Type:        string(m.Type),
		Status:      string(m.Status),
		Inputs:      m.Inputs.Raw(),
		Outputs:     m.Outputs.Raw(),
		Error:       m.Error,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		CompletedAt: m.CompletedAt,
	}
}
