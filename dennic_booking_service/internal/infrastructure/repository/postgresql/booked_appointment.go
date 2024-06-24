package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	appointment "booking_service/internal/entity/booked_appointments"
	"booking_service/internal/pkg/otlp"
	"booking_service/internal/pkg/postgres"
)

const (
	tableNameAppointment    = "booked_appointments"
	serviceNameAppointment  = "appointment"
	spanNameAppointmentRepo = "appointment"
)

type BookingAppointment struct {
	db *postgres.PostgresDB
}

func NewBookingAppointment(db *postgres.PostgresDB) *BookingAppointment {
	return &BookingAppointment{
		db: db,
	}
}

func tableColums() string {
	return `id, 
			department_id, 
			doctor_id, 
			patient_id, 
			doctor_service_id,
			appointment_date, 
			appointment_time, 
			duration, 
			key, 
			expires_at, 
			patient_problem,
			status,
			payment_type,
            payment_amount,
			user_id,
			created_at, 
			updated_at, 
			deleted_at`
}

func (r *BookingAppointment) CreateAppointment(
	ctx context.Context,
	req *appointment.CreateAppointment,
) (*appointment.Appointment, error) {
	fmt.Println(req.PatientProblem, "=============")

	ctx, span := otlp.Start(ctx, serviceNameAppointment, spanNameAppointmentRepo+"Create")
	defer span.End()
	var (
		response appointment.Appointment
		upAt     sql.NullTime
		delAt    sql.NullTime
	)
	toSql, args, err := r.db.Sq.Builder.
		Insert(tableNameAppointment).
		Columns(` 
				department_id, 
				doctor_id, 
				patient_id, 
				doctor_service_id,
				appointment_date, 
				appointment_time, 
				duration, 
				key, 
				expires_at, 
				patient_problem,
				status,
				payment_type,
				payment_amount,
				user_id`).
		Values(
			req.DepartmentId,
			req.DoctorId,
			req.PatientId,
			req.ServiceId,
			req.AppointmentDate.String(),
			req.AppointmentTime,
			req.Duration,
			req.Key,
			req.ExpiresAt,
			req.PatientProblem,
			req.Status,
			req.PaymentType,
			req.PaymentAmount,
			req.UserId).
		Suffix(fmt.Sprintf("RETURNING %s", tableColums())).
		ToSql()
	if err != nil {
		return nil, err
	}

	var (
		userId sql.NullString
	)
	if err = r.db.QueryRow(ctx, toSql, args...).Scan(
		&response.Id,
		&response.DepartmentId,
		&response.DoctorId,
		&response.PatientId,
		&response.ServiceId,
		&response.AppointmentDate,
		&response.AppointmentTime,
		&response.Duration,
		&response.Key,
		&response.ExpiresAt,
		&response.PatientProblem,
		&response.Status,
		&response.PaymentType,
		&response.PaymentAmount,
		&userId,
		&response.CreatedAt,
		&upAt,
		&delAt,
	); err != nil {
		return nil, err
	}

	if userId.Valid {
		response.UserId = userId.String
	}

	if upAt.Valid {
		response.UpdatedAt = upAt.Time
	}

	if delAt.Valid {
		response.DeletedAt = delAt.Time
	}

	return &response, nil
}

func (r *BookingAppointment) GetAppointment(
	ctx context.Context,
	req *appointment.FieldValueReq,
) (*appointment.Appointment, error) {
	ctx, span := otlp.Start(ctx, serviceNameAppointment, spanNameAppointmentRepo+"Get")
	defer span.End()

	var (
		response appointment.Appointment
		userId   sql.NullString
		upAt     sql.NullTime
		delAt    sql.NullTime
	)

	toSql := r.db.Sq.Builder.
		Select(tableColums()).
		From(tableNameAppointment)

	if !req.DeleteStatus {
		toSql = toSql.Where(r.db.Sq.Equal("deleted_at", nil))
	}
	toSqls, args, err := toSql.Where(r.db.Sq.Equal(req.Field, req.Value)).ToSql()
	if err != nil {
		return nil, err
	}
	if err = r.db.QueryRow(ctx, toSqls, args...).Scan(
		&response.Id,
		&response.DepartmentId,
		&response.DoctorId,
		&response.PatientId,
		&response.ServiceId,
		&response.AppointmentDate,
		&response.AppointmentTime,
		&response.Duration,
		&response.Key,
		&response.ExpiresAt,
		&response.PatientProblem,
		&response.Status,
		&response.PaymentType,
		&response.PaymentAmount,
		&userId,
		&response.CreatedAt,
		&upAt,
		&delAt,
	); err != nil {
		return nil, err
	}

	if userId.Valid {
		response.UserId = userId.String
	}

	if upAt.Valid {
		response.UpdatedAt = upAt.Time
	}

	if delAt.Valid {
		response.DeletedAt = delAt.Time
	}

	return &response, nil
}

func (r *BookingAppointment) GetAllAppointment(
	ctx context.Context,
	req *appointment.GetAllAppointment,
) (*appointment.AppointmentsType, error) {
	ctx, span := otlp.Start(ctx, serviceNameAppointment, spanNameAppointmentRepo+"List")
	defer span.End()

	var (
		response appointment.AppointmentsType
		upAt     sql.NullTime
		delAt    sql.NullTime
		count    int64
	)

	toSql := r.db.Sq.Builder.
		Select(tableColums()).
		From(tableNameAppointment)

	countBuilder := r.db.Sq.Builder.Select("count(*)").From(tableNameAppointment)

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
		countBuilder = countBuilder.Where(r.db.Sq.Equal("deleted_at", nil))
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
		var (
			res    appointment.Appointment
			userId sql.NullString
		)
		if err := rows.Scan(
			&res.Id,
			&res.DepartmentId,
			&res.DoctorId,
			&res.PatientId,
			&res.ServiceId,
			&res.AppointmentDate,
			&res.AppointmentTime,
			&res.Duration,
			&res.Key,
			&res.ExpiresAt,
			&res.PatientProblem,
			&res.Status,
			&res.PaymentType,
			&res.PaymentAmount,
			&userId,
			&res.CreatedAt,
			&upAt,
			&delAt,
		); err != nil {
			return nil, err
		}

		if userId.Valid {
			res.UserId = userId.String
		}

		if upAt.Valid {
			res.UpdatedAt = upAt.Time
		}

		if delAt.Valid {
			res.DeletedAt = delAt.Time
		}

		response.Appointments = append(response.Appointments, &res)
	}

	response.Count = count
	return &response, nil
}

func (r *BookingAppointment) GetFilteredAppointments(
	ctx context.Context,
	req *appointment.GetFilteredRequest,
) (*appointment.AppointmentsType, error) {
	ctx, span := otlp.Start(ctx, serviceNameAppointment, spanNameAppointmentRepo+"List")
	defer span.End()

	var (
		response appointment.AppointmentsType
		upAt     sql.NullTime
		delAt    sql.NullTime
		count    int64
	)

	toSql := r.db.Sq.Builder.
		Select(tableColums()).
		From(tableNameAppointment)

	countBuilder := r.db.Sq.Builder.Select("count(*)").From(tableNameAppointment)

	if req.Page >= 1 && req.Limit >= 1 {
		toSql = toSql.
			Limit(req.Limit).
			Offset(req.Limit * (req.Page - 1))
	}

	if req.Field != "" && req.Value != "" {
		toSql = toSql.Where(r.db.Sq.Equal(req.Field, req.Value))
		countBuilder = countBuilder.Where(r.db.Sq.Equal(req.Field, req.Value))
	}

	if req.Status != "" {
		toSql = toSql.Where(r.db.Sq.Equal("status", req.Status))
		countBuilder = countBuilder.Where(r.db.Sq.Equal("status", req.Status))

	}

	if req.OrderBy != "" {
		toSql = toSql.OrderBy(req.OrderBy)
	}
	if !req.DeleteStatus {
		countBuilder = countBuilder.Where(r.db.Sq.Equal("deleted_at", nil))
		toSql = toSql.Where(r.db.Sq.Equal("deleted_at", nil))
	}

	toSqls, args, err := toSql.ToSql()
	if err != nil {
		return nil, err
	}

	queryCount, argss, err := countBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, queryCount, argss...).Scan(&count)

	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, toSqls, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			res    appointment.Appointment
			userId sql.NullString
		)

		if err := rows.Scan(
			&res.Id,
			&res.DepartmentId,
			&res.DoctorId,
			&res.PatientId,
			&res.ServiceId,
			&res.AppointmentDate,
			&res.AppointmentTime,
			&res.Duration,
			&res.Key,
			&res.ExpiresAt,
			&res.PatientProblem,
			&res.Status,
			&res.PaymentType,
			&res.PaymentAmount,
			&userId,
			&res.CreatedAt,
			&upAt,
			&delAt,
		); err != nil {
			return nil, err
		}

		if userId.Valid {
			res.UserId = userId.String
		}

		if upAt.Valid {
			res.UpdatedAt = upAt.Time
		}

		if delAt.Valid {
			res.DeletedAt = delAt.Time
		}

		response.Appointments = append(response.Appointments, &res)
	}

	response.Count = count
	return &response, nil
}

func (r *BookingAppointment) UpdateAppointment(
	ctx context.Context,
	req *appointment.UpdateAppointment,
) (*appointment.Appointment, error) {
	ctx, span := otlp.Start(ctx, serviceNameAppointment, spanNameAppointmentRepo+"Update")
	defer span.End()

	var (
		response appointment.Appointment
		userId   sql.NullString
		upAt     sql.NullTime
		delAt    sql.NullTime
	)
	toSql, args, err := r.db.Sq.Builder.
		Update(tableNameAppointment).
		SetMap(map[string]interface{}{
			"department_id":     req.DepartmentId,
			"patient_id":        req.PatientId,
			"doctor_id":         req.DoctorId,
			"doctor_service_id": req.ServiceId,
			"appointment_date":  req.AppointmentDate.String(),
			"appointment_time":  req.AppointmentTime,
			"duration":          req.Duration,
			"key":               req.Key,
			"expires_at":        req.ExpiresAt,
			"patient_problem":   req.PatientProblem,
			"payment_type":      req.PaymentType,
			"payment_amount":    req.PaymentAmount,
			"status":            req.Status,
			"user_id":           req.UserId,
			"updated_at":        time.Now(),
		}).
		Where(r.db.Sq.Equal(req.Field, req.Value)).
		Suffix(fmt.Sprintf("RETURNING %s", tableColums())).
		ToSql()
	if err != nil {
		return nil, err
	}
	if err = r.db.QueryRow(ctx, toSql, args...).Scan(
		&response.Id,
		&response.DepartmentId,
		&response.DoctorId,
		&response.PatientId,
		&response.ServiceId,
		&response.AppointmentDate,
		&response.AppointmentTime,
		&response.Duration,
		&response.Key,
		&response.ExpiresAt,
		&response.PatientProblem,
		&response.Status,
		&response.PaymentType,
		&response.PaymentAmount,
		&userId,
		&response.CreatedAt,
		&upAt,
		&delAt,
	); err != nil {
		return nil, err
	}

	if userId.Valid {
		response.UserId = userId.String
	}

	if upAt.Valid {
		response.UpdatedAt = upAt.Time
	}

	if delAt.Valid {
		response.DeletedAt = delAt.Time
	}

	return &response, nil
}

func (r *BookingAppointment) DeleteAppointment(
	ctx context.Context,
	req *appointment.FieldValueReq,
) (*appointment.StatusRes, error) {
	ctx, span := otlp.Start(ctx, serviceNameAppointment, spanNameAppointmentRepo+"Delete")
	defer span.End()
	if !req.DeleteStatus {
		toSql, args, err := r.db.Sq.Builder.
			Update(tableNameAppointment).
			Set("deleted_at", time.Now()).
			Where(r.db.Sq.EqualMany(map[string]interface{}{
				"deleted_at": nil,
				req.Field:    req.Value,
			})).
			ToSql()
		if err != nil {
			return &appointment.StatusRes{Status: false}, err
		}

		resp, err := r.db.Exec(ctx, toSql, args...)
		if err != nil {
			return &appointment.StatusRes{Status: false}, err
		}
		if resp.RowsAffected() > 0 {
			return &appointment.StatusRes{Status: true}, nil
		}
		return &appointment.StatusRes{Status: false}, nil

	} else {
		toSql, args, err := r.db.Sq.Builder.
			Delete(tableNameAppointment).
			Where(r.db.Sq.Equal(req.Field, req.Value)).
			ToSql()
		if err != nil {
			return &appointment.StatusRes{Status: false}, err
		}

		resp, err := r.db.Exec(ctx, toSql, args...)
		if err != nil {
			return &appointment.StatusRes{Status: false}, err
		}

		if resp.RowsAffected() > 0 {
			return &appointment.StatusRes{Status: true}, nil
		}
		return &appointment.StatusRes{Status: false}, nil
	}
}
