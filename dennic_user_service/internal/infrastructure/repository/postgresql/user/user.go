package postgresql

import (
	"context"
	"database/sql"
	"dennic_user_service/internal/entity"
	"dennic_user_service/internal/pkg/otlp"
	"dennic_user_service/internal/pkg/postgres"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

const (
	userTableName      = "users"
	userServiceName    = "userService"
	userSpanRepoPrefix = "userRepo"
)

type userRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewUserRepo(db *postgres.PostgresDB) *userRepo {
	return &userRepo{
		tableName: userTableName,
		db:        db,
	}
}

func (p *userRepo) userSelectQueryPrefix() string {
	return `id,
			user_order,
			first_name,
			last_name,
			birth_date,
			phone_number,
			password,
			gender,
			image_url,
			created_at,
			updated_at,
			deleted_at`
}

func (p userRepo) Create(ctx context.Context, user *entity.User) error {
	ctx, span := otlp.Start(ctx, userServiceName, userSpanRepoPrefix+"Create")
	defer span.End()
	var userIDKey = attribute.Key("user_id")
	span.SetAttributes(userIDKey.String(user.Id))
	data := map[string]any{
		"id":            user.Id,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"birth_date":    user.BirthDate,
		"phone_number":  user.PhoneNumber,
		"password":      user.Password,
		"gender":        user.Gender,
		"refresh_token": user.RefreshToken,
		"image_url":     user.ImageUrl,
	}

	query, args, err := p.db.Sq.Builder.Insert(p.tableName).SetMap(data).ToSql()
	if err != nil {
		return p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "create"))
	}

	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		return p.db.Error(err)
	}

	return nil
}

func (p userRepo) Get(ctx context.Context, req *entity.FieldValueReq) (*entity.User, error) {
	ctx, span := otlp.Start(ctx, userServiceName, userSpanRepoPrefix+"Get")
	defer span.End()
	var (
		user entity.User
	)

	toSql := p.db.Sq.Builder.
		Select(p.userSelectQueryPrefix()).
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
		birthDate sql.NullTime
		updatedAt sql.NullTime
		deletedAt sql.NullTime
	)
	if err = p.db.QueryRow(ctx, toSqls, args...).Scan(
		&user.Id,
		&user.UserOrder,
		&user.FirstName,
		&user.LastName,
		&birthDate,
		&user.PhoneNumber,
		&user.Password,
		&user.Gender,
		&user.ImageUrl,
		&user.CreatedAt,
		&updatedAt,
		&deletedAt,
	); err != nil {
		return nil, p.db.Error(err)
	}

	if birthDate.Valid {
		user.BirthDate = birthDate.Time.String()
	}
	if updatedAt.Valid {
		user.UpdatedAt = updatedAt.Time
	}
	if deletedAt.Valid {
		user.DeletedAt = deletedAt.Time
	}
	var userIDKey = attribute.Key("user_id")
	span.SetAttributes(userIDKey.String(user.Id))
	return &user, nil
}

func (p userRepo) List(ctx context.Context, req *entity.GetAllReq) ([]*entity.User, error) {
	ctx, span := otlp.Start(ctx, userServiceName, userSpanRepoPrefix+"List")
	defer span.End()
	var (
		users []*entity.User
	)
	toSql := p.db.Sq.Builder.
		Select(p.userSelectQueryPrefix()).
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
	countBuilder := p.db.Sq.Builder.Select("count(*)").From(userTableName)

	if !req.DeleteStatus {
		toSql = toSql.Where(p.db.Sq.Equal("deleted_at", nil))
		countBuilder = countBuilder.Where("deleted_at IS NULL")
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
		birthDate sql.NullTime
		updatedAt sql.NullTime
		deletedAt sql.NullTime
		count     int64
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
		var user entity.User
		if err = rows.Scan(
			&user.Id,
			&user.UserOrder,
			&user.FirstName,
			&user.LastName,
			&birthDate,
			&user.PhoneNumber,
			&user.Password,
			&user.Gender,
			&user.ImageUrl,
			&user.CreatedAt,
			&updatedAt,
			&deletedAt,
		); err != nil {
			return nil, p.db.Error(err)
		}

		if birthDate.Valid {
			user.BirthDate = birthDate.Time.String()
		}
		if updatedAt.Valid {
			user.UpdatedAt = updatedAt.Time
		}
		if deletedAt.Valid {
			user.DeletedAt = deletedAt.Time
		}
		user.Count = count
		users = append(users, &user)
	}

	return users, nil
}

func (p userRepo) Update(ctx context.Context, user *entity.User) error {
	ctx, span := otlp.Start(ctx, userServiceName, userSpanRepoPrefix+"Update")
	defer span.End()
	clauses := map[string]any{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"birth_date": user.BirthDate,
		"gender":     user.Gender,
		"image_url":  user.ImageUrl,
		"updated_at": user.UpdatedAt,
	}

	updateBuilder := p.db.Sq.Builder.
		Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("id", user.Id))

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

func (p *userRepo) Delete(ctx context.Context, req *entity.FieldValueReq) (*entity.CheckDeleteResp, error) {
	ctx, span := otlp.Start(ctx, userServiceName, userSpanRepoPrefix+"Delete")
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

func (p *userRepo) CheckField(ctx context.Context, req *entity.CheckFieldReq) (*entity.CheckFieldResp, error) {
	ctx, span := otlp.Start(ctx, userServiceName, userSpanRepoPrefix+"CheckField")
	defer span.End()
	query := fmt.Sprintf(
		`SELECT count(1) 
		FROM users WHERE %s = $1 
		AND deleted_at IS NULL`, req.Field)

	var isExists int

	row := p.db.QueryRow(ctx, query, req.Value)
	if err := row.Scan(&isExists); err != nil {
		return nil, err
	}

	if isExists == 1 {
		return &entity.CheckFieldResp{
			Status: true,
		}, nil
	}

	return &entity.CheckFieldResp{
		Status: false,
	}, nil
}

func (p *userRepo) ChangePassword(ctx context.Context, req *entity.ChangeUserPasswordReq) (*entity.ChangePasswordResp, error) {
	ctx, span := otlp.Start(ctx, userServiceName, userSpanRepoPrefix+"ChangePassword")
	defer span.End()
	query := `
		UPDATE users 
		SET password = $1 
		WHERE phone_number = $2 
		AND deleted_at IS NULL
	`
	resp, err := p.db.Exec(ctx, query, req.Password, req.PhoneNumber)
	if err != nil {
		return nil, err
	}
	if resp.RowsAffected() == 0 {
		return &entity.ChangePasswordResp{Status: false}, nil
	}
	return &entity.ChangePasswordResp{Status: true}, nil
}

func (p *userRepo) UpdateRefreshToken(ctx context.Context, req *entity.UpdateRefreshTokenReq) (*entity.UpdateRefreshTokenResp, error) {
	ctx, span := otlp.Start(ctx, userServiceName, userSpanRepoPrefix+"UpdateRefreshToken")
	defer span.End()
	query := `
			UPDATE users 
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
