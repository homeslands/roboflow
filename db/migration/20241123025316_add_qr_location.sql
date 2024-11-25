-- +goose Up
-- +goose StatementBegin
CREATE TABLE qr_locations (
    "id" UUID PRIMARY KEY,
    "name" TEXT NOT NULL,
    "qr_code" TEXT NOT NULL UNIQUE,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "qr_locations";
-- +goose StatementEnd
