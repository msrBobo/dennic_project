package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"booking_service/internal/entity/patients"
	"booking_service/internal/pkg/otlp"
	"booking_service/internal/pkg/postgres"
)

const (
	tableNamePatients   = "patients"
	serviceNamePatient  = "patientsRepo"
	spanNamePatientRepo = "patientsRepo"
)

type BookingPatients struct {
	db *postgres.PostgresDB
}

func NewBookingPatients(db *postgres.PostgresDB) *BookingPatients {
	return &BookingPatients{
		db: db,
	}
}

func tableColumPatients() string {
	return `id,
			first_name,
			last_name,
			birth_date,
			gender,
			blood_group,
			phone_number,
			city,
			country,
			address,
			patient_problem,
			created_at,
			updated_at,
			deleted_at`
}

func (r *BookingPatients) CreatePatient(
	ctx context.Context,
	req *patients.CreatedPatient,
) (*patients.Patient, error) {
	ctx, span := otlp.Start(ctx, serviceNamePatient, spanNamePatientRepo+"Create")
	defer span.End()
	var (
		patient patients.Patient
		upTime  sql.NullTime
		delTime sql.NullTime
	)
	toSql, args, err := r.db.Sq.Builder.
		Insert(tableNamePatients).
		Columns(`id,
							first_name,
							last_name,
							birth_date,
							gender,
							blood_group,
							phone_number,
							city,
							country,
							address,
							patient_problem`).
		Values(req.Id,
			req.FirstName,
			req.LastName,
			req.BirthDate.String(),
			req.Gender,
			req.BloodGroup,
			req.PhoneNumber,
			req.City,
			req.Country,
			req.Address,
			req.PatientProblem).
		Suffix(fmt.Sprintf("RETURNING %s", tableColumPatients())).
		ToSql()
	if err != nil {
		return nil, err
	}

	if err = r.db.QueryRow(ctx, toSql, args...).Scan(
		&patient.Id,
		&patient.FirstName,
		&patient.LastName,
		&patient.BirthDate,
		&patient.Gender,
		&patient.BloodGroup,
		&patient.PhoneNumber,
		&patient.City,
		&patient.Country,
		&patient.Address,
		&patient.PatientProblem,
		&patient.CreatedAt,
		&upTime,
		&delTime,
	); err != nil {
		return nil, err
	}

	if upTime.Valid {
		patient.UpdatedAt = upTime.Time
	}

	if delTime.Valid {
		patient.DeletedAt = delTime.Time
	}

	return &patient, nil
}

func (r *BookingPatients) GetPatient(
	ctx context.Context,
	req *patients.FieldValueReq,
) (*patients.Patient, error) {
	ctx, span := otlp.Start(ctx, serviceNamePatient, spanNamePatientRepo+"Get")
	defer span.End()

	var (
		patient patients.Patient
		upTime  sql.NullTime
		delTime sql.NullTime
	)

	toSql := r.db.Sq.Builder.
		Select(tableColumPatients()).
		From(tableNamePatients).
		Where(r.db.Sq.Equal(req.Field, req.Value))

	if !req.DeleteStatus {
		toSql = toSql.Where(r.db.Sq.Equal("deleted_at", nil))
	}

	toSqls, args, err := toSql.ToSql()
	if err != nil {
		return nil, err
	}

	if err = r.db.QueryRow(ctx, toSqls, args...).Scan(
		&patient.Id,
		&patient.FirstName,
		&patient.LastName,
		&patient.BirthDate,
		&patient.Gender,
		&patient.BloodGroup,
		&patient.PhoneNumber,
		&patient.City,
		&patient.Country,
		&patient.Address,
		&patient.PatientProblem,
		&patient.CreatedAt,
		&upTime,
		&delTime,
	); err != nil {
		return nil, err
	}

	if upTime.Valid {
		patient.UpdatedAt = upTime.Time
	}

	if delTime.Valid {
		patient.DeletedAt = delTime.Time
	}

	return &patient, nil
}

func (r *BookingPatients) GetAllPatiens(
	ctx context.Context,
	req *patients.GetAllPatients,
) (*patients.PatientsType, error) {
	ctx, span := otlp.Start(ctx, serviceNamePatient, spanNamePatientRepo+"List")
	defer span.End()

	var (
		patientss patients.PatientsType
		upTime    sql.NullTime
		delTime   sql.NullTime
		count     int64
	)

	toSql := r.db.Sq.Builder.
		Select(tableColumPatients()).
		From(tableNamePatients)

	countBuilder := r.db.Sq.Builder.Select("count(*)").From(tableNamePatients)

	toSql = toSql.
		Limit(req.Limit).
		Offset(req.Limit * (req.Page - 1))

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
		var res patients.Patient
		if err := rows.Scan(
			&res.Id,
			&res.FirstName,
			&res.LastName,
			&res.BirthDate,
			&res.Gender,
			&res.BloodGroup,
			&res.PhoneNumber,
			&res.City,
			&res.Country,
			&res.Address,
			&res.PatientProblem,
			&res.CreatedAt,
			&upTime,
			&delTime,
		); err != nil {
			return nil, err
		}

		patientss.Patients = append(patientss.Patients, &res)
	}

	patientss.Count = count
	return &patientss, nil
}

func (r *BookingPatients) UpdatePatient(
	ctx context.Context,
	req *patients.UpdatePatient,
) (*patients.Patient, error) {
	ctx, span := otlp.Start(ctx, serviceNamePatient, spanNamePatientRepo+"Update")
	defer span.End()

	var (
		patient patients.Patient
		upTime  sql.NullTime
		delTime sql.NullTime
	)
	toSql, args, err := r.db.Sq.Builder.
		Update(tableNamePatients).
		SetMap(map[string]interface{}{
			"first_name":      req.FirstName,
			"last_name":       req.LastName,
			"birth_date":      req.BirthDate.String(),
			"gender":          req.Gender,
			"blood_group":     req.BloodGroup,
			"city":            req.City,
			"country":         req.Country,
			"address":         req.Address,
			"patient_problem": req.PatientProblem,
			"updated_at":      time.Now(),
		}).
		Where(r.db.Sq.Equal(req.Field, req.Value)).
		Suffix(fmt.Sprintf("RETURNING %s", tableColumPatients())).
		ToSql()
	if err != nil {
		return nil, err
	}

	if err = r.db.QueryRow(ctx, toSql, args...).Scan(
		&patient.Id,
		&patient.FirstName,
		&patient.LastName,
		&patient.BirthDate,
		&patient.Gender,
		&patient.BloodGroup,
		&patient.PhoneNumber,
		&patient.City,
		&patient.Country,
		&patient.Address,
		&patient.PatientProblem,
		&patient.CreatedAt,
		&upTime,
		&delTime,
	); err != nil {
		return nil, err
	}

	if upTime.Valid {
		patient.UpdatedAt = upTime.Time
	}

	if delTime.Valid {
		patient.DeletedAt = delTime.Time
	}

	return &patient, nil
}

func (r *BookingPatients) UpdatePhonePatient(
	ctx context.Context,
	req *patients.UpdatePhoneNumber,
) (*patients.StatusRes, error) {
	ctx, span := otlp.Start(ctx, serviceNamePatient, spanNamePatientRepo+"UpdatePhone")
	defer span.End()

	toSql, args, err := r.db.Sq.Builder.
		Update(tableNamePatients).
		SetMap(map[string]interface{}{
			"phone_number": req.PhoneNumber,
			"updated_at":   time.Now(),
		}).
		Where(r.db.Sq.Equal(req.Field, req.Value)).
		ToSql()
	if err != nil {
		return &patients.StatusRes{Status: false}, err
	}

	_, err = r.db.Exec(ctx, toSql, args...)
	if err != nil {
		return &patients.StatusRes{Status: false}, err
	}

	return &patients.StatusRes{Status: true}, nil
}

func (r *BookingPatients) DeletePatient(
	ctx context.Context,
	req *patients.FieldValueReq,
) (*patients.StatusRes, error) {
	ctx, span := otlp.Start(ctx, serviceNamePatient, spanNamePatientRepo+"Delete")
	defer span.End()

	if !req.DeleteStatus {
		toSql, args, err := r.db.Sq.Builder.
			Update(tableNamePatients).
			Set("deleted_at", time.Now()).
			Where(r.db.Sq.EqualMany(map[string]interface{}{
				"deleted_at": nil,
				req.Field:    req.Value,
			})).
			ToSql()
		if err != nil {
			return &patients.StatusRes{Status: false}, err
		}

		resp, err := r.db.Exec(ctx, toSql, args...)
		if err != nil {
			return &patients.StatusRes{Status: false}, err
		}
		if resp.RowsAffected() > 0 {
			return &patients.StatusRes{Status: true}, nil
		}
		return &patients.StatusRes{Status: false}, nil

	} else {
		toSql, args, err := r.db.Sq.Builder.
			Delete(tableNamePatients).
			Where(r.db.Sq.Equal(req.Field, req.Value)).
			ToSql()
		if err != nil {
			return &patients.StatusRes{Status: false}, err
		}

		resp, err := r.db.Exec(ctx, toSql, args...)
		if err != nil {
			return &patients.StatusRes{Status: false}, err
		}

		if resp.RowsAffected() > 0 {
			return &patients.StatusRes{Status: true}, nil
		}
		return &patients.StatusRes{Status: false}, nil
	}
}
