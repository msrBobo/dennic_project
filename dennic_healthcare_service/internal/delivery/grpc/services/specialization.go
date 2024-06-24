package services

import (
	pb "Healthcare_Evrone/genproto/healthcare-service"
	"Healthcare_Evrone/internal/entity"
	"Healthcare_Evrone/internal/pkg/minio"
	"Healthcare_Evrone/internal/pkg/otlp"
	"Healthcare_Evrone/internal/usecase"
	"Healthcare_Evrone/internal/usecase/event"
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type specializationRPC struct {
	logger         *zap.Logger
	specialization usecase.SpecializationUsecase
	brokerProducer event.BrokerProducer
}

const (
	serviceNameSpecializationDelivery           = "SpecializationDelivery"
	serviceNameSpecializationDeliveryRepoPrefix = "SpecializationDelivery"
)

func SpecializationRPC(logger *zap.Logger,
	specializationUsecase usecase.SpecializationUsecase,
	brokerProducer event.BrokerProducer) pb.SpecializationServiceServer {
	return &specializationRPC{
		logger,
		specializationUsecase,
		brokerProducer,
	}
}

func (r specializationRPC) CreateSpecialization(ctx context.Context, specializations *pb.Specializations) (*pb.Specializations, error) {
	ctx, span := otlp.Start(ctx, serviceNameSpecializationDelivery, serviceNameSpecializationDeliveryRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("CreateSpecialization").String(specializations.Id))
	defer span.End()

	reqImageUrl := minio.RemoveImageUrl(specializations.ImageUrl)

	req := entity.Specialization{
		ID:           specializations.Id,
		Order:        specializations.Order,
		Name:         specializations.Name,
		Description:  specializations.Description,
		DepartmentId: specializations.DepartmentId,
		ImageUrl:     reqImageUrl,
	}
	resp, err := r.specialization.CreateSpecialization(ctx, &req)
	if err != nil {
		return nil, err
	}
	respImageUrl := minio.AddImageUrl(resp.ImageUrl, cfg.MinioService.Bucket.Specialization)
	return &pb.Specializations{
		Id:           resp.ID,
		Order:        resp.Order,
		Name:         resp.Name,
		Description:  resp.Description,
		DepartmentId: resp.DepartmentId,
		ImageUrl:     respImageUrl,
		CreatedAt:    resp.CreatedAt.String(),
		UpdatedAt:    resp.UpdatedAt.String(),
		DeletedAt:    resp.DeletedAt.String(),
	}, nil
}

func (r specializationRPC) GetSpecializationById(ctx context.Context, str *pb.GetReqStrSpecialization) (*pb.Specializations, error) {

	ctx, span := otlp.Start(ctx, serviceNameSpecializationDelivery, serviceNameSpecializationDeliveryRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("GetSpecializationById").String(str.Value))
	defer span.End()
	spec, err := r.specialization.GetSpecializationById(ctx, &entity.GetReqStr{
		Field:    str.Field,
		Value:    str.Value,
		IsActive: str.IsActive,
	})
	if err != nil {
		return nil, err
	}
	respImageUrl := minio.AddImageUrl(spec.ImageUrl, cfg.MinioService.Bucket.Specialization)
	return &pb.Specializations{
		Id:           spec.ID,
		Order:        spec.Order,
		Name:         spec.Name,
		Description:  spec.Description,
		DepartmentId: spec.DepartmentId,
		ImageUrl:     respImageUrl,
		CreatedAt:    spec.CreatedAt.String(),
		UpdatedAt:    spec.UpdatedAt.String(),
		DeletedAt:    spec.DeletedAt.String(),
	}, nil
}

func (r specializationRPC) GetAllSpecializations(ctx context.Context, all *pb.GetAllSpecialization) (*pb.ListSpecializations, error) {
	ctx, span := otlp.Start(ctx, serviceNameSpecializationDelivery, serviceNameSpecializationDeliveryRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("GetAllSpecializations").String(all.Value))
	defer span.End()
	specializations, err := r.specialization.GetAllSpecializations(ctx, &entity.GetAllSpecializations{
		Page:         int64(all.Page),
		Limit:        int64(all.Limit),
		Field:        all.Field,
		Value:        all.Value,
		OrderBy:      all.OrderBy,
		IsActive:     all.IsActive,
		DepartmentId: all.DepartmentId,
	})
	if err != nil {
		return nil, err
	}
	var listSpec pb.ListSpecializations
	for _, s := range specializations.Specializations {
		respImageUrl := minio.AddImageUrl(s.ImageUrl, cfg.MinioService.Bucket.Specialization)
		listSpec.Specializations = append(listSpec.Specializations, &pb.Specializations{
			Id:           s.ID,
			Order:        s.Order,
			Name:         s.Name,
			Description:  s.Description,
			DepartmentId: s.DepartmentId,
			ImageUrl:     respImageUrl,
			CreatedAt:    s.CreatedAt.String(),
			UpdatedAt:    s.UpdatedAt.String(),
			DeletedAt:    s.DeletedAt.String(),
		})
	}
	listSpec.Count = specializations.Count
	return &listSpec, nil
}

func (r specializationRPC) UpdateSpecialization(ctx context.Context, in *pb.Specializations) (*pb.Specializations, error) {
	ctx, span := otlp.Start(ctx, serviceNameSpecializationDelivery, serviceNameSpecializationDeliveryRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("UpdateSpecialization").String(in.Id))
	defer span.End()
	reqImageUrl := minio.RemoveImageUrl(in.ImageUrl)
	resp, err := r.specialization.UpdateSpecialization(ctx, &entity.Specialization{
		ID:           in.Id,
		Order:        in.Order,
		Name:         in.Name,
		Description:  in.Description,
		DepartmentId: in.DepartmentId,
		ImageUrl:     reqImageUrl,
	})
	if err != nil {
		return nil, err
	}
	respImageUrl := minio.AddImageUrl(resp.ImageUrl, cfg.MinioService.Bucket.Specialization)
	return &pb.Specializations{
		Id:           resp.ID,
		Order:        resp.Order,
		Name:         resp.Name,
		Description:  resp.Description,
		DepartmentId: resp.DepartmentId,
		ImageUrl:     respImageUrl,
		CreatedAt:    resp.CreatedAt.String(),
		UpdatedAt:    resp.UpdatedAt.String(),
		DeletedAt:    resp.DeletedAt.String(),
	}, nil
}

func (r specializationRPC) DeleteSpecialization(ctx context.Context, in *pb.GetReqStrSpecialization) (*pb.StatusSpecialization, error) {
	ctx, span := otlp.Start(ctx, serviceNameSpecializationDelivery, serviceNameSpecializationDeliveryRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("DeleteSpecialization").String(in.Value))
	defer span.End()
	status, err := r.specialization.DeleteSpecialization(ctx, &entity.GetReqStr{Value: in.Value, Field: in.Field, IsActive: in.IsActive})
	if err != nil {
		return nil, err
	}
	return &pb.StatusSpecialization{Status: status}, nil
}
