package repository

import (
	"Healthcare_Evrone/internal/entity"
	"context"
)

type ReasonRepository interface {
	CreateReasons(ctx context.Context, in *entity.Reasons) (*entity.Reasons, error)
	GetReasonsById(ctx context.Context, in *entity.GetReqStr) (*entity.Reasons, error)
	GetAllReasons(context.Context, *entity.GetAllReas) (*entity.ListReasons, error)
	UpdateReasons(context.Context, *entity.Reasons) (*entity.Reasons, error)
	DeleteReasons(context.Context, *entity.GetReqStr) (*entity.StatusReasons, error)
}
