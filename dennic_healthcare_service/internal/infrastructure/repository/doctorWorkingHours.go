package repository

import (
	"Healthcare_Evrone/internal/entity"
	"context"
)

type DoctorWorkingHoursRepository interface {
	CreateDoctorWorkingHours(ctx context.Context, in *entity.DoctorWorkingHours) (*entity.DoctorWorkingHours, error)
	GetDoctorWorkingHoursById(ctx context.Context, in *entity.GetRequest) (*entity.DoctorWorkingHours, error)
	GetAllDoctorWorkingHours(ctx context.Context, all *entity.GetAll) (*entity.ListDoctorWorkingHours, error)
	UpdateDoctorWorkingHours(ctx context.Context, in *entity.DoctorWorkingHours) (*entity.DoctorWorkingHours, error)
	DeleteDoctorWorkingHours(ctx context.Context, in *entity.GetReqStr) (bool, error)
}
