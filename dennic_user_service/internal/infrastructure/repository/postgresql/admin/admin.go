package postgresql

import (
	"context"
	"database/sql"
	"dennic_user_service/internal/entity"
	"dennic_user_service/internal/pkg/otlp"
	"dennic_user_service/internal/pkg/postgres"
	"fmt"
	"time"
)

const (
	adminTableName      = "admins"
	adminServiceName    = "adminService"
	adminSpanRepoPrefix = "adminRepo"
)

type adminRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewAdminRepo(db *postgres.PostgresDB) *adminRepo {
	return &adminRepo{
		tableName: adminTableName,
		db:        db,
	}
}

func (p *adminRepo) adminSelectQueryPrefix() string {
	return `id,
			admin_order,
			role,
			first_name,
			last_name,
			birth_date,
			phone_number,
			email,
			password,
			gender,
			salary,
			biography,
			start_work_year,
			end_work_year,
			work_years,
			image_url,
			created_at,
			updated_at,
			deleted_at`
}

func (p adminRepo) Create(ctx context.Context, admin *entity.Admin) error {
	ctx, span := otlp.Start(ctx, adminServiceName, adminSpanRepoPrefix+"Create")
	defer span.End()
	data := map[string]any{
		"id":              admin.Id,
		"role":            admin.Role,
		"first_name":      admin.FirstName,
		"last_name":       admin.LastName,
		"birth_date":      admin.BirthDate,
		"phone_number":    admin.PhoneNumber,
		"email":           admin.Email,
		"password":        admin.Password,
		"gender":          admin.Gender,
		"salary":          admin.Salary,
		"biography":       admin.Biography,
		"start_work_year": admin.StartWorkYear,
		"work_years":      admin.WorkYears,
		"refresh_token":   admin.RefreshToken,
		"image_url":       admin.ImageUrl,
	}

	query, args, err := p.db.Sq.Builder.Insert(p.tableName).SetMap(data).ToSql()
	fmt.Println(err)

	if err != nil {
		return p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "create"))
	}
	_, err = p.db.Exec(ctx, query, args...)
	fmt.Println(err)

	if err != nil {
		return p.db.Error(err)
	}

	return nil
}

func (p adminRepo) Get(ctx context.Context, req *entity.FieldValueReq) (*entity.Admin, error) {
	ctx, span := otlp.Start(ctx, adminServiceName, adminSpanRepoPrefix+"Get")
	defer span.End()

	var (
		admin entity.Admin
	)

	toSql := p.db.Sq.Builder.
		Select(p.adminSelectQueryPrefix()).
		From(p.tableName).
		Where(p.db.Sq.Equal(req.Field, req.Value))

	if !req.DeleteStatus {
		toSql = toSql.Where(p.db.Sq.Equal("deleted_at", nil))
	}

	toSqls, args, err := toSql.ToSql()

	if err != nil {
		return nil, err
	}

	var (
		birthDate       sql.NullString
		updatedAt       sql.NullTime
		start_work_year sql.NullString
		end_work_year   sql.NullString
		deletedAt       sql.NullTime
	)
	if err = p.db.QueryRow(ctx, toSqls, args...).Scan(
		&admin.Id,
		&admin.AdminOrder,
		&admin.Role,
		&admin.FirstName,
		&admin.LastName,
		&birthDate,
		&admin.PhoneNumber,
		&admin.Email,
		&admin.Password,
		&admin.Gender,
		&admin.Salary,
		&admin.Biography,
		&start_work_year,
		&end_work_year,
		&admin.WorkYears,
		&admin.ImageUrl,
		&admin.CreatedAt,
		&updatedAt,
		&deletedAt,
	); err != nil {
		return nil, p.db.Error(err)
	}

	if updatedAt.Valid {
		admin.UpdatedAt = updatedAt.Time
	}
	if birthDate.Valid {
		admin.BirthDate = birthDate.String
	}
	if start_work_year.Valid {
		admin.StartWorkYear = start_work_year.String
	}
	if end_work_year.Valid {
		admin.EndWorkYear = end_work_year.String
	}
	if deletedAt.Valid {
		admin.DeletedAt = deletedAt.Time
	}

	return &admin, nil
}

func (p adminRepo) List(ctx context.Context, req *entity.GetAllReq) ([]*entity.Admin, error) {
	ctx, span := otlp.Start(ctx, adminServiceName, adminSpanRepoPrefix+"List")
	defer span.End()

	var (
		admins []*entity.Admin
	)

	toSql := p.db.Sq.Builder.
		Select(p.adminSelectQueryPrefix()).
		From(p.tableName)

	if req.Page >= 1 && req.Limit >= 1 {
		toSql = toSql.
			Limit(req.Limit).
			Offset(req.Limit * (req.Page - 1))
	}
	if req.Value != "" {
		toSql = toSql.Where(p.db.Sq.ILike(req.Field, req.Value+"%"))
	}
	if req.OrderBy != "" {
		toSql = toSql.OrderBy(req.OrderBy)
	}
	countBuilder := p.db.Sq.Builder.Select("count(*)").From(adminTableName)
	if !req.DeleteStatus {
		countBuilder = countBuilder.Where("deleted_at IS NULL")
		toSql = toSql.Where(p.db.Sq.Equal("deleted_at", nil))
	}
	toSqls, args, err := toSql.ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := p.db.Query(ctx, toSqls, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}
	defer rows.Close()
	var (
		birthDate       sql.NullString
		updatedAt       sql.NullTime
		start_work_year sql.NullString
		end_work_year   sql.NullString
		count           int64
		deletedAt       sql.NullTime
	)
	queryCount, _, err := countBuilder.ToSql()
	if err != nil {
		return nil, p.db.Error(err)
	}
	err = p.db.QueryRow(ctx, queryCount).Scan(&count)
	if err != nil {
		return nil, p.db.Error(err)
	}
	for rows.Next() {
		var admin entity.Admin
		if err = rows.Scan(
			&admin.Id,
			&admin.AdminOrder,
			&admin.Role,
			&admin.FirstName,
			&admin.LastName,
			&birthDate,
			&admin.PhoneNumber,
			&admin.Email,
			&admin.Password,
			&admin.Gender,
			&admin.Salary,
			&admin.Biography,
			&start_work_year,
			&end_work_year,
			&admin.WorkYears,
			&admin.ImageUrl,
			&admin.CreatedAt,
			&updatedAt,
			&deletedAt,
		); err != nil {
			return nil, p.db.Error(err)
		}

		if updatedAt.Valid {
			admin.UpdatedAt = updatedAt.Time
		}
		if birthDate.Valid {
			admin.BirthDate = birthDate.String
		}
		if start_work_year.Valid {
			admin.StartWorkYear = start_work_year.String
		}
		if end_work_year.Valid {
			admin.EndWorkYear = end_work_year.String
		}
		if deletedAt.Valid {
			admin.DeletedAt = deletedAt.Time
		}
		admins = append(admins, &admin)
		admin.Count = count
	}

	return admins, nil
}

func (p *adminRepo) Update(ctx context.Context, admin *entity.Admin) error {
	ctx, span := otlp.Start(ctx, adminServiceName, adminSpanRepoPrefix+"Update")
	defer span.End()

	clauses := map[string]interface{}{
		"first_name":      admin.FirstName,
		"last_name":       admin.LastName,
		"birth_date":      admin.BirthDate,
		"gender":          admin.Gender,
		"salary":          admin.Salary,
		"biography":       admin.Biography,
		"start_work_year": admin.StartWorkYear,
		"work_years":      admin.WorkYears,
		"image_url":       admin.ImageUrl,
		"updated_at":      admin.UpdatedAt,
	}

	if admin.EndWorkYear != "" {
		clauses["end_work_year"] = admin.EndWorkYear
	}

	updateBuilder := p.db.Sq.Builder.
		Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("id", admin.Id))

	updateBuilder = updateBuilder.Where("deleted_at IS NULL")

	sqlStr, args, err := updateBuilder.ToSql()
	if err != nil {
		return p.db.ErrSQLBuild(err, p.tableName+" update")
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return p.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return p.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}

func (p *adminRepo) Delete(ctx context.Context, req *entity.FieldValueReq) (*entity.CheckDeleteResp, error) {
	ctx, span := otlp.Start(ctx, adminServiceName, adminSpanRepoPrefix+"Delete")
	defer span.End()

	if !req.DeleteStatus {
		toSql, args, err := p.db.Sq.Builder.
			Update(p.tableName).
			Set("deleted_at", time.Now().Add(time.Hour*5)).
			Where(p.db.Sq.And(
				p.db.Sq.Equal("deleted_at", nil),
				p.db.Sq.Equal(req.Field, req.Value),
			)).
			ToSql()
		if err != nil {
			return nil, err
		}

		resp, err := p.db.Exec(ctx, toSql, args...)
		if err != nil {
			return nil, err
		}
		if resp.RowsAffected() > 0 {
			return &entity.CheckDeleteResp{Status: true}, nil
		}
		return &entity.CheckDeleteResp{Status: false}, nil
	}

	// If DeleteStatus is true, then perform a delete operation
	toSql, args, err := p.db.Sq.Builder.
		Delete(p.tableName).
		Where(p.db.Sq.Equal(req.Field, req.Value)).
		ToSql()
	if err != nil {
		return nil, err
	}

	resp, err := p.db.Exec(ctx, toSql, args...)
	if err != nil {
		return nil, err
	}

	if resp.RowsAffected() > 0 {
		return &entity.CheckDeleteResp{Status: true}, nil
	}
	return &entity.CheckDeleteResp{Status: false}, nil
}

func (p *adminRepo) CheckField(ctx context.Context, req *entity.CheckFieldReq) (*entity.CheckFieldResp, error) {
	ctx, span := otlp.Start(ctx, adminServiceName, adminSpanRepoPrefix+"CheckField")
	defer span.End()
	query := fmt.Sprintf(`
		SELECT count(1) 
			FROM admins WHERE %s = $1 AND 
			deleted_at IS NULL`, req.Field)

	var isExists int

	row := p.db.QueryRow(ctx, query, req.Value)
	if err := row.Scan(&isExists); err != nil {
		return nil, err
	}
	if isExists > 0 {
		return &entity.CheckFieldResp{
			Status: true,
		}, nil
	}

	return &entity.CheckFieldResp{
		Status: false,
	}, nil
}

func (p *adminRepo) ChangePassword(ctx context.Context, req *entity.ChangeAdminPasswordReq) (*entity.ChangeAdminPasswordResp, error) {
	ctx, span := otlp.Start(ctx, adminServiceName, adminSpanRepoPrefix+"ChangePassword")
	defer span.End()
	query := `
		UPDATE admins
		SET password = $1
		WHERE (email = $2 OR phone_number = $3)
		AND deleted_at IS NULL
	`

	resp, err := p.db.Exec(ctx, query, req.Password, req.Email, req.PhoneNumber)
	if err != nil {
		return nil, err
	}

	if resp.RowsAffected() == 0 {
		return &entity.ChangeAdminPasswordResp{Status: false}, nil
	}

	return &entity.ChangeAdminPasswordResp{Status: true}, nil
}

func (p *adminRepo) UpdateRefreshToken(ctx context.Context, req *entity.UpdateRefreshTokenReq) (*entity.UpdateRefreshTokenResp, error) {
	ctx, span := otlp.Start(ctx, adminServiceName, adminSpanRepoPrefix+"UpdateRefreshToken")
	defer span.End()
	query := `
			UPDATE admins 
			SET refresh_token = $1 
			WHERE id = $2 AND 
			deleted_at IS NULL`

	resp, err := p.db.Exec(ctx, query, req.RefreshToken, req.Id)
	if err != nil {
		return nil, err
	}
	if resp.RowsAffected() == 0 {
		return &entity.UpdateRefreshTokenResp{Status: false}, nil
	}

	return &entity.UpdateRefreshTokenResp{Status: true}, nil
}
