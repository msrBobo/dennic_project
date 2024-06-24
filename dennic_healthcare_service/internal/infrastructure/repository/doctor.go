package repository

import (
	"Healthcare_Evrone/internal/entity"
	"context"
)

type Doctor interface {
	CreateDoctor(ctx context.Context, doctor *entity.Doctor) (*entity.Doctor, error)
	GetDoctorById(ctx context.Context, get *entity.GetReqStr) (*entity.DoctorAndDoctorHours, error)
	GetAllDoctors(ctx context.Context, all *entity.GetAll) (*entity.ListDoctorsAndHours, error)
	UpdateDoctor(ctx context.Context, update *entity.Doctor) (*entity.Doctor, error)
	DeleteDoctor(ctx context.Context, del *entity.GetReqStr) (bool, error)
	ListDoctorsByDepartmentId(ctx context.Context, in *entity.GetReqStrDep) (*entity.ListDoctors, error)
	ListDoctorBySpecializationId(ctx context.Context, in *entity.GetReqStrSpec) (*entity.ListDoctorsAndHours, error)
}
