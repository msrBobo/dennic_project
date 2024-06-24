package postgresql

import (
	"Healthcare_Evrone/internal/entity"
	"Healthcare_Evrone/internal/pkg/otlp"
	"Healthcare_Evrone/internal/pkg/postgres"
	"context"
	"database/sql"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"time"
)

const (
	reasonsTableName         = "reasons"
	serviceReasons           = "reasons"
	serviceReasonsRepoPrefix = "reasons"
)

type Reasons struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewReasonsRepo(db *postgres.PostgresDB) *Reasons {
	return &Reasons{
		tableName: reasonsTableName,
		db:        db,
	}
}

func (p *Reasons) reasonsSelectQueryPrefix() string {
	return `id,
			name,
			specialization_id,
			image_url,
			created_at,
			updated_at,
			deleted_at
		`
}

func (r Reasons) CreateReasons(ctx context.Context, in *entity.Reasons) (*entity.Reasons, error) {
	ctx, span := otlp.Start(ctx, serviceReasons, serviceReasonsRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("CreateReasons").String(in.Id))

	defer span.End()
	data := map[string]any{
		"id":                in.Id,
		"name":              in.Name,
		"specialization_id": in.SpecializationId,
		"image_url":         in.ImageUrl,
	}
	query, args, err := r.db.Sq.Builder.Insert(r.tableName).SetMap(data).
		Suffix(fmt.Sprintf("RETURNING %s", r.reasonsSelectQueryPrefix())).ToSql()
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, "create"))
	}
	var updatedAt, deletedAt sql.NullTime
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&in.Id,
		&in.Name,
		&in.SpecializationId,
		&in.ImageUrl,
		&in.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	fmt.Println(query, args, "<+++++++++++++++++++++++++")
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, "create"))
	}
	return in, nil
}

func (r Reasons) GetReasonsById(ctx context.Context, in *entity.GetReqStr) (*entity.Reasons, error) {
	ctx, span := otlp.Start(ctx, serviceReasons, serviceReasonsRepoPrefix+"Get")
	span.SetAttributes(attribute.Key(in.Field).String(in.Value))
	defer span.End()

	queryBuilder := r.db.Sq.Builder.Select(r.reasonsSelectQueryPrefix()).From(r.tableName)
	if !in.IsActive {
		queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	}
	queryBuilder = queryBuilder.Where(r.db.Sq.Equal(in.Field, in.Value))
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}
	var reasons entity.Reasons
	var updatedAt, deletedAt sql.NullTime
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&reasons.Id,
		&reasons.Name,
		&reasons.SpecializationId,
		&reasons.ImageUrl,
		&reasons.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, "create"))
	}
	if updatedAt.Valid {
		reasons.UpdatedAt = updatedAt.Time
	}
	if deletedAt.Valid {
		reasons.DeletedAt = deletedAt.Time
	}
	return &reasons, nil
}

func (r *Reasons) GetAllReasons(ctx context.Context, reas *entity.GetAllReas) (*entity.ListReasons, error) {
	ctx, span := otlp.Start(ctx, serviceReasons, serviceReasonsRepoPrefix+"Get all")
	span.SetAttributes(attribute.Key(reas.Field).String(reas.Value))

	defer span.End()

	offset := reas.Limit * (reas.Page - 1)

	queryBuilder := r.db.Sq.Builder.Select(r.reasonsSelectQueryPrefix()).From(r.tableName)
	if reas.Field != "" {
		queryBuilder = queryBuilder.Where(fmt.Sprintf(`%s ILIKE '%s'`, reas.Field, reas.Value+"%"))
	}
	countBuilder := r.db.Sq.Builder.Select("count(*)").From(r.tableName)
	if !reas.IsActive {
		queryBuilder = queryBuilder.Where("deleted_at IS NULL")
		countBuilder = countBuilder.Where("deleted_at IS NULL")
	}
	if reas.OrderBy != "" {
		queryBuilder = queryBuilder.OrderBy(reas.OrderBy)
	}
	queryBuilder = queryBuilder.Limit(uint64(reas.Limit)).Offset(uint64(offset))
	query, args, err := queryBuilder.ToSql()

	if err != nil {

		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	var reasList entity.ListReasons
	for rows.Next() {
		var reasons entity.Reasons
		var updatedAt, deletedAt sql.NullTime
		err = rows.Scan(
			&reasons.Id,
			&reasons.Name,
			&reasons.SpecializationId,
			&reasons.ImageUrl,
			&reasons.CreatedAt,
			&updatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}
		if updatedAt.Valid {
			reasons.UpdatedAt = updatedAt.Time
		}

		if deletedAt.Valid {
			reasons.DeletedAt = deletedAt.Time
		}
		reasList.Reasons = append(reasList.Reasons, reasons)
	}
	var count int32
	queryCount, _, err := countBuilder.ToSql()
	if err != nil {
		return nil, r.db.Error(err)
	}

	err = r.db.QueryRow(ctx, queryCount).Scan(&count)
	if err != nil {
		return nil, r.db.Error(err)
	}
	reasList.Count = count
	return &reasList, nil
}

func (p *Reasons) UpdateReasons(ctx context.Context, reasons *entity.Reasons) (*entity.Reasons, error) {
	ctx, span := otlp.Start(ctx, serviceReasons, serviceReasonsRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("CreateReasons").String(reasons.Id))

	defer span.End()
	data := map[string]any{
		"name":              reasons.Name,
		"specialization_id": reasons.SpecializationId,
		"image_url":         reasons.ImageUrl,
		"updated_at":        time.Now().Add(time.Hour * 5),
	}

	query, args, err := p.db.Sq.Builder.Update(p.tableName).SetMap(data).
		Suffix(fmt.Sprintf("RETURNING %s", p.reasonsSelectQueryPrefix())).ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "update"))
	}
	var updatedAt, deletedAt sql.NullTime
	err = p.db.QueryRow(ctx, query, args...).Scan(
		&reasons.Id,
		&reasons.Name,
		&reasons.SpecializationId,
		&reasons.ImageUrl,
		&reasons.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		return nil, p.db.Error(err)
	}

	if updatedAt.Valid {
		reasons.UpdatedAt = updatedAt.Time
	}
	if deletedAt.Valid {
		reasons.DeletedAt = deletedAt.Time
	}
	return reasons, nil
}

func (p *Reasons) DeleteReasons(ctx context.Context, reasons *entity.GetReqStr) (*entity.StatusReasons, error) {
	ctx, span := otlp.Start(ctx, serviceNameDoctorServices, serviceNameDoctorServicesRepoPrefix+"Delete")
	defer span.End()

	data := map[string]any{
		"deleted_at": time.Now().Add(time.Hour * 5),
	}

	var args []interface{}
	var query string
	var err error
	if reasons.IsActive {
		query, args, err = p.db.Sq.Builder.Delete(p.tableName).From(p.tableName).
			Where(p.db.Sq.And(p.db.Sq.Equal(reasons.Field, reasons.Value))).ToSql()
		if err != nil {
			return &entity.StatusReasons{Status: false}, p.db.ErrSQLBuild(err, p.tableName+" delete")
		}
	} else {
		query, args, err = p.db.Sq.Builder.Update(p.tableName).SetMap(data).
			Where(p.db.Sq.And(p.db.Sq.Equal(reasons.Field, reasons.Value), p.db.Sq.Equal("deleted_at", nil))).ToSql()
		if err != nil {
			return &entity.StatusReasons{Status: false}, p.db.ErrSQLBuild(err, p.tableName+" delete")
		}
	}
	resp, err := p.db.Exec(ctx, query, args...)
	if err != nil {
		return &entity.StatusReasons{Status: false}, p.db.Error(err)
	}
	if resp.RowsAffected() > 0 {
		return &entity.StatusReasons{Status: true}, nil
	}
	return &entity.StatusReasons{Status: false}, nil
}
