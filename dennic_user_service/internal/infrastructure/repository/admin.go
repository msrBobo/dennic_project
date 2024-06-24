package repository

import (
	"context"
	"dennic_user_service/internal/entity"
)

type AdminStorageI interface {
	Create(ctx context.Context, admin *entity.Admin) (error)
	Get(ctx context.Context, req *entity.FieldValueReq) (*entity.Admin, error)
	List(ctx context.Context, req *entity.GetAllReq) ([]*entity.Admin, error)
	Update(ctx context.Context, kyc *entity.Admin) error
	Delete(ctx context.Context, req *entity.FieldValueReq) (*entity.CheckDeleteResp, error)
	CheckField(ctx context.Context, req *entity.CheckFieldReq) (*entity.CheckFieldResp, error)
	ChangePassword(ctx context.Context, req *entity.ChangeAdminPasswordReq) (*entity.ChangeAdminPasswordResp, error)
	UpdateRefreshToken(ctx context.Context, req *entity.UpdateRefreshTokenReq) (*entity.UpdateRefreshTokenResp, error)
}