package usecase

import (
	"Healthcare_Evrone/internal/entity"
	"Healthcare_Evrone/internal/infrastructure/repository"
	"Healthcare_Evrone/internal/pkg/otlp"
	"go.opentelemetry.io/otel/attribute"
	// "Healthcare_Evrone/internal/pkg/otlp"
	"context"
	"time"
)

const (
	serviceNameDoctorUseCase           = "doctorUseCase"
	serviceNameDoctorUseCaseRepoPrefix = "doctorUseCase"
)

type DoctorUsecase interface {
	CreateDoctor(ctx context.Context, doctor *entity.Doctor) (*entity.Doctor, error)
	GetDoctorById(ctx context.Context, get *entity.GetReqStr) (*entity.DoctorAndDoctorHours, error)
	GetAllDoctors(ctx context.Context, all *entity.GetAll) (*entity.ListDoctorsAndHours, error)
	UpdateDoctor(ctx context.Context, update *entity.Doctor) (*entity.Doctor, error)
	DeleteDoctor(ctx context.Context, del *entity.GetReqStr) (bool, error)
	ListDoctorsByDepartmentId(ctx context.Context, in *entity.GetReqStrDep) (*entity.ListDoctors, error)
	ListDoctorBySpecializationId(ctx context.Context, spec *entity.GetReqStrSpec) (*entity.ListDoctorsAndHours, error)
}

type newsService struct {
	BaseUseCase
	repo       repository.Doctor
	ctxTimeout time.Duration
}

func NewDoctorService(ctxTimeout time.Duration, repo repository.Doctor) newsService {
	return newsService{
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}

func (u newsService) CreateDoctor(ctx context.Context, doctor *entity.Doctor) (*entity.Doctor, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDoctorUseCase, serviceNameDoctorUseCaseRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("CreateDoctor").String(doctor.Id))
	defer span.End()

	return u.repo.CreateDoctor(ctx, doctor)
}

func (u newsService) GetDoctorById(ctx context.Context, get *entity.GetReqStr) (*entity.DoctorAndDoctorHours, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)

	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDoctorUseCase, serviceNameDoctorUseCaseRepoPrefix+"Get")
	span.SetAttributes(attribute.Key(get.Field).String(get.Value))
	defer span.End()

	return u.repo.GetDoctorById(ctx, get)
}

func (u newsService) GetAllDoctors(ctx context.Context, all *entity.GetAll) (*entity.ListDoctorsAndHours, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDoctorUseCase, serviceNameDoctorUseCaseRepoPrefix+"Get all")
	span.SetAttributes(attribute.Key(all.Field).String(all.Value))

	defer span.End()

	return u.repo.GetAllDoctors(ctx, all)
}

func (u newsService) UpdateDoctor(ctx context.Context, update *entity.Doctor) (*entity.Doctor, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDoctorUseCase, serviceNameDoctorUseCaseRepoPrefix+"Update")
	span.SetAttributes(attribute.Key("UpdateDoctor").String(update.Id))
	defer span.End()

	return u.repo.UpdateDoctor(ctx, update)
}

func (u newsService) DeleteDoctor(ctx context.Context, del *entity.GetReqStr) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDoctorUseCase, serviceNameDoctorUseCaseRepoPrefix+"Delete")
	span.SetAttributes(attribute.Key("DeleteDoctor").String(del.Value))
	defer span.End()

	return u.repo.DeleteDoctor(ctx, del)
}

func (u newsService) ListDoctorsByDepartmentId(ctx context.Context, in *entity.GetReqStrDep) (*entity.ListDoctors, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDoctorUseCase, serviceNameDoctorUseCaseRepoPrefix+"Get all")
	span.SetAttributes(attribute.Key("ListDoctorsByDepartmentId").String(in.DepartmentId))

	defer span.End()

	return u.repo.ListDoctorsByDepartmentId(ctx, in)
}

func (u newsService) ListDoctorBySpecializationId(ctx context.Context, in *entity.GetReqStrSpec) (*entity.ListDoctorsAndHours, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDoctorUseCase, serviceNameDoctorUseCaseRepoPrefix+"Get all")
	span.SetAttributes(attribute.Key("ListDoctorBySpecializationId").String(in.SpecializationId))

	defer span.End()

	return u.repo.ListDoctorBySpecializationId(ctx, in)
}
