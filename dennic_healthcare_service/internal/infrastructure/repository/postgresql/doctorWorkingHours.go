package postgresql

import (
	"Healthcare_Evrone/internal/entity"
	"Healthcare_Evrone/internal/pkg/otlp"
	"Healthcare_Evrone/internal/pkg/postgres"
	"context"
	"database/sql"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

const (
	doctorWorkingHoursTableName             = "doctor_working_hours"
	serviceNameDoctorWorkingHours           = "doctor_working_hours"
	serviceNameDoctorWorkingHoursRepoPrefix = "doctor_working_hours"
)

type Dwh struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewDoctorWorkingHoursRepo(db *postgres.PostgresDB) *Dwh {
	return &Dwh{
		tableName: doctorWorkingHoursTableName,
		db:        db,
	}
}

func (p *Dwh) doctorWorkingHoursSelectQueryPrefix() string {
	return `id,
			doctor_id,
			day_of_week,
			start_time,
			finish_time,
			created_at,
			updated_at,
			deleted_at
		`
}

func (p Dwh) CreateDoctorWorkingHours(ctx context.Context, in *entity.DoctorWorkingHours) (*entity.DoctorWorkingHours, error) {
	ctx, span := otlp.Start(ctx, serviceNameDoctorWorkingHours, serviceNameDoctorWorkingHoursRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("CreateDoctorWorkingHours").String(string(in.Id)))
	defer span.End()
	data := map[string]any{
		"doctor_id":   in.DoctorId,
		"day_of_week": in.DayOfWeek,
		"start_time":  in.StartTime,
		"finish_time": in.FinishTime,
	}
	query, args, err := p.db.Sq.Builder.Insert(p.tableName).SetMap(data).
		Suffix(fmt.Sprintf("RETURNING %s", p.doctorWorkingHoursSelectQueryPrefix())).ToSql()

	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "create"))
	}
	var updatedAt, deletedAt, startTime, finishTime sql.NullTime
	err = p.db.QueryRow(ctx, query, args...).Scan(
		&in.Id,
		&in.DoctorId,
		&in.DayOfWeek,
		&startTime,
		&finishTime,
		&in.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		return nil, p.db.Error(err)
	}
	in.StartTime = startTime.Time.Format("15:04:05")
	in.FinishTime = finishTime.Time.Format("15:04:05")

	if updatedAt.Valid {
		in.UpdatedAt = updatedAt.Time
	}
	if deletedAt.Valid {
		in.DeletedAt = deletedAt.Time
	}
	return in, nil
}

func (p Dwh) GetDoctorWorkingHoursById(ctx context.Context, in *entity.GetRequest) (*entity.DoctorWorkingHours, error) {

	ctx, span := otlp.Start(ctx, serviceNameDoctorWorkingHours, serviceNameDoctorWorkingHoursRepoPrefix+"Get")
	span.SetAttributes(attribute.Key(in.Field).String(in.Value))

	defer span.End()

	var doctorWorkingHours entity.DoctorWorkingHours
	queryBuilder := p.db.Sq.Builder.Select(p.doctorWorkingHoursSelectQueryPrefix()).From(p.tableName)
	if !in.IsActive {
		queryBuilder = queryBuilder.Where(p.db.Sq.Equal("deleted_at", nil))
	}
	queryBuilder = queryBuilder.Where(p.db.Sq.Equal(in.Field, in.Value))

	if in.DayOfWeek != "" {
		queryBuilder = queryBuilder.Where(p.db.Sq.Equal("day_of_week", in.DayOfWeek))
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}
	var updatedAt, deletedAt, startTime, finishTime sql.NullTime
	err = p.db.QueryRow(ctx, query, args...).Scan(
		&doctorWorkingHours.Id,
		&doctorWorkingHours.DoctorId,
		&doctorWorkingHours.DayOfWeek,
		&startTime,
		&finishTime,
		&doctorWorkingHours.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	doctorWorkingHours.StartTime = startTime.Time.Format("15:04:05")
	doctorWorkingHours.FinishTime = finishTime.Time.Format("15:04:05")

	if updatedAt.Valid {
		doctorWorkingHours.UpdatedAt = updatedAt.Time
	}
	if deletedAt.Valid {
		doctorWorkingHours.DeletedAt = deletedAt.Time
	}
	if err != nil {
		return nil, err
	}
	return &doctorWorkingHours, nil
}

func (p Dwh) GetAllDoctorWorkingHours(ctx context.Context, all *entity.GetAll) (*entity.ListDoctorWorkingHours, error) {

	ctx, span := otlp.Start(ctx, serviceNameDoctorWorkingHours, serviceNameDoctorWorkingHoursRepoPrefix+"Get all")
	span.SetAttributes(attribute.Key(all.Field).String(all.Value))

	defer span.End()

	offset := all.Limit * (all.Page - 1)

	queryBuilder := p.db.Sq.Builder.Select(p.doctorWorkingHoursSelectQueryPrefix()).From(p.tableName)
	if all.Field != "" {
		queryBuilder = queryBuilder.Where(fmt.Sprintf(`%s ILIKE '%s'`, all.Field, all.Value+"%"))
	}
	if all.OrderBy != "" {
		queryBuilder = queryBuilder.OrderBy(all.OrderBy)
	}
	countBuilder := p.db.Sq.Builder.Select("count(*)").From(p.tableName)
	if !all.IsActive {
		queryBuilder = queryBuilder.Where("deleted_at IS NULL")
		countBuilder = countBuilder.Where("deleted_at IS NULL")
	}
	queryBuilder = queryBuilder.Limit(uint64(all.Limit)).Offset(uint64(offset))
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}
	var doctorWorkHour entity.ListDoctorWorkingHours
	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var Dwhour entity.DoctorWorkingHours
		var deletedAt, updatedAt, startTime, finishTime sql.NullTime
		err = rows.Scan(
			&Dwhour.Id,
			&Dwhour.DoctorId,
			&Dwhour.DayOfWeek,
			&startTime,
			&finishTime,
			&Dwhour.CreatedAt,
			&updatedAt,
			&deletedAt,
		)

		Dwhour.StartTime = startTime.Time.Format("15:04:05")
		Dwhour.FinishTime = finishTime.Time.Format("15:04:05")
		if updatedAt.Valid {
			Dwhour.UpdatedAt = updatedAt.Time
		}
		if deletedAt.Valid {
			Dwhour.DeletedAt = deletedAt.Time
		}
		if err != nil {
			return nil, err
		}
		doctorWorkHour.DoctorWhs = append(doctorWorkHour.DoctorWhs, Dwhour)
	}
	var count int32
	queryCount, _, err := countBuilder.ToSql()
	if err != nil {
		return nil, p.db.Error(err)
	}

	err = p.db.QueryRow(ctx, queryCount).Scan(&count)
	if err != nil {
		return nil, p.db.Error(err)
	}
	doctorWorkHour.Count = count
	return &doctorWorkHour, nil
}

func (p Dwh) UpdateDoctorWorkingHours(ctx context.Context, in *entity.DoctorWorkingHours) (*entity.DoctorWorkingHours, error) {

	ctx, span := otlp.Start(ctx, serviceNameDoctorWorkingHours, serviceNameDoctorWorkingHoursRepoPrefix+"Update")
	span.SetAttributes(attribute.Key("UpdateDoctorWorkingHours").String(string(in.Id)))

	defer span.End()

	data := map[string]any{
		"id":          in.Id,
		"doctor_id":   in.DoctorId,
		"day_of_week": in.DayOfWeek,
		"start_time":  in.StartTime,
		"finish_time": in.FinishTime,
		"updated_at":  time.Now().Add(time.Hour * 5),
	}
	query, args, err := p.db.Sq.Builder.Update(p.tableName).
		SetMap(data).Where(p.db.Sq.Equal("id", in.Id)).Suffix(fmt.Sprintf("RETURNING %s", p.doctorWorkingHoursSelectQueryPrefix())).ToSql()

	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "create"))
	}
	var updatedAt, deletedAt, startTime, finishTime sql.NullTime
	err = p.db.QueryRow(ctx, query, args...).Scan(
		&in.Id,
		&in.DoctorId,
		&in.DayOfWeek,
		&startTime,
		&finishTime,
		&in.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		return nil, p.db.Error(err)
	}
	in.StartTime = startTime.Time.Format("15:04:05")
	in.FinishTime = finishTime.Time.Format("15:04:05")

	if updatedAt.Valid {
		in.UpdatedAt = updatedAt.Time
	}
	if deletedAt.Valid {
		in.DeletedAt = deletedAt.Time
	}
	return in, nil
}

func (p Dwh) DeleteDoctorWorkingHours(ctx context.Context, in *entity.GetReqStr) (bool, error) {

	ctx, span := otlp.Start(ctx, serviceNameDoctorWorkingHours, serviceNameDoctorWorkingHoursRepoPrefix+"Delete")
	span.SetAttributes(attribute.Key("DeleteDoctorWorkingHours").String(in.Value))

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
