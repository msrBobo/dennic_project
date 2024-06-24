package usecase

import (
	"Healthcare_Evrone/internal/entity"
	"Healthcare_Evrone/internal/infrastructure/repository"
	"Healthcare_Evrone/internal/pkg/otlp"
	"context"
	"go.opentelemetry.io/otel/attribute"
	"time"
)

type ReasonUseCase interface {
	CreateReasons(ctx context.Context, in *entity.Reasons) (*entity.Reasons, error)
	GetReasonsById(ctx context.Context, in *entity.GetReqStr) (*entity.Reasons, error)
	GetAllReasons(context.Context, *entity.GetAllReas) (*entity.ListReasons, error)
	UpdateReasons(context.Context, *entity.Reasons) (*entity.Reasons, error)
	DeleteReasons(context.Context, *entity.GetReqStr) (*entity.StatusReasons, error)
}

type newsReasons struct {
	BaseUseCase
	repo       repository.ReasonRepository
	ctxTimeout time.Duration
}

const (
	serviceNameReasonsUseCase           = "reasonsUseCase"
	serviceNameReasonsUseCaseRepoPrefix = "reasonsUseCase"
)

func NewReasons(ctxTimeout time.Duration, reasons repository.ReasonRepository) newsReasons {
	return newsReasons{
		ctxTimeout: ctxTimeout,
		repo:       reasons,
	}
}

func (n newsReasons) CreateReasons(ctx context.Context, in *entity.Reasons) (*entity.Reasons, error) {
	ctx, cancel := context.WithTimeout(ctx, n.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameReasonsUseCase, serviceNameReasonsUseCaseRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("CreateReasons").String(in.Id))
	defer span.End()

	return n.repo.CreateReasons(ctx, in)
}

func (n newsReasons) GetReasonsById(ctx context.Context, in *entity.GetReqStr) (*entity.Reasons, error) {
	ctx, cancel := context.WithTimeout(ctx, n.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameReasonsUseCase, serviceNameReasonsUseCaseRepoPrefix+"Get")
	span.SetAttributes(attribute.Key(in.Field).String(in.Value))
	defer span.End()

	return n.repo.GetReasonsById(ctx, in)
}

func (n newsReasons) GetAllReasons(ctx context.Context, reas *entity.GetAllReas) (*entity.ListReasons, error) {
	ctx, cancel := context.WithTimeout(ctx, n.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameReasonsUseCase, serviceNameReasonsUseCaseRepoPrefix+"Get all")
	span.SetAttributes(attribute.Key(reas.Field).String(reas.Value))
	defer span.End()

	return n.repo.GetAllReasons(ctx, reas)
}

func (n newsReasons) UpdateReasons(ctx context.Context, reasons *entity.Reasons) (*entity.Reasons, error) {
	ctx, cancel := context.WithTimeout(ctx, n.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameReasonsUseCase, serviceNameReasonsUseCaseRepoPrefix+"update")
	span.SetAttributes(attribute.Key("UpdateReasons").String(reasons.Id))
	defer span.End()

	return n.repo.UpdateReasons(ctx, reasons)
}

func (n newsReasons) DeleteReasons(ctx context.Context, reasons *entity.GetReqStr) (*entity.StatusReasons, error) {
	ctx, cancel := context.WithTimeout(ctx, n.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameReasonsUseCase, serviceNameReasonsUseCaseRepoPrefix+"delete")
	span.SetAttributes(attribute.Key("DeleteReasons").String(reasons.Value))
	defer span.End()

	return n.repo.DeleteReasons(ctx, reasons)
}
