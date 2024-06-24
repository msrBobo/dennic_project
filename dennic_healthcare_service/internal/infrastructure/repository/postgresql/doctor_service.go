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
	doctorServicesTableName             = "doctor_service"
	serviceNameDoctorServices           = "doctors_service"
	serviceNameDoctorServicesRepoPrefix = "doctors_service"
)

type Ds struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewDoctorServicesRepo(db *postgres.PostgresDB) *Ds {
	return &Ds{
		tableName: doctorServicesTableName,
		db:        db,
	}
}

func (d Ds) doctorServicesSelectQueryPrefix() string {
	return `
			id,
			doctor_service_order,
			doctor_id,
			specialization_id,
			online_price,
			offline_price,
			name,
			duration,
			created_at,
			updated_at,
			deleted_at
		`
}

func (d Ds) CreateDoctorServices(ctx context.Context, in *entity.DoctorServices) (*entity.DoctorServices, error) {

	ctx, span := otlp.Start(ctx, serviceNameDoctorServices, serviceNameDoctorServicesRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("CreateDoctorServices").String(in.Id))
	defer span.End()

	data := map[string]any{
		"id":                in.Id,
		"doctor_id":         in.DoctorId,
		"specialization_id": in.SpecializationId,
		"online_price":      in.OnlinePrice,
		"offline_price":     in.OfflinePrice,
		"name":              in.Name,
		"duration":          in.Duration,
	}
	query, args, err := d.db.Sq.Builder.Insert(d.tableName).
		SetMap(data).Suffix(fmt.Sprintf("RETURNING %s", d.doctorServicesSelectQueryPrefix())).ToSql()

	if err != nil {
		return nil, d.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", d.tableName, "create"))
	}
	var updatedAt, deletedAt sql.NullTime
	err = d.db.QueryRow(ctx, query, args...).Scan(
		&in.Id,
		&in.Order,
		&in.DoctorId,
		&in.SpecializationId,
		&in.OnlinePrice,
		&in.OfflinePrice,
		&in.Name,
		&in.Duration,
		&in.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	if updatedAt.Valid {
		in.UpdatedAt = updatedAt.Time
	}
	if deletedAt.Valid {
		in.DeletedAt = deletedAt.Time
	}

	if err != nil {
		return nil, d.db.Error(err)
	}
	return in, nil
}

func (d Ds) GetDoctorServiceByID(ctx context.Context, in *entity.GetReqStr) (*entity.DoctorServices, error) {

	ctx, span := otlp.Start(ctx, serviceNameDoctorServices, serviceNameDoctorServicesRepoPrefix+"Get")
	span.SetAttributes(attribute.Key(in.Field).String(in.Value))
	defer span.End()

	var doctorService entity.DoctorServices
	queryBuilder := d.db.Sq.Builder.Select(d.doctorServicesSelectQueryPrefix()).From(d.tableName)
	if !in.IsActive {
		queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	}
	queryBuilder = queryBuilder.Where(d.db.Sq.Equal(in.Field, in.Value))
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}
	var updatedAt, deletedAt sql.NullTime
	err = d.db.QueryRow(ctx, query, args...).Scan(
		&doctorService.Id,
		&doctorService.Order,
		&doctorService.DoctorId,
		&doctorService.SpecializationId,
		&doctorService.OnlinePrice,
		&doctorService.OfflinePrice,
		&doctorService.Name,
		&doctorService.Duration,
		&doctorService.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		return nil, err
	}
	if updatedAt.Valid {
		doctorService.UpdatedAt = updatedAt.Time
	}
	if deletedAt.Valid {
		doctorService.DeletedAt = deletedAt.Time
	}
	return &doctorService, nil
}

func (d Ds) GetAllDoctorServices(ctx context.Context, all *entity.GetAll) (*entity.ListDoctorServices, error) {

	ctx, span := otlp.Start(ctx, serviceNameDoctorServices, serviceNameDoctorServicesRepoPrefix+"Get all")
	span.SetAttributes(attribute.Key(all.Field).String(all.Value))
	defer span.End()

	offset := all.Limit * (all.Page - 1)
	countBuilder := d.db.Sq.Builder.Select("count(*)").From(d.tableName)

	queryBuilder := d.db.Sq.Builder.Select(d.doctorServicesSelectQueryPrefix()).From(d.tableName)
	if all.Field != "" {
		queryBuilder = queryBuilder.Where(fmt.Sprintf(`%s = '%s'`, all.Field, all.Value))
		countBuilder = countBuilder.Where(fmt.Sprintf(`%s = '%s'`, all.Field, all.Value))
	}
	if all.OrderBy != "" {
		queryBuilder = queryBuilder.OrderBy(all.OrderBy)
	}
	if !all.IsActive {
		queryBuilder = queryBuilder.Where("deleted_at IS NULL")
		countBuilder = countBuilder.Where("deleted_at IS NULL")
	}
	queryBuilder = queryBuilder.Limit(uint64(all.Limit)).Offset(uint64(offset))
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}
	var doctorServices entity.ListDoctorServices
	rows, err := d.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var dService entity.DoctorServices
		var deletedAt, updatedAt sql.NullTime
		err = rows.Scan(
			&dService.Id,
			&dService.Order,
			&dService.DoctorId,
			&dService.SpecializationId,
			&dService.OnlinePrice,
			&dService.OfflinePrice,
			&dService.Name,
			&dService.Duration,
			&dService.CreatedAt,
			&updatedAt,
			&deletedAt,
		)
		if updatedAt.Valid {
			dService.UpdatedAt = updatedAt.Time
		}
		if deletedAt.Valid {
			dService.DeletedAt = deletedAt.Time
		}
		if err != nil {
			return nil, err
		}
		doctorServices.DoctorServices = append(doctorServices.DoctorServices, dService)
	}
	var count int64
	queryCount, _, err := countBuilder.ToSql()
	if err != nil {
		return nil, d.db.Error(err)
	}

	err = d.db.QueryRow(ctx, queryCount).Scan(&count)
	if err != nil {
		return nil, d.db.Error(err)
	}
	doctorServices.Count = count
	return &doctorServices, nil
}

func (d Ds) UpdateDoctorServices(ctx context.Context, services *entity.DoctorServices) (*entity.DoctorServices, error) {

	ctx, span := otlp.Start(ctx, serviceNameDoctorServices, serviceNameDoctorServicesRepoPrefix+"Update")
	span.SetAttributes(attribute.Key("UpdateDoctorServices").String(services.Id))
	defer span.End()

	data := map[string]any{
		"doctor_id":         services.DoctorId,
		"specialization_id": services.SpecializationId,
		"online_price":      services.OnlinePrice,
		"offline_price":     services.OfflinePrice,
		"name":              services.Name,
		"duration":          services.Duration,
		"updated_at":        time.Now().Add(time.Hour * 5),
	}
	query, args, err := d.db.Sq.Builder.Update(d.tableName).
		SetMap(data).Where(d.db.Sq.Equal("id", services.Id)).
		Suffix(fmt.Sprintf("RETURNING %s", d.doctorServicesSelectQueryPrefix())).ToSql()

	if err != nil {
		return nil, d.db.ErrSQLBuild(err, d.tableName+" update")
	}
	var deletedAt, updatedAt sql.NullTime
	err = d.db.QueryRow(ctx, query, args...).Scan(
		&services.Id,
		&services.Order,
		&services.DoctorId,
		&services.SpecializationId,
		&services.OnlinePrice,
		&services.OfflinePrice,
		&services.Name,
		&services.Duration,
		&services.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	if updatedAt.Valid {
		services.UpdatedAt = updatedAt.Time
	}
	if deletedAt.Valid {
		services.DeletedAt = deletedAt.Time
	}
	if err != nil {
		return nil, d.db.Error(err)
	}
	return services, nil
}

func (d Ds) DeleteDoctorService(ctx context.Context, in *entity.GetReqStr) (bool, error) {

	ctx, span := otlp.Start(ctx, serviceNameDoctorServices, serviceNameDoctorServicesRepoPrefix+"Delete")
	span.SetAttributes(attribute.Key(in.Field).String(in.Value))
	defer span.End()

	data := map[string]any{
		"deleted_at": time.Now().Add(time.Hour * 5),
	}

	var args []interface{}
	var query string
	var err error
	if in.IsActive {
		query, args, err = d.db.Sq.Builder.Delete(d.tableName).From(d.tableName).
			Where(d.db.Sq.And(d.db.Sq.Equal(in.Field, in.Value))).ToSql()
		if err != nil {
			return false, d.db.ErrSQLBuild(err, d.tableName+" delete")
		}
	} else {
		query, args, err = d.db.Sq.Builder.Update(d.tableName).SetMap(data).
			Where(d.db.Sq.And(d.db.Sq.Equal(in.Field, in.Value), d.db.Sq.Equal("deleted_at", nil))).ToSql()
		if err != nil {
			return false, d.db.ErrSQLBuild(err, d.tableName+" delete")
		}
	}
	resp, err := d.db.Exec(ctx, query, args...)
	if err != nil {
		return false, d.db.Error(err)
	}
	if resp.RowsAffected() > 0 {
		return true, nil
	}
	return false, nil
}
