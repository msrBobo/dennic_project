package repo

import (
	"booking_service/internal/entity/archive"
	"booking_service/internal/pkg/otlp"
	"booking_service/internal/pkg/postgres"
	"context"
	"database/sql"
	"fmt"
	"time"
)

const (
	tableNameArchive    = "archive"
	serviceNameArchive  = "bookingService"
	spanNameArchiveRepo = "archiveRepo"
)

type BookingArchive struct {
	db *postgres.PostgresDB
}

func NewBookingArchive(db *postgres.PostgresDB) *BookingArchive {
	return &BookingArchive{
		db: db,
	}
}

func tableColumArchive() string {
	return `id,
			doctor_availability_id,
			start_time,
			patient_problem,
			end_time,
			status,
			payment_type,
			payment_amount,
			created_at,
			updated_at,
			deleted_at`
}

func (r *BookingArchive) CreateArchive(ctx context.Context, req *archive.CreatedArchive) (*archive.Archive, error) {
	var (
		archiveRes archive.Archive
		upTime     sql.NullTime
		delTime    sql.NullTime
	)

	ctx, span := otlp.Start(ctx, serviceNameArchive, spanNameArchiveRepo+"Create")
	defer span.End()

	toSql, args, err := r.db.Sq.Builder.
		Insert(tableNameArchive).
		Columns(`doctor_availability_id,
			start_time,
			patient_problem,
			end_time,
			status,
			payment_type,
			payment_amount`).
		Values(
			req.DoctorAvailabilityId,
			req.StartTime,
			req.PatientProblem,
			req.EndTime,
			req.Status,
			req.PaymentType,
			req.PaymentAmount,
		).
		Suffix(fmt.Sprintf("RETURNING %s", tableColumArchive())).
		ToSql()

	if err != nil {
		return nil, err
	}

	if err = r.db.QueryRow(ctx, toSql, args...).Scan(
		&archiveRes.Id,
		&archiveRes.DoctorAvailabilityId,
		&archiveRes.StartTime,
		&archiveRes.PatientProblem,
		&archiveRes.EndTime,
		&archiveRes.Status,
		&archiveRes.PaymentType,
		&archiveRes.PaymentAmount,
		&archiveRes.CreatedAt,
		&upTime,
		&delTime,
	); err != nil {
		return nil, err
	}

	if upTime.Valid {
		archiveRes.UpdatedAt = upTime.Time
	}

	if delTime.Valid {
		archiveRes.DeletedAt = delTime.Time
	}

	return &archiveRes, nil
}

func (r *BookingArchive) GetArchive(ctx context.Context, req *archive.FieldValueReq) (*archive.Archive, error) {
	ctx, span := otlp.Start(ctx, serviceNameArchive, spanNameArchiveRepo+"Get")
	defer span.End()
	var (
		archiveRes archive.Archive
		upTime     sql.NullTime
		delTime    sql.NullTime
	)

	toSql := r.db.Sq.Builder.
		Select(tableColumArchive()).
		From(tableNameArchive).
		Where(r.db.Sq.Equal(req.Field, req.Value))

	if !req.DeleteStatus {
		toSql = toSql.Where(r.db.Sq.Equal("deleted_at", nil))
	}

	toSqls, args, err := toSql.ToSql()
	if err != nil {
		return nil, err
	}

	if err = r.db.QueryRow(ctx, toSqls, args...).Scan(
		&archiveRes.Id,
		&archiveRes.DoctorAvailabilityId,
		&archiveRes.StartTime,
		&archiveRes.PatientProblem,
		&archiveRes.EndTime,
		&archiveRes.Status,
		&archiveRes.PaymentType,
		&archiveRes.PaymentAmount,
		&archiveRes.CreatedAt,
		&upTime,
		&delTime,
	); err != nil {
		return nil, err
	}

	if upTime.Valid {
		archiveRes.UpdatedAt = upTime.Time
	}

	if delTime.Valid {
		archiveRes.DeletedAt = delTime.Time
	}

	return &archiveRes, nil

}

func (r *BookingArchive) GetAllArchive(ctx context.Context, req *archive.GetAllArchives) (*archive.ArchivesType, error) {
	ctx, span := otlp.Start(ctx, serviceNameArchive, spanNameArchiveRepo+"List")
	defer span.End()
	var (
		archivesRes archive.ArchivesType
		upTime      sql.NullTime
		delTime     sql.NullTime
	)
	toSql := r.db.Sq.Builder.
		Select(tableColumArchive()).
		From(tableNameArchive)

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

	rows, err := r.db.Query(ctx, toSqls, args...)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var archiveRes archive.Archive
		if err = rows.Scan(
			&archiveRes.Id,
			&archiveRes.DoctorAvailabilityId,
			&archiveRes.StartTime,
			&archiveRes.PatientProblem,
			&archiveRes.EndTime,
			&archiveRes.Status,
			&archiveRes.PaymentType,
			&archiveRes.PaymentAmount,
			&archiveRes.CreatedAt,
			&upTime,
			&delTime,
		); err != nil {
			return nil, err
		}

		if upTime.Valid {
			archiveRes.UpdatedAt = upTime.Time
		}

		if delTime.Valid {
			archiveRes.DeletedAt = delTime.Time
		}
		archivesRes.Archives = append(archivesRes.Archives, &archiveRes)
		archivesRes.Count += 1
	}

	return &archivesRes, nil
}

func (r *BookingArchive) UpdateArchive(ctx context.Context, req *archive.UpdateArchive) (*archive.Archive, error) {
	ctx, span := otlp.Start(ctx, serviceNameArchive, spanNameArchiveRepo+"Update")
	defer span.End()

	var (
		archiveRes archive.Archive
		upTime     sql.NullTime
		delTime    sql.NullTime
	)
	toSql, args, err := r.db.Sq.Builder.
		Update(tableNameArchive).
		SetMap(map[string]interface{}{
			"doctor_availability_id": req.DoctorAvailabilityId,
			"start_time":             req.StartTime,
			"end_time":               req.EndTime,
			"patient_problem":        req.PatientProblem,
			"status":                 req.Status,
			"payment_type":           req.PaymentType,
			"payment_amount":         req.PaymentAmount,
			"updated_at":             time.Now(),
		}).
		Where(r.db.Sq.Equal(req.Field, req.Value)).
		Suffix(fmt.Sprintf("RETURNING %s", tableColumArchive())).
		ToSql()

	if err != nil {
		return nil, err
	}

	if err = r.db.QueryRow(ctx, toSql, args...).Scan(
		&archiveRes.Id,
		&archiveRes.DoctorAvailabilityId,
		&archiveRes.StartTime,
		&archiveRes.PatientProblem,
		&archiveRes.EndTime,
		&archiveRes.Status,
		&archiveRes.PaymentType,
		&archiveRes.PaymentAmount,
		&archiveRes.CreatedAt,
		&upTime,
		&delTime,
	); err != nil {
		return nil, err
	}

	if upTime.Valid {
		archiveRes.UpdatedAt = upTime.Time
	}

	if delTime.Valid {
		archiveRes.DeletedAt = delTime.Time
	}

	return &archiveRes, nil
}

func (r *BookingArchive) DeleteArchive(ctx context.Context, req *archive.FieldValueReq) (*archive.StatusRes, error) {
	ctx, span := otlp.Start(ctx, serviceNameArchive, spanNameArchiveRepo+"Delete")
	defer span.End()
	if !req.DeleteStatus {
		toSql, args, err := r.db.Sq.Builder.
			Update(tableNameArchive).
			Set("deleted_at", time.Now()).
			Where(r.db.Sq.EqualMany(map[string]interface{}{
				"deleted_at": nil,
				req.Field:    req.Value,
			})).
			ToSql()
		if err != nil {
			return &archive.StatusRes{Status: false}, err
		}

		_, err = r.db.Exec(ctx, toSql, args...)

		if err != nil {
			return &archive.StatusRes{Status: false}, err
		}
		return &archive.StatusRes{Status: true}, nil

	} else {
		toSql, args, err := r.db.Sq.Builder.
			Delete(tableNameArchive).
			Where(r.db.Sq.Equal(req.Field, req.Value)).
			ToSql()

		if err != nil {
			return &archive.StatusRes{Status: false}, err
		}

		_, err = r.db.Exec(ctx, toSql, args...)

		if err != nil {
			return &archive.StatusRes{Status: false}, err
		}
		return &archive.StatusRes{Status: true}, nil
	}
}
