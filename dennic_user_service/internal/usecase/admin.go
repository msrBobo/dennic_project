package usecase

import (
	"context"
	"dennic_user_service/internal/entity"
	"dennic_user_service/internal/infrastructure/repository"
	"dennic_user_service/internal/pkg/otlp"
	"time"
)

const (
	AdminServiceName = "adminService"
	AdinSpanName    = "adminUsecase"
)

type AdminStorageI interface {
	Create(ctx context.Context, admin *entity.Admin) (string, error)
	Get(ctx context.Context, req *entity.FieldValueReq) (*entity.Admin, error)
	List(ctx context.Context, req *entity.GetAllReq) ([]*entity.Admin, error)
	Update(ctx context.Context, kyc *entity.Admin) error
	Delete(ctx context.Context, req *entity.FieldValueReq) (*entity.CheckDeleteResp, error)
	CheckField(ctx context.Context, req *entity.CheckFieldReq) (*entity.CheckFieldResp, error)
	ChangePassword(ctx context.Context, req *entity.ChangeAdminPasswordReq) (*entity.ChangeAdminPasswordResp, error)
	UpdateRefreshToken(ctx context.Context, req *entity.UpdateRefreshTokenReq) (*entity.UpdateRefreshTokenResp, error)
}


type adminService struct {
	repo       repository.AdminStorageI
	ctxTimeout time.Duration
}

func NewAdminService(ctxTimeout time.Duration, repo repository.AdminStorageI) adminService {
	return adminService{
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}

func (a adminService) Create(ctx context.Context, admin *entity.Admin) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, AdminServiceName, AdinSpanName+"Create")
	defer span.End()
	adminId := admin.Id

	return adminId, a.repo.Create(ctx, admin)
}

func (a adminService) Get(ctx context.Context, req *entity.FieldValueReq) (*entity.Admin, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, AdminServiceName, AdinSpanName+"Get")
	defer span.End()

	return a.repo.Get(ctx, req)
}

func (a adminService) List(ctx context.Context, req *entity.GetAllReq) ([]*entity.Admin, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, AdminServiceName, AdinSpanName+"List")
	defer span.End()

	return a.repo.List(ctx, req)
}

func (a adminService) Update(ctx context.Context, req *entity.Admin) error {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, AdminServiceName, AdinSpanName+"Update")
	defer span.End()

	return a.repo.Update(ctx, req)
}

func (a adminService) Delete(ctx context.Context, req *entity.FieldValueReq) (*entity.CheckDeleteResp, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, AdminServiceName, AdinSpanName+"Delete")
	defer span.End()

	return a.repo.Delete(ctx, req)
}

func (a adminService) CheckField(ctx context.Context, req *entity.CheckFieldReq) (*entity.CheckFieldResp, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, AdminServiceName, AdinSpanName+"CheckField")
	defer span.End()

	return a.repo.CheckField(ctx, req)
}

func (a adminService) ChangePassword(ctx context.Context, req *entity.ChangeAdminPasswordReq) (*entity.ChangeAdminPasswordResp, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, AdminServiceName, AdinSpanName+"ChangePassword")
	defer span.End()

	return a.repo.ChangePassword(ctx, req)
}

func (a adminService) UpdateRefreshToken(ctx context.Context, req *entity.UpdateRefreshTokenReq) (*entity.UpdateRefreshTokenResp, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, AdminServiceName, AdinSpanName+"UpdateRefreshToken")
	defer span.End()

	return a.repo.UpdateRefreshToken(ctx, req)
}
