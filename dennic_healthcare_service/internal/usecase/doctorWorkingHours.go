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
	serviceNameDoctorWorkingHoursUseCase           = "doctor_working_hoursUseCase"
	serviceNameDoctorWorkingHoursUseCaseRepoPrefix = "doctor_working_hoursUseCase"
)

type DoctorWorkingHoursUseCase interface {
	CreateDoctorWorkingHours(ctx context.Context, in *entity.DoctorWorkingHours) (*entity.DoctorWorkingHours, error)
	GetDoctorWorkingHoursById(ctx context.Context, in *entity.GetRequest) (*entity.DoctorWorkingHours, error)
	GetAllDoctorWorkingHours(ctx context.Context, in *entity.GetAll) (*entity.ListDoctorWorkingHours, error)
	UpdateDoctorWorkingHours(ctx context.Context, in *entity.DoctorWorkingHours) (*entity.DoctorWorkingHours, error)
	DeleteDoctorWorkingHours(ctx context.Context, in *entity.GetReqStr) (bool, error)
}

type dwhService struct {
	BaseUseCase
	repo       repository.DoctorWorkingHoursRepository
	ctxTimeout time.Duration
}

func NewDoctorWorkingHoursService(ctxTimeout time.Duration, repo repository.DoctorWorkingHoursRepository) dwhService {
	return dwhService{
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}
func (d dwhService) CreateDoctorWorkingHours(ctx context.Context, in *entity.DoctorWorkingHours) (*entity.DoctorWorkingHours, error) {
	ctx, cancel := context.WithTimeout(ctx, d.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDoctorWorkingHoursUseCase, serviceNameDoctorWorkingHoursUseCaseRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("CreateDoctorWorkingHours").String(string(in.Id)))
	defer span.End()

	return d.repo.CreateDoctorWorkingHours(ctx, in)
}

func (d dwhService) GetDoctorWorkingHoursById(ctx context.Context, in *entity.GetRequest) (*entity.DoctorWorkingHours, error) {
	ctx, cancel := context.WithTimeout(ctx, d.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDoctorWorkingHoursUseCase, serviceNameDoctorWorkingHoursUseCaseRepoPrefix+"Get")
	span.SetAttributes(attribute.Key(in.Field).String(in.Value))

	defer span.End()

	return d.repo.GetDoctorWorkingHoursById(ctx, in)
}

func (d dwhService) GetAllDoctorWorkingHours(ctx context.Context, all *entity.GetAll) (*entity.ListDoctorWorkingHours, error) {
	ctx, cancel := context.WithTimeout(ctx, d.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDoctorWorkingHoursUseCase, serviceNameDoctorWorkingHoursUseCaseRepoPrefix+"Get all")
	span.SetAttributes(attribute.Key(all.Field).String(all.Value))

	defer span.End()

	return d.repo.GetAllDoctorWorkingHours(ctx, all)
}

func (d dwhService) UpdateDoctorWorkingHours(ctx context.Context, in *entity.DoctorWorkingHours) (*entity.DoctorWorkingHours, error) {
	ctx, cancel := context.WithTimeout(ctx, d.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDoctorWorkingHoursUseCase, serviceNameDoctorWorkingHoursUseCaseRepoPrefix+"Update")
	span.SetAttributes(attribute.Key("UpdateDoctorWorkingHours").String(string(in.Id)))
	defer span.End()

	return d.repo.UpdateDoctorWorkingHours(ctx, in)
}

func (d dwhService) DeleteDoctorWorkingHours(ctx context.Context, in *entity.GetReqStr) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, d.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDoctorWorkingHoursUseCase, serviceNameDoctorWorkingHoursUseCaseRepoPrefix+"Delete")
	span.SetAttributes(attribute.Key("DeleteDoctorWorkingHours").String(in.Value))
	defer span.End()

	return d.repo.DeleteDoctorWorkingHours(ctx, in)
}
