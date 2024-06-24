package postgresql

import (
	"context"
	"database/sql"
	"dennic_session_service/internal/entity"
	"dennic_session_service/internal/pkg/otlp"
	"dennic_session_service/internal/pkg/postgres"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"time"
)

type SessionRepository struct {
	tableName string
	db        *postgres.PostgresDB
}

const (
	sessionTableName             = "sessions"
	serviceNameSession           = "postgres_session_service"
	serviceNameSessionRepoPrefix = "postgres_session_service"
)

func NewSessionRepository(db *postgres.PostgresDB) *SessionRepository {
	return &SessionRepository{
		tableName: sessionTableName,
		db:        db,
	}
}

func (s *SessionRepository) sessionSelectQueryPrefix() string {
	return `id,
			session_order,
			ip_address,
			user_id,
			fcm_token,
			platform_name,
			platform_type,
			login_at,
			created_at,
			updated_at,
			deleted_at`
}

func (s *SessionRepository) CreateSession(ctx context.Context, session *entity.SessionRequests) (*entity.Session, error) {
	ctx, span := otlp.Start(ctx, serviceNameSession, serviceNameSessionRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("CreateSession").String(session.Id))
	defer span.End()

	data := map[string]any{
		"id":            session.Id,
		"ip_address":    session.IpAddress,
		"user_id":       session.UserId,
		"fcm_token":     session.FcmToken,
		"platform_name": session.PlatformName,
		"platform_type": session.PlatformType,
	}
	query, args, err := s.db.Sq.Builder.Insert(s.tableName).SetMap(data).Suffix(fmt.Sprintf("RETURNING %s", s.sessionSelectQueryPrefix())).ToSql()
	if err != nil {

		return nil, s.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", s.tableName, " create"))
	}

	var resp entity.Session
	var updatedAt, deletedAt sql.NullTime

	err = s.db.QueryRow(ctx, query, args...).Scan(
		&resp.Id,
		&resp.Order,
		&resp.IpAddress,
		&resp.UserId,
		&resp.FcmToken,
		&resp.PlatformName,
		&resp.PlatformType,
		&resp.LoginAt,
		&resp.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		return nil, s.db.Error(err)
	}

	if updatedAt.Valid {
		resp.UpdatedAt = updatedAt.Time
	}
	if deletedAt.Valid {
		resp.DeletedAt = deletedAt.Time
	}
	return &resp, nil
}

func (s *SessionRepository) GetSessionById(ctx context.Context, req *entity.StrReq) (*entity.Session, error) {
	ctx, span := otlp.Start(ctx, serviceNameSession, serviceNameSessionRepoPrefix+"Get")
	span.SetAttributes(attribute.Key("GetSessionById").String(req.Id))
	defer span.End()

	query, args, err := s.db.Sq.Builder.Select(s.sessionSelectQueryPrefix()).From(s.tableName).
		Where(s.db.Sq.Equal("id", req.Id)).ToSql()
	if err != nil {
		return nil, s.db.Error(err)
	}
	var resp entity.Session
	var updatedAt, deletedAt sql.NullTime

	err = s.db.QueryRow(ctx, query, args...).Scan(
		&resp.Id,
		&resp.Order,
		&resp.IpAddress,
		&resp.UserId,
		&resp.FcmToken,
		&resp.PlatformName,
		&resp.PlatformType,
		&resp.LoginAt,
		&resp.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		return nil, s.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", s.tableName, " get"))
	}

	if updatedAt.Valid {
		resp.UpdatedAt = updatedAt.Time
	}
	if deletedAt.Valid {
		resp.DeletedAt = deletedAt.Time
	}
	return &resp, nil
}

func (s *SessionRepository) DeleteSessionById(ctx context.Context, req *entity.StrReq) (*entity.Empty, error) {
	ctx, span := otlp.Start(ctx, serviceNameSession, serviceNameSessionRepoPrefix+"Delete")
	span.SetAttributes(attribute.Key("DeleteSessionById").String(req.Id))
	defer span.End()
	data := map[string]any{
		"deleted_at": time.Now().Add(time.Hour * 5),
	}
	query, args, err := s.db.Sq.Builder.Update(s.tableName).SetMap(data).Where(s.db.Sq.Equal("id", req.Id)).ToSql()
	if err != nil {
		return nil, s.db.ErrSQLBuild(err, s.tableName+" delete")
	}

	_, err = s.db.Exec(ctx, query, args...)
	if err != nil {

		return nil, s.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", s.tableName, " delete"))
	}
	return &entity.Empty{}, nil
}

func (s *SessionRepository) DeleteSessionByUserId(ctx context.Context, req *entity.StrUserReq) (*entity.Empty, error) {
	ctx, span := otlp.Start(ctx, serviceNameSession, serviceNameSessionRepoPrefix+"Delete")
	span.SetAttributes(attribute.Key("DeleteSessionById").String(req.UserId))
	defer span.End()
	data := map[string]any{
		"deleted_at": time.Now().Add(time.Hour * 5),
	}
	query, args, err := s.db.Sq.Builder.Update(s.tableName).SetMap(data).Where(s.db.Sq.Equal("user_id", req.UserId)).ToSql()
	if err != nil {
		return nil, s.db.ErrSQLBuild(err, s.tableName+" delete")
	}

	_, err = s.db.Exec(ctx, query, args...)
	if err != nil {

		return nil, s.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", s.tableName, " delete"))
	}
	return &entity.Empty{}, nil
}

func (s *SessionRepository) GetUserSessions(ctx context.Context, req *entity.StrUserReq) (*entity.UserSessionsList, error) {
	ctx, span := otlp.Start(ctx, serviceNameSession, serviceNameSessionRepoPrefix+"Get")
	span.SetAttributes(attribute.Key("GetUserSessions").String(req.UserId))
	defer span.End()
	countBuilder := s.db.Sq.Builder.Select("count(*)").From(s.tableName)
	queryBuilder := s.db.Sq.Builder.Select(s.sessionSelectQueryPrefix()).From(s.tableName)
	if !req.IsActive {
		queryBuilder = queryBuilder.Where("deleted_at IS NULL")
		countBuilder = countBuilder.Where("deleted_at IS NULL")
	}

	query, args, err := queryBuilder.Where(s.db.Sq.Equal("user_id", req.UserId)).ToSql()
	if err != nil {
		return nil, s.db.Error(err)
	}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, s.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", s.tableName, " get"))
	}

	var userSessions entity.UserSessionsList
	for rows.Next() {
		var resp entity.Session
		var updatedAt, deletedAt sql.NullTime
		err := rows.Scan(
			&resp.Id,
			&resp.Order,
			&resp.IpAddress,
			&resp.UserId,
			&resp.FcmToken,
			&resp.PlatformName,
			&resp.PlatformType,
			&resp.LoginAt,
			&resp.CreatedAt,
			&updatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, s.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", s.tableName, " get"))
		}

		if updatedAt.Valid {
			resp.UpdatedAt = updatedAt.Time
		}
		if deletedAt.Valid {
			resp.DeletedAt = deletedAt.Time
		}
		userSessions.UserSessions = append(userSessions.UserSessions, &resp)
		var count int32
		queryCount, _, err := countBuilder.ToSql()
		err = s.db.QueryRow(ctx, queryCount).Scan(&count)
		if err != nil {
			return nil, s.db.Error(err)
		}
		userSessions.Count = count
	}
	return &userSessions, nil
}
