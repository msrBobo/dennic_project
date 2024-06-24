package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"time"

	"booking_service/internal/entity/doctor_availability"
	"booking_service/internal/pkg/otlp"
	"booking_service/internal/pkg/postgres"
)

const (
	tableNameDoctorAvailability    = "doctor_availability"
	serviceNameDoctorAvailability  = "doctor_availability"
	spanNameDoctorAvailabilityRepo = "doctor_availability"
)

type DoctorAvailability struct {
	db *postgres.PostgresDB
}

func NewDoctorAvailability(db *postgres.PostgresDB) *DoctorAvailability {
	return &DoctorAvailability{
		db: db,
	}
}

func tableColumDoctorAvailability() string {
	return `id,
			department_id,
			doctor_id,
			doctor_date,
			start_time,
			end_time,
			status,
			created_at,
			updated_at,
			deleted_at`
}

func (r *DoctorAvailability) CreateDoctorAvailability(
	ctx context.Context,
	req *doctor_availability.CreateDoctorAvailability,
) (*doctor_availability.DoctorAvailability, error) {
	ctx, span := otlp.Start(
		ctx,
		serviceNameDoctorAvailability,
		spanNameDoctorAvailabilityRepo+"Create",
	)
	defer span.End()

	var (
		docAvail doctor_availability.DoctorAvailability
		upAt     sql.NullTime
		delAt    sql.NullTime
	)

	toSql, args, err := r.db.Sq.Builder.
		Insert(tableNameDoctorAvailability).
		Columns(`department_id,
			doctor_id,
			doctor_date,
			start_time,
			end_time,
			status`).
		Values(
			req.DepartmentId,
			req.DoctorId,
			req.DoctorDate.String(),
			req.StartTime,
			req.EndTime,
			req.Status,
		).
		Suffix(fmt.Sprintf("RETURNING %s", tableColumDoctorAvailability())).
		ToSql()
	if err != nil {
		return nil, err
	}

	if err = r.db.QueryRow(ctx, toSql, args...).Scan(
		&docAvail.Id,
		&docAvail.DepartmentId,
		&docAvail.DoctorId,
		&docAvail.DoctorDate,
		&docAvail.StartTime,
		&docAvail.EndTime,
		&docAvail.Status,
		&docAvail.CreatedAt,
		&upAt,
		&delAt,
	); err != nil {
		return nil, err
	}

	if upAt.Valid {
		docAvail.UpdatedAt = upAt.Time
	}

	if delAt.Valid {
		docAvail.DeletedAt = delAt.Time
	}

	return &docAvail, nil
}

func (r *DoctorAvailability) GetDoctorAvailability(
	ctx context.Context,
	req *doctor_availability.FieldValueReq,
) (*doctor_availability.DoctorAvailability, error) {
	ctx, span := otlp.Start(
		ctx,
		serviceNameDoctorAvailability,
		spanNameDoctorAvailabilityRepo+"Get",
	)
	defer span.End()

	var (
		docAvail doctor_availability.DoctorAvailability
		upAt     sql.NullTime
		delAt    sql.NullTime
	)

	toSql := r.db.Sq.Builder.
		Select(tableColumDoctorAvailability()).
		From(tableNameDoctorAvailability).
		Where(r.db.Sq.Equal(req.Field, req.Value))

	if !req.DeleteStatus {
		toSql = toSql.Where(r.db.Sq.Equal("deleted_at", nil))
	}

	toSqls, args, err := toSql.ToSql()
	if err != nil {
		return nil, err
	}

	if err = r.db.QueryRow(ctx, toSqls, args...).Scan(
		&docAvail.Id,
		&docAvail.DepartmentId,
		&docAvail.DoctorId,
		&docAvail.DoctorDate,
		&docAvail.StartTime,
		&docAvail.EndTime,
		&docAvail.Status,
		&docAvail.CreatedAt,
		&upAt,
		&delAt,
	); err != nil {
		return nil, err
	}

	if upAt.Valid {
		docAvail.UpdatedAt = upAt.Time
	}

	if delAt.Valid {
		docAvail.DeletedAt = delAt.Time
	}

	return &docAvail, nil
}

func (r *DoctorAvailability) GetAllDoctorAvailability(
	ctx context.Context,
	req *doctor_availability.GetAllReq,
) (*doctor_availability.DoctorAvailabilityType, error) {
	ctx, span := otlp.Start(
		ctx,
		serviceNameDoctorAvailability,
		spanNameDoctorAvailabilityRepo+"List",
	)
	defer span.End()

	var (
		docAvails doctor_availability.DoctorAvailabilityType
		upAt      sql.NullTime
		delAt     sql.NullTime
		count     int64
	)

	toSql := r.db.Sq.Builder.
		Select(tableColumDoctorAvailability()).
		From(tableNameDoctorAvailability)

	countBuilder := r.db.Sq.Builder.Select("count(*)").From(tableNameDoctorAvailability)

	if req.Page >= 1 && req.Limit >= 1 {
		toSql = toSql.
			Limit(req.Limit).
			Offset(req.Limit * (req.Page - 1))
	}
	if req.Value != "" {
		toSql = toSql.Where(r.db.Sq.ILike(req.Field, req.Value+"%"))
	}
	if req.OrderBy != "" {
		toSql = toSql.OrderBy(req.OrderBy)
	}
	if !req.DeleteStatus {
		toSql = toSql.Where(r.db.Sq.Equal("deleted_at", nil))
	}
	toSqls, args, err := toSql.ToSql()
	if err != nil {
		return nil, err
	}

	queryCount, _, err := countBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, queryCount).Scan(&count)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, toSqls, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var docAvail doctor_availability.DoctorAvailability
		if err = rows.Scan(
			&docAvail.Id,
			&docAvail.DepartmentId,
			&docAvail.DoctorId,
			&docAvail.DoctorDate,
			&docAvail.StartTime,
			&docAvail.EndTime,
			&docAvail.Status,
			&docAvail.CreatedAt,
			&upAt,
			&delAt,
		); err != nil {
			return nil, err
		}

		if upAt.Valid {
			docAvail.UpdatedAt = upAt.Time
		}

		if delAt.Valid {
			docAvail.DeletedAt = delAt.Time
		}

		docAvails.DoctorAvailabilitys = append(docAvails.DoctorAvailabilitys, &docAvail)
	}
	docAvails.Count = count
	return &docAvails, nil
}

func (r *DoctorAvailability) UpdateDoctorAvailability(
	ctx context.Context,
	req *doctor_availability.UpdateDoctorAvailability,
) (*doctor_availability.DoctorAvailability, error) {
	ctx, span := otlp.Start(
		ctx,
		serviceNameDoctorAvailability,
		spanNameDoctorAvailabilityRepo+"Update",
	)
	defer span.End()

	var (
		docAvail doctor_availability.DoctorAvailability
		upAt     sql.NullTime
		delAt    sql.NullTime
	)

	toSql, args, err := r.db.Sq.Builder.
		Update(tableNameDoctorAvailability).
		SetMap(map[string]interface{}{
			"department_id": req.DepartmentId,
			"doctor_id":     req.DoctorId,
			"doctor_date":   req.DoctorDate.String(),
			"start_time":    req.StartTime,
			"end_time":      req.EndTime,
			"status":        req.Status,
			"updated_at":    time.Now(),
		}).
		Where(r.db.Sq.Equal(req.Field, req.Value)).
		Suffix(fmt.Sprintf("RETURNING %s", tableColumDoctorAvailability())).
		ToSql()
	if err != nil {
		return &doctor_availability.DoctorAvailability{}, err
	}

	if err = r.db.QueryRow(ctx, toSql, args...).Scan(
		&docAvail.Id,
		&docAvail.DepartmentId,
		&docAvail.DoctorId,
		&docAvail.DoctorDate,
		&docAvail.StartTime,
		&docAvail.EndTime,
		&docAvail.Status,
		&docAvail.CreatedAt,
		&upAt,
		&delAt,
	); err != nil {
		return nil, err
	}

	if upAt.Valid {
		docAvail.UpdatedAt = upAt.Time
	}

	if delAt.Valid {
		docAvail.DeletedAt = delAt.Time
	}

	return &docAvail, nil
}

func (r *DoctorAvailability) DeleteDoctorAvailability(
	ctx context.Context,
	req *doctor_availability.FieldValueReq,
) (*doctor_availability.StatusRes, error) {
	ctx, span := otlp.Start(
		ctx,
		serviceNameDoctorAvailability,
		spanNameDoctorAvailabilityRepo+"Delete",
	)
	defer span.End()

	if !req.DeleteStatus {
		toSql, args, err := r.db.Sq.Builder.
			Update(tableNameDoctorAvailability).
			Set("deleted_at", time.Now()).
			Where(r.db.Sq.EqualMany(map[string]interface{}{
				"deleted_at": nil,
				req.Field:    req.Value,
			})).
			ToSql()
		if err != nil {
			return &doctor_availability.StatusRes{Status: false}, err
		}

		resp, err := r.db.Exec(ctx, toSql, args...)
		if err != nil {
			return &doctor_availability.StatusRes{Status: false}, err
		}
		if resp.RowsAffected() > 0 {
			return &doctor_availability.StatusRes{Status: true}, nil
		}
		return &doctor_availability.StatusRes{Status: false}, nil

	} else {
		toSql, args, err := r.db.Sq.Builder.
			Delete(tableNameDoctorAvailability).
			Where(r.db.Sq.Equal(req.Field, req.Value)).
			ToSql()
		if err != nil {
			return &doctor_availability.StatusRes{Status: false}, err
		}

		resp, err := r.db.Exec(ctx, toSql, args...)
		if err != nil {
			return &doctor_availability.StatusRes{Status: false}, err
		}

		if resp.RowsAffected() > 0 {
			return &doctor_availability.StatusRes{Status: true}, nil
		}
		return &doctor_availability.StatusRes{Status: false}, nil
	}
}

func (r *DoctorAvailability) GeDoctorTimeStatus(
	ctx context.Context,
	req *doctor_availability.DoctorTimeStatusReq,
) (*doctor_availability.TimeStatus, error) {
	ctx, span := otlp.Start(
		ctx,
		serviceNameDoctorAvailability,
		spanNameDoctorAvailabilityRepo+"GetStatus",
	)
	defer span.End()

	var status string

	toSql := r.db.Sq.Builder.
		Select("status").
		From(tableNameDoctorAvailability).
		Where(r.db.Sq.Equal("department_id", req.DepartmentID),
			r.db.Sq.Equal("doctor_id", req.DoctorID),
			r.db.Sq.Equal("doctor_date", req.DoctorDate),
			r.db.Sq.Equal("start_time", req.StartTime))

	toSqls, args, err := toSql.ToSql()
	if err != nil {
		return nil, err
	}

	if err = r.db.QueryRow(ctx, toSqls, args...).Scan(
		&status,
	); err != nil {
		return nil, err
	}

	return &doctor_availability.TimeStatus{Status: status}, nil
}

func (r *DoctorAvailability) UpdateTimeStatus(
	ctx context.Context,
	req *doctor_availability.DoctorTimeUpdateStatusReq,
) (*doctor_availability.TimeStatus, error) {
	ctx, span := otlp.Start(
		ctx,
		serviceNameDoctorAvailability,
		spanNameDoctorAvailabilityRepo+"UpdateStatus",
	)
	defer span.End()

	toSql := r.db.Sq.Builder.
		Update(tableNameDoctorAvailability).
		Set("status", req.Status).
		Where(squirrel.And{
			squirrel.Eq{"status": req.Status},
			squirrel.Eq{"department_id": req.DepartmentID},
			squirrel.Eq{"doctor_id": req.DoctorID},
			squirrel.Eq{"doctor_date": req.DoctorDate},
			squirrel.Eq{"start_time": req.StartTime},
		})

	toSqls, args, err := toSql.ToSql()
	if err != nil {
		return nil, err
	}

	_, err = r.db.Exec(ctx, toSqls, args...)

	if err != nil {
		return nil, err
	}

	return &doctor_availability.TimeStatus{Status: req.Status}, nil
}

func (r *DoctorAvailability) GetDoctorTimeByDoctorID(
	ctx context.Context,
	req *doctor_availability.DoctorTimeFieldValueByIdReq,
) (*doctor_availability.Hours, error) {
	ctx, span := otlp.Start(
		ctx,
		serviceNameDoctorAvailability,
		spanNameDoctorAvailabilityRepo+"Get",
	)
	defer span.End()

	var hours doctor_availability.Hours

	toSql := r.db.Sq.Builder.
		Select("start_time").
		From(tableNameDoctorAvailability).
		Where(squirrel.And{
			squirrel.Eq{"doctor_id": req.DoctorID},
			squirrel.Eq{"doctor_date": req.Date},
		})

	toSqls, args, err := toSql.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, toSqls, args...)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var hour time.Time
		err = rows.Scan(
			&hour)
		if err != nil {
			return nil, err
		}
		hours.Hours = append(hours.Hours, hour.Format("15:04"))
	}
	return &hours, nil
}
