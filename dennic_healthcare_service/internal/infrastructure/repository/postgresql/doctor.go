package postgresql

import (
	"Healthcare_Evrone/internal/entity"
	"Healthcare_Evrone/internal/pkg/otlp"
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"go.opentelemetry.io/otel/attribute"

	// "Healthcare_Evrone/internal/pkg/otlp"
	"Healthcare_Evrone/internal/pkg/postgres"
)

const (
	doctorTableName             = "doctors"
	serviceNameDoctor           = "doctors"
	serviceNameDoctorRepoPrefix = "doctors"
)

type DocTor struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewDoctorRepo(db *postgres.PostgresDB) *DocTor {
	return &DocTor{
		tableName: doctorTableName,
		db:        db,
	}
}

func (h *DocTor) docTorSelectQueryPrefix() string {
	return `id,
      doctor_order,
      first_name,
      last_name,
      image_url,
      gender,
      birth_date,
      phone_number,
      email,
      password,
      address,
      city,
      country,
      salary,
      biography,
      start_work_date,
      end_work_date,
      work_years,
      department_id,
      room_number,
      created_at,
      updated_at
    `
}

func (h *DocTor) getDocTorSelectQueryPrefix() string {
	return `      d.id,                                                            
                        d.doctor_order,
                        d.first_name,
                        d.last_name,
                        d.image_url,
                        d.gender,
                        d.birth_date,
                        d.phone_number,
                        d.email,
                        d.password,
                        d.address,
                        d.city,
                        d.country,
                        d.salary,
            dwh.start_time,
            dwh.finish_time,
            dwh.day_of_week,
                        d.biography,
                        d.start_work_date,
                        d.end_work_date,
                        d.work_years,
                        d.department_id,
                        d.room_number,
                        d.created_at,
                        d.updated_at
            
`
}

func (h *DocTor) CreateDoctor(ctx context.Context, req *entity.Doctor) (*entity.Doctor, error) {

	ctx, span := otlp.Start(ctx, serviceNameDoctor, serviceNameDoctorRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("CreateDoctor").String(req.Id))
	defer span.End()

	data := map[string]any{
		"id":              req.Id,
		"first_name":      req.FirstName,
		"last_name":       req.LastName,
		"image_url":       req.ImageUrl,
		"gender":          req.Gender,
		"birth_date":      req.BirthDate,
		"phone_number":    req.PhoneNumber,
		"email":           req.Email,
		"password":        req.Password,
		"address":         req.Address,
		"city":            req.City,
		"country":         req.Country,
		"salary":          req.Salary,
		"biography":       req.Bio,
		"start_work_date": req.StartWorkDate,
		"work_years":      req.WorkYears,
		"department_id":   req.DepartmentId,
		"room_number":     req.RoomNumber,
	}
	query, args, err := h.db.Sq.Builder.Insert(h.tableName).
		SetMap(data).Suffix(fmt.Sprintf("RETURNING %s", h.docTorSelectQueryPrefix())).ToSql()

	if err != nil {
		return nil, h.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", h.tableName, "create"))
	}
	var startWorkYear, endWorkYear, updatedAt, birthDate sql.NullTime
	err = h.db.QueryRow(ctx, query, args...).Scan(
		&req.Id,
		&req.Order,
		&req.FirstName,
		&req.LastName,
		&req.ImageUrl,
		&req.Gender,
		&birthDate,
		&req.PhoneNumber,
		&req.Email,
		&req.Password,
		&req.Address,
		&req.City,
		&req.Country,
		&req.Salary,
		&req.Bio,
		&startWorkYear,
		&endWorkYear,
		&req.WorkYears,
		&req.DepartmentId,
		&req.RoomNumber,
		&req.CreatedAt,
		&updatedAt,
	)

	req.BirthDate = birthDate.Time.Format("2006-01-02")
	req.StartWorkDate = startWorkYear.Time.String()
	if endWorkYear.Valid {
		req.EndWorkDate = endWorkYear.Time.String()
	}
	if err != nil {
		return nil, h.db.Error(err)
	}

	return req, nil
}

func (h *DocTor) GetDoctorById(ctx context.Context, get *entity.GetReqStr) (*entity.DoctorAndDoctorHours, error) {

	ctx, span := otlp.Start(ctx, serviceNameDoctor, serviceNameDoctorRepoPrefix+"Get")
	span.SetAttributes(attribute.Key(get.Field).String(get.Value))

	defer span.End()

	var doctor entity.DoctorAndDoctorHours
	var startWorkYear, endWorkYear, updatedAt, birthDate, startTime, finishTime sql.NullTime
	queryBuilder := h.db.Sq.Builder.Select(h.getDocTorSelectQueryPrefix()).
		From(h.tableName + " d ").Join("doctor_working_hours dwh ON dwh.doctor_id = d.id")

	if !get.IsActive {
		queryBuilder = queryBuilder.Where("d.deleted_at IS NULL")
	}

	queryBuilder = queryBuilder.Where(h.db.Sq.Equal("d."+get.Field, get.Value))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}
	err = h.db.QueryRow(ctx, query, args...).Scan(
		&doctor.Id,
		&doctor.Order,
		&doctor.FirstName,
		&doctor.LastName,
		&doctor.ImageUrl,
		&doctor.Gender,
		&birthDate,
		&doctor.PhoneNumber,
		&doctor.Email,
		&doctor.Password,
		&doctor.Address,
		&doctor.City,
		&doctor.Country,
		&doctor.Salary,
		&startTime,
		&finishTime,
		&doctor.DayOfWeek,
		&doctor.Bio,
		&startWorkYear,
		&endWorkYear,
		&doctor.WorkYears,
		&doctor.DepartmentId,
		&doctor.RoomNumber,
		&doctor.CreatedAt,
		&updatedAt,
	)
	if updatedAt.Valid {
		doctor.UpdatedAt = updatedAt.Time
	}
	if startTime.Valid {
		doctor.StartTime = startTime.Time.Format("15:04")
	}
	if finishTime.Valid {
		doctor.FinishTime = finishTime.Time.Format("15:04")
	}
	doctor.StartWorkDate = startWorkYear.Time.String()
	if endWorkYear.Valid {
		doctor.EndWorkDate = endWorkYear.Time.String()
	}

	doctor.BirthDate = birthDate.Time.Format("2006-01-02")

	if err != nil {
		return nil, err
	}

	querySpecBuilder := h.db.Sq.Builder.Select("s.id, s.name").
		From("doctor_service ds").Join("doctors d ON d.id = ds.doctor_id").
		Join("specializations s ON ds.specialization_id = s.id").Where(h.db.Sq.Equal("ds.doctor_id", get.Value))

	query, args, err = querySpecBuilder.ToSql()
	if err != nil {
		return nil, err
	}
	specs, err := h.db.Query(ctx, query, args...)
	if err != nil {
		return nil, h.db.Error(err)
	}
	var doctorSpec []entity.DoctorSpec
	for specs.Next() {
		var spec entity.DoctorSpec
		err = specs.Scan(
			&spec.Id, &spec.Name,
		)
		doctorSpec = append(doctorSpec, spec)
		if err != nil {
			return nil, err
		}
		doctor.Specializations = doctorSpec
	}

	return &doctor, nil
}

func (h *DocTor) GetAllDoctors(ctx context.Context, all *entity.GetAll) (*entity.ListDoctorsAndHours, error) {

	ctx, span := otlp.Start(ctx, serviceNameDoctor, serviceNameDoctorRepoPrefix+"Get all")
	span.SetAttributes(attribute.Key(all.Field).String(all.Value))

	defer span.End()

	queryBuilder := h.db.Sq.Builder.Select(h.getDocTorSelectQueryPrefix()).
		From(h.tableName + " d ").Join("doctor_working_hours dwh ON dwh.doctor_id = d.id")
	if all.Field != "" {
		queryBuilder = queryBuilder.Where(fmt.Sprintf(`%s ILIKE '%s'`, all.Field, all.Value+"%"))
	}
	if all.OrderBy != "" {
		queryBuilder = queryBuilder.OrderBy(all.OrderBy)
	}
	countBuilder := h.db.Sq.Builder.Select("count(*)").From(h.tableName)
	if !all.IsActive {
		queryBuilder = queryBuilder.Where("d.deleted_at IS NULL")
		countBuilder = countBuilder.Where("deleted_at IS NULL")
	}

	if all.Limit != 0 && all.Page == 0 {
		offset := all.Limit * (all.Page - 1)
		queryBuilder = queryBuilder.Limit(uint64(all.Limit)).Offset(uint64(offset))
	}

	query, args, err := queryBuilder.ToSql()
	fmt.Println(query)

	if err != nil {
		return nil, err
	}
	var doctors entity.ListDoctorsAndHours
	rows, err := h.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var doctor entity.DoctorAndDoctorHours
		var startWorkYear, endWorkYear, birthDate, updatedAt, startTime, finishTime sql.NullTime
		err = rows.Scan(
			&doctor.Id,
			&doctor.Order,
			&doctor.FirstName,
			&doctor.LastName,
			&doctor.ImageUrl,
			&doctor.Gender,
			&birthDate,
			&doctor.PhoneNumber,
			&doctor.Email,
			&doctor.Password,
			&doctor.Address,
			&doctor.City,
			&doctor.Country,
			&doctor.Salary,
			&startTime,
			&finishTime,
			&doctor.DayOfWeek,
			&doctor.Bio,
			&startWorkYear,
			&endWorkYear,
			&doctor.WorkYears,
			&doctor.DepartmentId,
			&doctor.RoomNumber,
			&doctor.CreatedAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}
		if updatedAt.Valid {
			doctor.UpdatedAt = updatedAt.Time
		}
		if startTime.Valid {
			doctor.StartTime = startTime.Time.Format("15:04")
		}
		if finishTime.Valid {
			doctor.FinishTime = finishTime.Time.Format("15:04")
		}
		doctor.BirthDate = birthDate.Time.Format("2006-01-02")
		doctor.StartWorkDate = startWorkYear.Time.String()
		if endWorkYear.Valid {
			doctor.EndWorkDate = endWorkYear.Time.String()
		}

		querySpecBuilder := h.db.Sq.Builder.Select("s.id, s.name").
			From("doctor_service ds").Join("doctors d ON d.id = ds.doctor_id").
			Join("specializations s ON ds.specialization_id = s.id").Where(h.db.Sq.Equal("ds.doctor_id", doctor.Id))

		query, args, err = querySpecBuilder.ToSql()
		if err != nil {
			return nil, err
		}
		specs, err := h.db.Query(ctx, query, args...)
		if err != nil {
			return nil, h.db.Error(err)
		}
		var doctorSpec []entity.DoctorSpec
		for specs.Next() {
			var spec entity.DoctorSpec
			err = specs.Scan(
				&spec.Id, &spec.Name,
			)
			doctorSpec = append(doctorSpec, spec)
			if err != nil {
				return nil, err
			}
			doctor.Specializations = doctorSpec
		}
		doctors.Doctors = append(doctors.Doctors, doctor)
	}
	var count int64
	queryCount, _, err := countBuilder.ToSql()
	if err != nil {
		return nil, h.db.Error(err)
	}
	err = h.db.QueryRow(ctx, queryCount).Scan(&count)
	if err != nil {
		return nil, h.db.Error(err)
	}
	doctors.Count = count
	return &doctors, nil
}

func (h *DocTor) UpdateDoctor(ctx context.Context, update *entity.Doctor) (*entity.Doctor, error) {

	ctx, span := otlp.Start(ctx, serviceNameDoctor, serviceNameDoctorRepoPrefix+"Update")
	span.SetAttributes(attribute.Key("UpdateDoctor").String(update.Id))

	defer span.End()

	data := map[string]any{
		"first_name":      update.FirstName,
		"last_name":       update.LastName,
		"image_url":       update.ImageUrl,
		"gender":          update.Gender,
		"birth_date":      update.BirthDate,
		"phone_number":    update.PhoneNumber,
		"email":           update.Email,
		"password":        update.Password,
		"address":         update.Address,
		"city":            update.City,
		"country":         update.Country,
		"salary":          update.Salary,
		"biography":       update.Bio,
		"start_work_date": update.StartWorkDate,
		"end_work_date":   update.EndWorkDate,
		"work_years":      update.WorkYears,
		"department_id":   update.DepartmentId,
		"room_number":     update.RoomNumber,
		"updated_at":      time.Now().Add(time.Hour * 5),
	}

	query, args, err := h.db.Sq.Builder.Update(h.tableName).SetMap(data).
		Where(h.db.Sq.Equal("id", update.Id)).Suffix(fmt.Sprintf("RETURNING %s", h.docTorSelectQueryPrefix())).ToSql()

	if err != nil {

		return nil, h.db.ErrSQLBuild(err, h.tableName+" update")
	}
	var startWorkYear, endWorkYear, birthDate sql.NullTime
	err = h.db.QueryRow(ctx, query, args...).Scan(
		&update.Id,
		&update.Order,
		&update.FirstName,
		&update.LastName,
		&update.ImageUrl,
		&update.Gender,
		&birthDate,
		&update.PhoneNumber,
		&update.Email,
		&update.Password,
		&update.Address,
		&update.City,
		&update.Country,
		&update.Salary,
		&update.Bio,
		&startWorkYear,
		&endWorkYear,
		&update.WorkYears,
		&update.DepartmentId,
		&update.RoomNumber,
		&update.CreatedAt,
		&update.UpdatedAt,
	)

	update.StartWorkDate = startWorkYear.Time.String()
	if endWorkYear.Valid {
		update.EndWorkDate = endWorkYear.Time.String()
	}
	update.BirthDate = birthDate.Time.Format("2006-01-02")
	if err != nil {
		return nil, h.db.Error(err)
	}

	return update, nil
}

func (h *DocTor) DeleteDoctor(ctx context.Context, del *entity.GetReqStr) (bool, error) {

	ctx, span := otlp.Start(ctx, serviceNameDoctor, serviceNameDoctorRepoPrefix+"Delete")
	span.SetAttributes(attribute.Key("DeleteDoctor").String(del.Value))

	defer span.End()

	data := map[string]any{
		"deleted_at":    time.Now().Add(time.Hour * 5),
		"end_work_date": time.Now().Add(time.Hour * 5),
	}
	var args []interface{}
	var query string
	var err error
	if del.IsActive {
		query, args, err = h.db.Sq.Builder.Delete(h.tableName).From(h.tableName).
			Where(h.db.Sq.And(h.db.Sq.Equal(del.Field, del.Value))).ToSql()
		if err != nil {
			return false, h.db.ErrSQLBuild(err, h.tableName+" delete")
		}
	} else {
		query, args, err = h.db.Sq.Builder.Update(h.tableName).SetMap(data).
			Where(h.db.Sq.And(h.db.Sq.Equal(del.Field, del.Value), h.db.Sq.Equal("deleted_at", nil))).ToSql()
		if err != nil {
			return false, h.db.ErrSQLBuild(err, h.tableName+" delete")
		}
	}
	resp, err := h.db.Exec(ctx, query, args...)
	if err != nil {
		return false, h.db.Error(err)
	}
	if resp.RowsAffected() > 0 {
		return true, nil
	}

	return false, nil
}

func (h *DocTor) ListDoctorsByDepartmentId(ctx context.Context, in *entity.GetReqStrDep) (*entity.ListDoctors, error) {
	ctx, span := otlp.Start(ctx, serviceNameDoctor, serviceNameDoctorRepoPrefix+"Get all by department_id")
	span.SetAttributes(attribute.Key("ListDoctorsByDepartmentId").String(in.DepartmentId))

	defer span.End()

	offset := in.Limit * (in.Page - 1)

	queryBuilder := h.db.Sq.Builder.Select(h.docTorSelectQueryPrefix()).From(doctorTableName)
	if in.Field != "" {
		queryBuilder = queryBuilder.Where(fmt.Sprintf(`%s ILIKE '%s'`, in.Field, in.Value+"%"))
	}
	if in.OrderBy != "" {
		queryBuilder = queryBuilder.OrderBy(in.OrderBy)
	}
	countBuilder := h.db.Sq.Builder.Select("count(*)").From(h.tableName)
	if !in.IsActive {
		countBuilder = countBuilder.Where("deleted_at IS NULL")
		queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	}

	queryBuilder = queryBuilder.Where(h.db.Sq.Equal("department_id", in.DepartmentId)).Limit(uint64(in.Limit)).Offset(uint64(offset))
	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := h.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	var doctors entity.ListDoctors
	for rows.Next() {
		var doctor entity.Doctor
		var startWorkYear, endWorkYear, birthDate, updatedAt, startTime, finishTime sql.NullTime
		err = rows.Scan(
			&doctor.Id,
			&doctor.Order,
			&doctor.FirstName,
			&doctor.LastName,
			&doctor.ImageUrl,
			&doctor.Gender,
			&birthDate,
			&doctor.PhoneNumber,
			&doctor.Email,
			&doctor.Password,
			&doctor.Address,
			&doctor.City,
			&doctor.Country,
			&doctor.Salary,
			&startTime,
			&finishTime,
			&doctor.Bio,
			&startWorkYear,
			&endWorkYear,
			&doctor.WorkYears,
			&doctor.DepartmentId,
			&doctor.RoomNumber,
			&doctor.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		if updatedAt.Valid {
			doctor.UpdatedAt = updatedAt.Time
		}
		doctor.BirthDate = birthDate.Time.Format("2006-01-02")
		doctor.StartWorkDate = startWorkYear.Time.Format("2006-01-02")
		if endWorkYear.Valid {
			doctor.EndWorkDate = endWorkYear.Time.Format("2006-01-02")
		}

		querySpecBuilder := h.db.Sq.Builder.Select("s.id, s.name").
			From("doctor_service ds").Join("doctors d ON d.id = ds.doctor_id").
			Join("specializations s ON ds.specialization_id = s.id").Where(h.db.Sq.Equal("ds.doctor_id", doctor.Id))

		query, args, err = querySpecBuilder.ToSql()
		if err != nil {
			return nil, err
		}
		specs, err := h.db.Query(ctx, query, args...)
		if err != nil {
			return nil, h.db.Error(err)
		}
		var doctorSpec []entity.DoctorSpec
		for specs.Next() {
			var spec entity.DoctorSpec
			err = specs.Scan(
				&spec.Id, &spec.Name,
			)
			doctorSpec = append(doctorSpec, spec)
			if err != nil {
				return nil, err
			}
			doctor.Specializations = doctorSpec
			doctors.Doctors = append(doctors.Doctors, doctor)
		}
	}
	var count int64
	queryCount, _, err := countBuilder.ToSql()
	if err != nil {
		return nil, h.db.Error(err)
	}

	err = h.db.QueryRow(ctx, queryCount).Scan(&count)
	if err != nil {
		return nil, h.db.Error(err)
	}
	doctors.Count = count
	return &doctors, nil
}

func (h *DocTor) ListDoctorBySpecializationId(ctx context.Context, in *entity.GetReqStrSpec) (*entity.ListDoctorsAndHours, error) {

	ctx, span := otlp.Start(ctx, serviceNameDoctor, serviceNameDoctorRepoPrefix+"Get all by specialization_id")
	span.SetAttributes(attribute.Key("ListDoctorBySpecializationId").String(in.SpecializationId))

	defer span.End()

	offset := in.Limit * (in.Page - 1)

	queryBuilder := h.db.Sq.Builder.Select(h.getDocTorSelectQueryPrefix()).
		From(h.tableName + " d ").
		Join("doctor_service ds ON ds.doctor_id = d.id").
		Join("doctor_working_hours dwh ON dwh.doctor_id = d.id")
	if in.Field != "" {
		queryBuilder = queryBuilder.Where(fmt.Sprintf(`%s ILIKE '%s'`, in.Field, in.Value+"%"))
	}
	if in.OrderBy != "" {
		queryBuilder = queryBuilder.OrderBy(in.OrderBy)
	}
	countBuilder := h.db.Sq.Builder.Select("count(*)").From(h.tableName + " d ").
		Join("doctor_service ds ON ds.doctor_id = d.id").
		Where(fmt.Sprintf("ds.specialization_id = '%s'", in.SpecializationId))
	if !in.IsActive {
		countBuilder = countBuilder.Where("d.deleted_at IS NULL")
		queryBuilder = queryBuilder.Where("d.deleted_at IS NULL")
	}
	queryBuilder = queryBuilder.Where(h.db.Sq.Equal("ds.specialization_id", in.SpecializationId)).
		Limit(uint64(in.Limit)).Offset(uint64(offset))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := h.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	var doctors entity.ListDoctorsAndHours
	for rows.Next() {
		var doctor entity.DoctorAndDoctorHours
		var startWorkYear, endWorkYear, birthDate, updatedAt, startTime, finishTime sql.NullTime
		err = rows.Scan(
			&doctor.Id,
			&doctor.Order,
			&doctor.FirstName,
			&doctor.LastName,
			&doctor.ImageUrl,
			&doctor.Gender,
			&birthDate,
			&doctor.PhoneNumber,
			&doctor.Email,
			&doctor.Password,
			&doctor.Address,
			&doctor.City,
			&doctor.Country,
			&doctor.Salary,
			&startTime,
			&finishTime,
			&doctor.DayOfWeek,
			&doctor.Bio,
			&startWorkYear,
			&endWorkYear,
			&doctor.WorkYears,
			&doctor.DepartmentId,
			&doctor.RoomNumber,
			&doctor.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		if updatedAt.Valid {
			doctor.UpdatedAt = updatedAt.Time
		}
		doctor.BirthDate = birthDate.Time.Format("2006-01-02")
		doctor.StartWorkDate = startWorkYear.Time.Format("2006-01-02")
		if endWorkYear.Valid {
			doctor.EndWorkDate = endWorkYear.Time.Format("2006-01-02")
		}
		if startTime.Valid {
			doctor.StartTime = startTime.Time.Format("2006-01-02")
		}
		if finishTime.Valid {
			doctor.FinishTime = finishTime.Time.Format("2006-01-02")
		}

		querySpecBuilder := h.db.Sq.Builder.Select("s.id, s.name").
			From("doctor_service ds").Join("doctors d ON d.id = ds.doctor_id").
			Join("specializations s ON ds.specialization_id = s.id").Where(h.db.Sq.Equal("ds.doctor_id", doctor.Id))

		query, args, err = querySpecBuilder.ToSql()
		if err != nil {
			return nil, err
		}
		specs, err := h.db.Query(ctx, query, args...)
		if err != nil {
			return nil, h.db.Error(err)
		}
		var doctorSpec []entity.DoctorSpec
		for specs.Next() {
			var spec entity.DoctorSpec
			err = specs.Scan(
				&spec.Id, &spec.Name,
			)
			doctorSpec = append(doctorSpec, spec)
			if err != nil {
				return nil, err
			}
			doctor.Specializations = doctorSpec
			doctors.Doctors = append(doctors.Doctors, doctor)
		}
	}
	var count int64
	queryCount, _, err := countBuilder.ToSql()
	if err != nil {
		return nil, h.db.Error(err)
	}

	err = h.db.QueryRow(ctx, queryCount).Scan(&count)
	if err != nil {
		return nil, h.db.Error(err)
	}

	doctors.Count = count
	return &doctors, nil
}
