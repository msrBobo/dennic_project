package usecase

import (
	"Healthcare_Evrone/internal/entity"
	"Healthcare_Evrone/internal/infrastructure/repository"
	"Healthcare_Evrone/internal/pkg/otlp"
	"context"
	"go.opentelemetry.io/otel/attribute"
	"time"
)

const (
	serviceNameDoctorServiceUseCase           = "doctor_serviceUseCase"
	serviceNameDoctorServiceUseCaseRepoPrefix = "doctor_serviceUseCase"
)

type DoctorServices interface {
	CreateDoctorServices(ctx context.Context, in *entity.DoctorServices) (*entity.DoctorServices, error)
	GetDoctorServiceByID(ctx context.Context, in *entity.GetReqStr) (*entity.DoctorServices, error)
	GetAllDoctorServices(ctx context.Context, all *entity.GetAll) (*entity.ListDoctorServices, error)
	UpdateDoctorServices(ctx context.Context, in *entity.DoctorServices) (*entity.DoctorServices, error)
	DeleteDoctorService(ctx context.Context, in *entity.GetReqStr) (bool, error)
}

type newsDoctorSService struct {
	BaseUseCase
	repo       repository.DoctorServices
	ctxTimeout time.Duration
}

func NewDoctorServices(ctxTimeout time.Duration, ds repository.DoctorServices) newsDoctorSService {
	return newsDoctorSService{
		ctxTimeout: ctxTimeout,
		repo:       ds,
	}
}

func (u newsDoctorSService) CreateDoctorServices(ctx context.Context, in *entity.DoctorServices) (*entity.DoctorServices, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDoctorServiceUseCase, serviceNameDoctorServiceUseCaseRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("CreateDoctorServices").String(in.Id))
	defer span.End()

	return u.repo.CreateDoctorServices(ctx, in)
}

func (u newsDoctorSService) GetDoctorServiceByID(ctx context.Context, in *entity.GetReqStr) (*entity.DoctorServices, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDoctorServiceUseCase, serviceNameDoctorServiceUseCaseRepoPrefix+"Get")
	span.SetAttributes(attribute.Key(in.Field).String(in.Value))
	defer span.End()

	return u.repo.GetDoctorServiceByID(ctx, in)
}

func (u newsDoctorSService) GetAllDoctorServices(ctx context.Context, all *entity.GetAll) (*entity.ListDoctorServices, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)

	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDoctorServiceUseCase, serviceNameDoctorServiceUseCaseRepoPrefix+"Get all")
	span.SetAttributes(attribute.Key(all.Field).String(all.Value))

	defer span.End()

	return u.repo.GetAllDoctorServices(ctx, all)
}

func (u newsDoctorSService) UpdateDoctorServices(ctx context.Context, in *entity.DoctorServices) (*entity.DoctorServices, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDoctorServiceUseCase, serviceNameDoctorServiceUseCaseRepoPrefix+"Update")
	span.SetAttributes(attribute.Key("doctorService").String(in.Id))
	defer span.End()

	return u.repo.UpdateDoctorServices(ctx, in)
}

func (u newsDoctorSService) DeleteDoctorService(ctx context.Context, in *entity.GetReqStr) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDoctorServiceUseCase, serviceNameDoctorServiceUseCaseRepoPrefix+"Delete")
	span.SetAttributes(attribute.Key("doctorService").String(in.Value))
	defer span.End()

	return u.repo.DeleteDoctorService(ctx, in)
}
