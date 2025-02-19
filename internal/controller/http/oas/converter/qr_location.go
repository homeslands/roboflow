package converter

import (
	"github.com/tuanvumaihuynh/roboflow/internal/controller/http/oas/gen"
	qrlocation "github.com/tuanvumaihuynh/roboflow/internal/model/qr_location"
)

func ToQRLocationResponse(m qrlocation.QRLocation) gen.QRLocationResponse {
	return gen.QRLocationResponse{
		Id:        m.ID,
		Name:      m.Name,
		QrCode:    m.QRCode,
		Metadata:  m.Metadata,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
