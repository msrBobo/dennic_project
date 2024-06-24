package repository

import (
	"context"
	"dennic_user_service/internal/entity"
)

type UserStorageI interface {
	Create(ctx context.Context, user *entity.User) (error)
	Get(ctx context.Context, req *entity.FieldValueReq) (*entity.User, error)
	List(ctx context.Context, req *entity.GetAllReq) ([]*entity.User, error)
	Update(ctx context.Context, kyc *entity.User) error
	Delete(ctx context.Context, req *entity.FieldValueReq) (*entity.CheckDeleteResp, error)
	CheckField(ctx context.Context, req *entity.CheckFieldReq) (*entity.CheckFieldResp, error)
	ChangePassword(ctx context.Context, req *entity.ChangeUserPasswordReq) (*entity.ChangePasswordResp, error)
	UpdateRefreshToken(ctx context.Context, req *entity.UpdateRefreshTokenReq) (*entity.UpdateRefreshTokenResp, error)
}
