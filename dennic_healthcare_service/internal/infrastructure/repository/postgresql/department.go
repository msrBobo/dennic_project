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
	departmentTableName             = "departments"
	serviceNameDepartment           = "department"
	serviceNameDepartmentRepoPrefix = "department"
)

type DepartMent struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewDepartmentRepo(db *postgres.PostgresDB) *DepartMent {
	return &DepartMent{
		tableName: departmentTableName,
		db:        db,
	}
}

func (h *DepartMent) departmentSelectQueryPrefix() string {
	return `id,
			department_order,
			name,
			description,
			image_url,
			floor_number,
			short_description,
			created_at,
			updated_at,
			deleted_at`
}

func (h *DepartMent) CreateDepartment(ctx context.Context, dep *entity.Department) (*entity.Department, error) {

	ctx, span := otlp.Start(ctx, serviceNameDepartment, serviceNameDepartmentRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("CreateDepartment").String(dep.Id))
	defer span.End()

	data := map[string]any{
		"id":                dep.Id,
		"name":              dep.Name,
		"description":       dep.Description,
		"image_url":         dep.ImageUrl,
		"floor_number":      dep.FloorNumber,
		"short_description": dep.ShortDescription,
	}
	query, args, err := h.db.Sq.Builder.
		Insert(h.tableName).
		SetMap(data).
		Suffix(fmt.Sprintf("RETURNING %s", h.departmentSelectQueryPrefix())).
		ToSql()

	if err != nil {
		return nil, h.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", h.tableName, " create"))
	}
	var updatedAt, deletedAt sql.NullTime
	var resp entity.Department
	err = h.db.QueryRow(ctx, query, args...).Scan(
		&resp.Id,
		&resp.Order,
		&resp.Name,
		&resp.Description,
		&resp.ImageUrl,
		&resp.FloorNumber,
		&resp.ShortDescription,
		&resp.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		return nil, h.db.Error(err)
	}
	if updatedAt.Valid {
		resp.UpdatedAt = updatedAt.Time
	}
	if deletedAt.Valid {
		resp.DeletedAt = deletedAt.Time
	}

	return &resp, nil
}

func (p *DepartMent) GetDepartmentById(ctx context.Context, get *entity.GetReqStr) (*entity.Department, error) {

	ctx, span := otlp.Start(ctx, serviceNameDepartment, serviceNameDepartmentRepoPrefix+"Get")
	span.SetAttributes(attribute.Key(get.Field).String(get.Value))
	defer span.End()

	var dep entity.Department
	queryBuilder := p.db.Sq.Builder.Select(p.departmentSelectQueryPrefix()).From(departmentTableName)
	if !get.IsActive {
		queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	}
	queryBuilder = queryBuilder.Where(p.db.Sq.Equal(get.Field, get.Value))
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}
	var updatedAt, deletedAt sql.NullTime
	err = p.db.QueryRow(ctx, query, args...).Scan(
		&dep.Id,
		&dep.Order,
		&dep.Name,
		&dep.Description,
		&dep.ImageUrl,
		&dep.FloorNumber,
		&dep.ShortDescription,
		&dep.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	if updatedAt.Valid {
		dep.UpdatedAt = updatedAt.Time
	}
	if deletedAt.Valid {
		dep.DeletedAt = deletedAt.Time
	}
	if err != nil {
		return nil, err
	}
	return &dep, nil
}

func (p *DepartMent) GetAllDepartments(ctx context.Context, all *entity.GetAll) (*entity.ListDepartments, error) {

	ctx, span := otlp.Start(ctx, serviceNameDepartment, serviceNameDepartmentRepoPrefix+"Get All")
	span.SetAttributes(attribute.Key(all.Field).String(all.Value))
	defer span.End()

	offset := all.Limit * (all.Page - 1)

	queryBuilder := p.db.Sq.Builder.Select(p.departmentSelectQueryPrefix()).From(departmentTableName)
	if all.Field != "" {
		queryBuilder = queryBuilder.Where(fmt.Sprintf(`%s ILIKE '%s'`, all.Field, all.Value+"%"))
	}
	if all.OrderBy != "" {
		queryBuilder = queryBuilder.OrderBy(all.OrderBy)
	}
	countBuilder := p.db.Sq.Builder.Select("count(*)").From(departmentTableName)
	if !all.IsActive {
		queryBuilder = queryBuilder.Where("deleted_at IS NULL")
		countBuilder = countBuilder.Where("deleted_at IS NULL")
	}
	queryBuilder = queryBuilder.Limit(uint64(all.Limit)).Offset(uint64(offset))
	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	var departments entity.ListDepartments
	for rows.Next() {
		var dep entity.Department
		var updatedAt, deletedAt sql.NullTime
		err = rows.Scan(
			&dep.Id,
			&dep.Order,
			&dep.Name,
			&dep.Description,
			&dep.ImageUrl,
			&dep.FloorNumber,
			&dep.ShortDescription,
			&dep.CreatedAt,
			&updatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}
		if updatedAt.Valid {
			dep.UpdatedAt = updatedAt.Time
		}
		if deletedAt.Valid {
			dep.DeletedAt = deletedAt.Time
		}
		departments.Departments = append(departments.Departments, dep)
	}
	var count int64
	queryCount, _, err := countBuilder.ToSql()
	if err != nil {
		return nil, p.db.Error(err)
	}

	err = p.db.QueryRow(ctx, queryCount).Scan(&count)
	if err != nil {
		return nil, p.db.Error(err)
	}
	departments.Count = count
	return &departments, nil
}

func (p *DepartMent) UpdateDepartment(ctx context.Context, up *entity.Department) (*entity.Department, error) {

	ctx, span := otlp.Start(ctx, serviceNameDepartment, serviceNameDepartmentRepoPrefix+"Update")
	span.SetAttributes(attribute.Key("UpdateDepartment").String(up.Id))
	defer span.End()
	data := map[string]any{
		"id":                up.Id,
		"name":              up.Name,
		"description":       up.Description,
		"image_url":         up.ImageUrl,
		"floor_number":      up.FloorNumber,
		"short_description": up.ShortDescription,
		"updated_at":        time.Now().Add(time.Hour * 5),
	}
	query, args, err := p.db.Sq.Builder.Update(p.tableName).SetMap(data).Where(p.db.Sq.Equal("id", up.Id)).Suffix(fmt.Sprintf("RETURNING %s", p.departmentSelectQueryPrefix())).ToSql()

	if err != nil {
		return nil, p.db.ErrSQLBuild(err, p.tableName+" update")
	}
	var deletedAt sql.NullTime
	err = p.db.QueryRow(ctx, query, args...).Scan(
		&up.Id,
		&up.Order,
		&up.Name,
		&up.Description,
		&up.ImageUrl,
		&up.FloorNumber,
		&up.ShortDescription,
		&up.CreatedAt,
		&up.UpdatedAt,
		&deletedAt,
	)
	if deletedAt.Valid {
		up.DeletedAt = deletedAt.Time
	}
	if err != nil {
		return nil, p.db.Error(err)
	}
	return up, nil
}

func (p *DepartMent) DeleteDepartment(ctx context.Context, get *entity.GetReqStr) (bool, error) {

	ctx, span := otlp.Start(ctx, serviceNameDepartment, serviceNameDepartmentRepoPrefix+"Delete")
	span.SetAttributes(attribute.Key("DeleteDepartment").String(get.Value))
	defer span.End()

	data := map[string]any{
		"deleted_at": time.Now().Add(time.Hour * 5),
	}
	var args []interface{}
	var query string
	var err error
	if get.IsActive {
		query, args, err = p.db.Sq.Builder.Delete(p.tableName).From(p.tableName).
			Where(p.db.Sq.And(p.db.Sq.Equal(get.Field, get.Value))).ToSql()
		if err != nil {
			return false, p.db.ErrSQLBuild(err, p.tableName+" delete")
		}
	} else {
		query, args, err = p.db.Sq.Builder.Update(p.tableName).SetMap(data).
			Where(p.db.Sq.And(p.db.Sq.Equal(get.Field, get.Value), p.db.Sq.Equal("deleted_at", nil))).ToSql()
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
