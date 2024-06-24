package usecase

import (
	"context"
	"dennic_user_service/internal/entity"
	"dennic_user_service/internal/infrastructure/repository"
	"dennic_user_service/internal/pkg/otlp"
	"time"
)

const (
	UserServiceName = "userService"
	UserSpanName    = "userUsecase"
)

type UserStorageI interface {
	Create(ctx context.Context, user *entity.User) (string, error)
	Get(ctx context.Context, req *entity.FieldValueReq) (*entity.User, error)
	List(ctx context.Context, req *entity.GetAllReq) ([]*entity.User, error)
	Update(ctx context.Context, kyc *entity.User) error
	Delete(ctx context.Context, req *entity.FieldValueReq) (*entity.CheckDeleteResp, error)
	CheckField(ctx context.Context, req *entity.CheckFieldReq) (*entity.CheckFieldResp, error)
	ChangePassword(ctx context.Context, req *entity.ChangeUserPasswordReq) (*entity.ChangePasswordResp, error)
	UpdateRefreshToken(ctx context.Context, req *entity.UpdateRefreshTokenReq) (*entity.UpdateRefreshTokenResp, error)
}

type userService struct {
	repo       repository.UserStorageI
	ctxTimeout time.Duration
}

func NewUserService(ctxTimeout time.Duration, repo repository.UserStorageI) userService {
	return userService{
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}

func (u userService) Create(ctx context.Context, user *entity.User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, UserServiceName, UserSpanName+"Create")
	defer span.End()
	userId := user.Id
	return userId, u.repo.Create(ctx, user)
}

func (u userService) Get(ctx context.Context, req *entity.FieldValueReq) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, UserServiceName, UserSpanName+"Get")
	defer span.End()

	return u.repo.Get(ctx, req)
}

func (u userService) List(ctx context.Context, req *entity.GetAllReq) ([]*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, UserServiceName, UserSpanName+"List")
	defer span.End()

	return u.repo.List(ctx, req)
}

func (u userService) Update(ctx context.Context, articleCategory *entity.User) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, UserServiceName, UserSpanName+"Update")
	defer span.End()

	return u.repo.Update(ctx, articleCategory)
}

func (u userService) Delete(ctx context.Context, req *entity.FieldValueReq) (*entity.CheckDeleteResp, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, UserServiceName, UserSpanName+"Delete")
	defer span.End()

	return u.repo.Delete(ctx, req)
}

func (u userService) CheckField(ctx context.Context, req *entity.CheckFieldReq) (*entity.CheckFieldResp, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, UserServiceName, UserSpanName+"CheckField")
	defer span.End()

	return u.repo.CheckField(ctx, req)
}

func (u userService) ChangePassword(ctx context.Context, req *entity.ChangeUserPasswordReq) (*entity.ChangePasswordResp, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, UserServiceName, UserSpanName+"ChangePassword")
	defer span.End()

	return u.repo.ChangePassword(ctx, req)
}

func (u userService) UpdateRefreshToken(ctx context.Context, req *entity.UpdateRefreshTokenReq) (*entity.UpdateRefreshTokenResp, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, UserServiceName, UserSpanName+"UpdateRefreshToken")
	defer span.End()

	return u.repo.UpdateRefreshToken(ctx, req)
}
