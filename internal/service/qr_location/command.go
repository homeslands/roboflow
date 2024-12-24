package qrlocation

import (
	"github.com/google/uuid"
	"github.com/tuanvumaihuynh/roboflow/pkg/validator"
)

type CreateQrLocationCommand struct {
	Name     string                 `validate:"required,alphanumspace,min=1,max=100"`
	QRCode   string                 `validate:"required,qrcode,min=1,max=100"`
	Metadata map[string]interface{} `validate:"required"`
}

func (c CreateQrLocationCommand) Validate() error {
	return validator.Validate(c)
}

type DeleteQrLocationCommand struct {
	ID uuid.UUID `validate:"required,uuid"`
}

func (c DeleteQrLocationCommand) Validate() error {
	return validator.Validate(c)
}

type UpdateQRLocationCommand struct {
	ID       uuid.UUID              `validate:"required,uuid"`
	Name     string                 `validate:"required,alphanumspace,min=1,max=100"`
	QRCode   string                 `validate:"required,qrcode,min=1,max=100"`
	Metadata map[string]interface{} `validate:"required"`
}

func (c UpdateQRLocationCommand) Validate() error {
	return validator.Validate(c)
}
