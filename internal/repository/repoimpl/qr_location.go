package repoimpl

import (
	"context"
	"encoding/json"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/tuanvumaihuynh/roboflow/internal/db/sqlcpg"
	"github.com/tuanvumaihuynh/roboflow/internal/db/sqldb"
	qrlocation "github.com/tuanvumaihuynh/roboflow/internal/model/qr_location"
	"github.com/tuanvumaihuynh/roboflow/internal/repository"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/sort"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerror"
)

var (
	_ repository.QRLocationRepository = (*qrLocationRepository)(nil)

	qrLocationQRCodeConstraint = "qr_code"

	ErrQRLocationNotFound  = xerror.NotFound(nil, "qrLocation.notFound", "qr_location not found")
	ErrQRCodeAlreadyExists = xerror.Conflict(nil, "qrLocation.qrCodeAlreadyExists", "qr_code already exists")
)

type qrLocationRepository struct {
	queries sqlcpg.Queries
}

func newQRLocationRepository(queries sqlcpg.Queries) *qrLocationRepository {
	return &qrLocationRepository{queries: queries}
}

func (r qrLocationRepository) GetQRLocation(ctx context.Context, db sqldb.SQLDB, id string) (qrlocation.QRLocation, error) {
	row, err := r.queries.QRLocationGetByID(ctx, db, id)
	if err != nil {
		if sqldb.IsNoRowsError(err) {
			return qrlocation.QRLocation{}, ErrQRLocationNotFound
		}
		return qrlocation.QRLocation{}, fmt.Errorf("queries get qr location by id: %w", err)
	}

	return qrLocationRowToModel(row)
}

func (r qrLocationRepository) ListQRLocations(ctx context.Context, db sqldb.SQLDB, pagingParams paging.Params, sorts []sort.Sort) (paging.List[qrlocation.QRLocation], error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("*").
		From("qr_locations").
		Limit(uint64(pagingParams.Limit())).
		Offset(uint64(pagingParams.Offset()))

	for _, s := range sorts {
		query = s.Attach(query)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return paging.List[qrlocation.QRLocation]{}, fmt.Errorf("build query: %w", err)
	}

	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		return paging.List[qrlocation.QRLocation]{}, fmt.Errorf("queries list qr locations: %w", err)
	}
	defer rows.Close()

	items := make([]qrlocation.QRLocation, 0, pagingParams.Limit())
	for rows.Next() {
		var i sqlcpg.QrLocation
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.QrCode,
			&i.Metadata,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return paging.List[qrlocation.QRLocation]{}, fmt.Errorf("scan qr location: %w", err)
		}

		qrLocation, err := qrLocationRowToModel(i)
		if err != nil {
			return paging.List[qrlocation.QRLocation]{}, fmt.Errorf("convert qr location row to model: %w", err)
		}

		items = append(items, qrLocation)
	}
	if err := rows.Err(); err != nil {
		return paging.List[qrlocation.QRLocation]{}, fmt.Errorf("rows error: %w", err)
	}

	countQuery := psql.Select("COUNT(*)").From("qr_locations")

	countSQL, countArgs, err := countQuery.ToSql()
	if err != nil {
		return paging.List[qrlocation.QRLocation]{}, fmt.Errorf("build count query: %w", err)
	}

	var count int64
	if err := db.QueryRow(ctx, countSQL, countArgs...).Scan(&count); err != nil {
		return paging.List[qrlocation.QRLocation]{}, fmt.Errorf("queries count qr locations: %w", err)
	}

	return paging.NewList(items, count), nil
}

func (r qrLocationRepository) CreateQRLocation(ctx context.Context, db sqldb.SQLDB, qrLocation qrlocation.QRLocation) error {
	metadata, err := json.Marshal(qrLocation.Metadata)
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}

	err = r.queries.QRLocationInsert(ctx, db, sqlcpg.QRLocationInsertParams{
		ID:       qrLocation.ID,
		Name:     qrLocation.Name,
		QrCode:   qrLocation.QRCode,
		Metadata: metadata,
	})
	if err != nil {
		if sqldb.IsNoRowsError(err) {
			return ErrQRLocationNotFound
		}
		if sqldb.IsUniqueViolationError(err, qrLocationQRCodeConstraint) {
			return ErrQRCodeAlreadyExists
		}
		return fmt.Errorf("queries create qr location: %w", err)
	}

	return nil
}

func (r qrLocationRepository) UpdateQRLocation(ctx context.Context, db sqldb.SQLDB, params repository.UpdateQRLocationParams) (qrlocation.QRLocation, error) {
	metadata, err := json.Marshal(params.Metadata)
	if err != nil {
		return qrlocation.QRLocation{}, fmt.Errorf("marshal metadata: %w", err)
	}

	row, err := r.queries.QRLocationUpdate(ctx, db, sqlcpg.QRLocationUpdateParams{
		ID:          params.ID,
		Name:        params.Name,
		SetName:     params.SetName,
		QrCode:      params.QRCode,
		SetQrCode:   params.SetQRCode,
		Metadata:    metadata,
		SetMetadata: params.SetMetadata,
	})
	if err != nil {
		if sqldb.IsNoRowsError(err) {
			return qrlocation.QRLocation{}, ErrQRLocationNotFound
		}
		if sqldb.IsUniqueViolationError(err, qrLocationQRCodeConstraint) {
			return qrlocation.QRLocation{}, ErrQRCodeAlreadyExists
		}
		return qrlocation.QRLocation{}, fmt.Errorf("queries update qr location: %w", err)
	}

	return qrLocationRowToModel(row)
}

func (r qrLocationRepository) DeleteQRLocation(ctx context.Context, db sqldb.SQLDB, id string) error {
	if err := r.queries.QRLocationDelete(ctx, db, id); err != nil {
		if sqldb.IsNoRowsError(err) {
			return ErrQRLocationNotFound
		}
		return fmt.Errorf("queries delete qr location: %w", err)
	}

	return nil
}

func qrLocationRowToModel(row sqlcpg.QrLocation) (qrlocation.QRLocation, error) {
	var metadata map[string]interface{}
	if err := json.Unmarshal(row.Metadata, &metadata); err != nil {
		return qrlocation.QRLocation{}, fmt.Errorf("unmarshal metadata: %w", err)
	}

	return qrlocation.QRLocation{
		ID:        row.ID,
		Name:      row.Name,
		QRCode:    row.QrCode,
		Metadata:  metadata,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}
