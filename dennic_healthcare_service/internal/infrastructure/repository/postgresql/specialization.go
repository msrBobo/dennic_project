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
	specTableName                       = "specializations"
	serviceNameSpecialization           = "doctor_working_hours"
	serviceNameSpecializationRepoPrefix = "doctor_working_hours"
)

type Specialization struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewSpecializationRepo(db *postgres.PostgresDB) *Specialization {
	return &Specialization{
		tableName: specTableName,
		db:        db,
	}
}

func (p *Specialization) specializationSelectQueryPrefix() string {
	return ` id,
			Specializations_order,
			name,
			description,
			department_id,
			image_url,
			created_at,
			updated_at,
			deleted_at
		`
}

func (p *Specialization) CreateSpecialization(ctx context.Context, specialization *entity.Specialization) (*entity.Specialization, error) {

	ctx, span := otlp.Start(ctx, serviceNameSpecialization, serviceNameSpecializationRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("CreateSpecialization").String(specialization.ID))
	defer span.End()

	data := map[string]any{
		"id":            specialization.ID,
		"name":          specialization.Name,
		"description":   specialization.Description,
		"department_id": specialization.DepartmentId,
		"image_url":     specialization.ImageUrl,
	}
	query, args, err := p.db.Sq.Builder.Insert(p.tableName).
		SetMap(data).Suffix(fmt.Sprintf("RETURNING %s", p.specializationSelectQueryPrefix())).ToSql()

	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "create"))
	}
	var updatedAt, deletedAt sql.NullTime
	err = p.db.QueryRow(ctx, query, args...).Scan(
		&specialization.ID,
		&specialization.Order,
		&specialization.Name,
		&specialization.Description,
		&specialization.DepartmentId,
		&specialization.ImageUrl,
		&specialization.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		return nil, p.db.Error(err)
	}
	if updatedAt.Valid {
		specialization.UpdatedAt = updatedAt.Time
	}
	if deletedAt.Valid {
		specialization.DeletedAt = deletedAt.Time
	}

	return specialization, nil
}

func (p *Specialization) GetSpecializationById(ctx context.Context, in *entity.GetReqStr) (*entity.Specialization, error) {

	ctx, span := otlp.Start(ctx, serviceNameSpecialization, serviceNameSpecializationRepoPrefix+"Get")
	span.SetAttributes(attribute.Key(in.Field).String(in.Value))
	defer span.End()

	var spec entity.Specialization
	queryBuilder := p.db.Sq.Builder.Select(p.specializationSelectQueryPrefix()).From(p.tableName)
	if !in.IsActive {
		queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	}
	queryBuilder = queryBuilder.Where(p.db.Sq.Equal(in.Field, in.Value))
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}
	var updatedAt, deletedAt sql.NullTime
	err = p.db.QueryRow(ctx, query, args...).Scan(
		&spec.ID,
		&spec.Order,
		&spec.Name,
		&spec.Description,
		&spec.DepartmentId,
		&spec.ImageUrl,
		&spec.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	if updatedAt.Valid {
		spec.UpdatedAt = updatedAt.Time
	}
	if deletedAt.Valid {
		spec.DeletedAt = deletedAt.Time
	}
	if err != nil {
		return nil, err
	}
	return &spec, nil
}

func (p *Specialization) GetAllSpecializations(ctx context.Context, all *entity.GetAllSpecializations) (*entity.ListSpecializations, error) {

	ctx, span := otlp.Start(ctx, serviceNameSpecialization, serviceNameSpecializationRepoPrefix+"Get all")
	span.SetAttributes(attribute.Key(all.Field).String(all.Value))

	defer span.End()

	offset := all.Limit * (all.Page - 1)

	queryBuilder := p.db.Sq.Builder.Select(p.specializationSelectQueryPrefix()).From(p.tableName)

	countBuilder := p.db.Sq.Builder.Select("count(*)").From(p.tableName)
	if all.DepartmentId != "" {
		queryBuilder = queryBuilder.Where(p.db.Sq.Equal("department_id", all.DepartmentId))
		countBuilder = countBuilder.Where(fmt.Sprintf("department_id = '%s'", all.DepartmentId))
	}

	if all.Field != "" {
		queryBuilder = queryBuilder.Where(fmt.Sprintf(`%s ILIKE '%s'`, all.Field, all.Value+"%"))
	}
	if !all.IsActive {
		queryBuilder = queryBuilder.Where("deleted_at IS NULL")
		countBuilder = countBuilder.Where("deleted_at IS NULL")
	}
	if all.OrderBy != "" {
		queryBuilder = queryBuilder.OrderBy(all.OrderBy)
	}
	queryBuilder = queryBuilder.Limit(uint64(all.Limit)).Offset(uint64(offset))
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}
	var Specializations entity.ListSpecializations
	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var spec entity.Specialization
		var deletedAt, updatedAt sql.NullTime
		err = rows.Scan(
			&spec.ID,
			&spec.Order,
			&spec.Name,
			&spec.Description,
			&spec.DepartmentId,
			&spec.ImageUrl,
			&spec.CreatedAt,
			&updatedAt,
			&deletedAt,
		)
		if updatedAt.Valid {
			spec.UpdatedAt = updatedAt.Time
		}
		if deletedAt.Valid {
			spec.DeletedAt = deletedAt.Time
		}
		if err != nil {
			return nil, err
		}
		Specializations.Specializations = append(Specializations.Specializations, spec)
	}
	var count int32
	queryCount, _, err := countBuilder.ToSql()
	if err != nil {
		return nil, p.db.Error(err)
	}

	fmt.Println(queryCount)
	err = p.db.QueryRow(ctx, queryCount).Scan(&count)
	if err != nil {
		return nil, p.db.Error(err)
	}
	Specializations.Count = count
	return &Specializations, nil
}

func (p *Specialization) UpdateSpecialization(ctx context.Context, in *entity.Specialization) (*entity.Specialization, error) {

	ctx, span := otlp.Start(ctx, serviceNameSpecialization, serviceNameSpecializationRepoPrefix+"Update")
	span.SetAttributes(attribute.Key("UpdateSpecialization").String(in.ID))

	defer span.End()

	data := map[string]any{
		"name":          in.Name,
		"description":   in.Description,
		"department_id": in.DepartmentId,
		"image_url":     in.ImageUrl,
		"updated_at":    time.Now().Add(time.Hour * 5),
	}
	query, args, err := p.db.Sq.Builder.Update(p.tableName).
		SetMap(data).Where(p.db.Sq.Equal("id", in.ID)).Suffix(fmt.Sprintf("RETURNING %s", p.specializationSelectQueryPrefix())).ToSql()

	if err != nil {
		return nil, p.db.ErrSQLBuild(err, p.tableName+" update")
	}
	var deletedAt sql.NullTime
	err = p.db.QueryRow(ctx, query, args...).Scan(
		&in.ID,
		&in.Order,
		&in.Name,
		&in.Description,
		&in.DepartmentId,
		&in.ImageUrl,
		&in.CreatedAt,
		&in.UpdatedAt,
		&deletedAt,
	)
	if deletedAt.Valid {
		in.DeletedAt = deletedAt.Time
	}
	if err != nil {
		return nil, p.db.Error(err)
	}
	return in, nil
}

func (p *Specialization) DeleteSpecialization(ctx context.Context, in *entity.GetReqStr) (bool, error) {

	ctx, span := otlp.Start(ctx, serviceNameSpecialization, serviceNameSpecializationRepoPrefix+"Delete")
	span.SetAttributes(attribute.Key("DeleteSpecialization").String(in.Value))

	defer span.End()

	data := map[string]any{
		"deleted_at": time.Now().Add(time.Hour * 5),
	}

	var args []interface{}
	var query string
	var err error
	if in.IsActive {
		query, args, err = p.db.Sq.Builder.Delete(p.tableName).From(p.tableName).
			Where(p.db.Sq.And(p.db.Sq.Equal(in.Field, in.Value))).ToSql()
		if err != nil {
			return false, p.db.ErrSQLBuild(err, p.tableName+" delete")
		}
	} else {
		query, args, err = p.db.Sq.Builder.Update(p.tableName).SetMap(data).
			Where(p.db.Sq.And(p.db.Sq.Equal(in.Field, in.Value), p.db.Sq.Equal("deleted_at", nil))).ToSql()
		if err != nil {
			return false, p.db.ErrSQLBuild(err, p.tableName+" delete")
		}
	}
	resp, err := p.db.Exec(ctx, query, args...)
	if err != nil {
		return false, p.db.Error(err)
	}
	if resp.RowsAffected() > 0 {
		return true, nil
	}
	return false, nil
}
