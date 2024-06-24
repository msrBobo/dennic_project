package usecase

import (
	"context"
	"dennic_session_service/internal/entity"
	"dennic_session_service/internal/infrastructure/repository"
	"dennic_session_service/internal/pkg/otlp"
	"go.opentelemetry.io/otel/attribute"
	"time"
)

const (
	serviceNameSessionUsecase           = "UseCase_session_service"
	serviceNameSessionUsecaseRepoPrefix = "UseCase_session_service"
)

type SessionUsecase interface {
	CreateSession(context.Context, *entity.SessionRequests) (*entity.Session, error)
	GetSessionById(context.Context, *entity.StrReq) (*entity.Session, error)
	DeleteSessionById(context.Context, *entity.StrReq) (*entity.Empty, error)
	DeleteSessionByUserId(context.Context, *entity.StrUserReq) (*entity.Empty, error)
	GetUserSessions(context.Context, *entity.StrUserReq) (*entity.UserSessionsList, error)
}

type newsDepService struct {
	BaseUseCase
	repo       repository.SessionRepo
	ctxTimeout time.Duration
}

func (u newsDepService) CreateSession(ctx context.Context, req *entity.SessionRequests) (*entity.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.ctxTimeout)
	defer cancel()
	ctx, span := otlp.Start(ctx, serviceNameSessionUsecase, serviceNameSessionUsecaseRepoPrefix)
	span.SetAttributes(attribute.Key("CreateSession").String(req.Id))
	defer span.End()
	return u.repo.CreateSession(ctx, req)
}

func (u newsDepService) GetSessionById(ctx context.Context, req *entity.StrReq) (*entity.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.ctxTimeout)
	defer cancel()
	ctx, span := otlp.Start(ctx, serviceNameSessionUsecase, serviceNameSessionUsecaseRepoPrefix)
	span.SetAttributes(attribute.Key("CreateSession").String(req.Id))
	defer span.End()
	return u.repo.GetSessionById(ctx, req)
}

func (u newsDepService) DeleteSessionById(ctx context.Context, req *entity.StrReq) (*entity.Empty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.ctxTimeout)
	defer cancel()
	ctx, span := otlp.Start(ctx, serviceNameSessionUsecase, serviceNameSessionUsecaseRepoPrefix)
	span.SetAttributes(attribute.Key("DeleteSessionById").String(req.Id))
	defer span.End()
	return u.repo.DeleteSessionById(ctx, req)
}

func (u newsDepService) DeleteSessionByUserId(ctx context.Context, req *entity.StrUserReq) (*entity.Empty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.ctxTimeout)
	defer cancel()
	ctx, span := otlp.Start(ctx, serviceNameSessionUsecase, serviceNameSessionUsecaseRepoPrefix)
	span.SetAttributes(attribute.Key("DeleteSessionByUserId").String(req.UserId))
	defer span.End()
	return u.repo.DeleteSessionByUserId(ctx, req)
}

func (u newsDepService) GetUserSessions(ctx context.Context, req *entity.StrUserReq) (*entity.UserSessionsList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.ctxTimeout)
	defer cancel()
	ctx, span := otlp.Start(ctx, serviceNameSessionUsecase, serviceNameSessionUsecaseRepoPrefix)
	span.SetAttributes(attribute.Key("GetUserSessions").String(req.UserId))
	defer span.End()
	return u.repo.GetUserSessions(ctx, req)
}

func NewSessionService(ctxTimeout time.Duration, repo repository.SessionRepo) newsDepService {
	return newsDepService{
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}
