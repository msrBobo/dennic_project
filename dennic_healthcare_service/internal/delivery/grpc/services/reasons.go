package services

import (
	pb "Healthcare_Evrone/genproto/healthcare-service"
	"Healthcare_Evrone/internal/entity"
	"Healthcare_Evrone/internal/pkg/minio"
	"Healthcare_Evrone/internal/pkg/otlp"
	"Healthcare_Evrone/internal/usecase"
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type reasonsServiceRPC struct {
	logger  *zap.Logger
	reasons usecase.ReasonUseCase
}

const (
	serviceNameReasonsDelivery           = "reasonsDelivery"
	serviceNameReasonsDeliveryRepoPrefix = "reasonsDelivery"
)

func ReasonsServiceRPC(logger *zap.Logger, dwhUsecase usecase.ReasonUseCase) pb.ReasonsServiceServer {
	return &reasonsServiceRPC{
		logger,
		dwhUsecase,
	}
}

func (r reasonsServiceRPC) CreateReasons(ctx context.Context, reasons *pb.Reasons) (*pb.Reasons, error) {
	ctx, span := otlp.Start(ctx, serviceNameReasonsDelivery, serviceNameReasonsDeliveryRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("CreateReasons").String(reasons.Id))
	defer span.End()

	imageUrl := minio.RemoveImageUrl(reasons.ImageUrl)

	reas := &entity.Reasons{
		Id:               reasons.Id,
		Name:             reasons.Name,
		SpecializationId: reasons.SpecializationId,
		ImageUrl:         imageUrl,
	}
	resp, err := r.reasons.CreateReasons(ctx, reas)
	if err != nil {
		return nil, err
	}
	respImageUrl := minio.AddImageUrl(resp.ImageUrl, cfg.MinioService.Bucket.Reasons)

	return &pb.Reasons{
		Id:               resp.Id,
		Name:             resp.Name,
		SpecializationId: resp.SpecializationId,
		ImageUrl:         respImageUrl,
		CreatedAt:        resp.CreatedAt.String(),
		UpdatedAt:        resp.UpdatedAt.String(),
		DeletedAt:        resp.DeletedAt.String(),
	}, nil
}

func (r reasonsServiceRPC) GetReasonsById(ctx context.Context, reasons *pb.GetReqStrReasons) (*pb.Reasons, error) {
	ctx, span := otlp.Start(ctx, serviceNameReasonsDelivery, serviceNameReasonsDeliveryRepoPrefix+"Get")
	span.SetAttributes(attribute.Key("GetReasonsById").String(reasons.Value))
	defer span.End()

	resp, err := r.reasons.GetReasonsById(ctx, &entity.GetReqStr{
		Field:    reasons.Field,
		Value:    reasons.Value,
		IsActive: reasons.IsActive,
	})
	if err != nil {
		return nil, err
	}
	respImageUrl := minio.AddImageUrl(resp.ImageUrl, cfg.MinioService.Bucket.Reasons)

	return &pb.Reasons{
		Id:               resp.Id,
		Name:             resp.Name,
		SpecializationId: resp.SpecializationId,
		ImageUrl:         respImageUrl,
		CreatedAt:        resp.CreatedAt.String(),
		UpdatedAt:        resp.UpdatedAt.String(),
		DeletedAt:        resp.DeletedAt.String(),
	}, nil
}

func (r reasonsServiceRPC) GetAllReasons(ctx context.Context, all *pb.GetAllReas) (*pb.ListReasons, error) {
	ctx, span := otlp.Start(ctx, serviceNameDoctorServiceDelivery, serviceNameDoctorServiceDeliveryRepoPrefix+"Get all")
	span.SetAttributes(attribute.Key("GetAllDoctorServices").String(all.Field))
	defer span.End()
	reasonss, err := r.reasons.GetAllReasons(ctx, &entity.GetAllReas{
		Page:     all.Page,
		Limit:    all.Limit,
		IsActive: all.IsActive,
		Field:    all.Field,
		Value:    all.Value,
		OrderBy:  all.OrderBy,
	})
	if err != nil {
		return nil, err
	}

	var listReasonss pb.ListReasons
	for _, reasons := range reasonss.Reasons {
		respImageUrl := minio.AddImageUrl(reasons.ImageUrl, cfg.MinioService.Bucket.Reasons)

		listReasonss.Reasons = append(listReasonss.Reasons, &pb.Reasons{
			Id:               reasons.Id,
			Name:             reasons.Name,
			SpecializationId: reasons.SpecializationId,
			ImageUrl:         respImageUrl,
			CreatedAt:        reasons.CreatedAt.String(),
			UpdatedAt:        reasons.UpdatedAt.String(),
			DeletedAt:        reasons.DeletedAt.String(),
		})
	}
	listReasonss.Count += reasonss.Count
	return &listReasonss, nil
}

func (r reasonsServiceRPC) UpdateReasons(ctx context.Context, reasons *pb.Reasons) (*pb.Reasons, error) {
	ctx, span := otlp.Start(ctx, serviceNameReasonsDelivery, serviceNameReasonsDeliveryRepoPrefix+"Update")
	span.SetAttributes(attribute.Key("UpdateReasons").String(reasons.Id))
	defer span.End()

	reqImageUrl := minio.RemoveImageUrl(reasons.ImageUrl)

	reas := &entity.Reasons{
		Id:               reasons.Id,
		Name:             reasons.Name,
		SpecializationId: reasons.SpecializationId,
		ImageUrl:         reqImageUrl,
	}
	resp, err := r.reasons.UpdateReasons(ctx, reas)
	if err != nil {
		return nil, err
	}

	respImageUrl := minio.AddImageUrl(resp.ImageUrl, cfg.MinioService.Bucket.Reasons)

	return &pb.Reasons{
		Id:               resp.Id,
		Name:             resp.Name,
		SpecializationId: resp.SpecializationId,
		ImageUrl:         respImageUrl,
		CreatedAt:        resp.CreatedAt.String(),
		UpdatedAt:        resp.UpdatedAt.String(),
		DeletedAt:        resp.DeletedAt.String(),
	}, nil
}

func (r reasonsServiceRPC) DeleteReasons(ctx context.Context, reasons *pb.GetReqStrReasons) (*pb.StatusReasons, error) {
	ctx, span := otlp.Start(ctx, serviceNameReasonsDelivery, serviceNameReasonsDeliveryRepoPrefix+"Delete")
	span.SetAttributes(attribute.Key("UpdateReasons").String(reasons.Value))
	defer span.End()

	resp, err := r.reasons.DeleteReasons(ctx, &entity.GetReqStr{
		Field:    reasons.Field,
		Value:    reasons.Value,
		IsActive: reasons.IsActive,
	})
	if err != nil {
		return nil, err
	}
	return &pb.StatusReasons{Status: resp.Status}, nil
}
