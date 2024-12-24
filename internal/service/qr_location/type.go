package qrlocation

import (
	"regexp"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/validator"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
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

type GetQrLocationByIDQuery struct {
	ID uuid.UUID `validate:"required,uuid"`
}

func (q GetQrLocationByIDQuery) Validate() error {
	return validator.Validate(q)
}

type ListQrLocationQuery struct {
	PagingParams paging.Params
	Sorts        []xsort.Sort
}

var (
	allowedSortFieldsRegexp = regexp.MustCompile("^(name|qr_code|created_at|updated_at)$")
)

func (q ListQrLocationQuery) Validate() error {
	for _, sort := range q.Sorts {
		if !allowedSortFieldsRegexp.MatchString(sort.Col) {
			return xerrors.ThrowInvalidArgument(nil, "invalid sort field")
		}
	}
	return nil
}
