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
	serviceNameDepartmentUseCase           = "departmentUseCase"
	serviceNameDepartmentUseCaseRepoPrefix = "departmentUseCase"
)

type DepartmentsUsecase interface {
	CreateDepartment(ctx context.Context, dep *entity.Department) (*entity.Department, error)
	GetDepartmentById(ctx context.Context, get *entity.GetReqStr) (*entity.Department, error)
	GetAllDepartments(ctx context.Context, all *entity.GetAll) (*entity.ListDepartments, error)
	UpdateDepartment(ctx context.Context, update *entity.Department) (*entity.Department, error)
	DeleteDepartment(ctx context.Context, del *entity.GetReqStr) (bool, error)
}

type newsDepService struct {
	BaseUseCase
	repo       repository.DepartmentRepository
	ctxTimeout time.Duration
}

func NewDepartmentService(ctxTimeout time.Duration, repo repository.DepartmentRepository) newsDepService {
	return newsDepService{
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}

func (u newsDepService) CreateDepartment(ctx context.Context, dep *entity.Department) (*entity.Department, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDepartmentUseCase, serviceNameDepartmentUseCaseRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("CreateDepartment").String(dep.Id))
	defer span.End()

	res, err := u.repo.CreateDepartment(ctx, dep)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (u newsDepService) GetDepartmentById(ctx context.Context, get *entity.GetReqStr) (*entity.Department, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDepartmentUseCase, serviceNameDepartmentUseCaseRepoPrefix+"Get")
	span.SetAttributes(attribute.Key(get.Field).String(get.Value))

	defer span.End()

	return u.repo.GetDepartmentById(ctx, get)
}

func (u newsDepService) GetAllDepartments(ctx context.Context, all *entity.GetAll) (*entity.ListDepartments, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDepartmentUseCase, serviceNameDepartmentUseCaseRepoPrefix+"Get all")
	span.SetAttributes(attribute.Key(all.Field).String(all.Value))

	defer span.End()

	return u.repo.GetAllDepartments(ctx, all)
}

func (u newsDepService) UpdateDepartment(ctx context.Context, update *entity.Department) (*entity.Department, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDepartmentUseCase, serviceNameDepartmentUseCaseRepoPrefix+"Update")
	span.SetAttributes(attribute.Key("UpdateDepartment").String(update.Id))

	defer span.End()

	return u.repo.UpdateDepartment(ctx, update)
}

func (u newsDepService) DeleteDepartment(ctx context.Context, del *entity.GetReqStr) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDepartmentUseCase, serviceNameDepartmentUseCaseRepoPrefix+"Delete")
	span.SetAttributes(attribute.Key("DeleteDepartment").String(del.Value))
	defer span.End()

	return u.repo.DeleteDepartment(ctx, del)
}
