package services

import (
	pb "Healthcare_Evrone/genproto/healthcare-service"
	"Healthcare_Evrone/internal/entity"
	"Healthcare_Evrone/internal/pkg/otlp"
	"Healthcare_Evrone/internal/usecase"
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"time"
)

const (
	serviceNameDoctorServiceDelivery           = "doctorServiceDelivery"
	serviceNameDoctorServiceDeliveryRepoPrefix = "doctorServiceDelivery"
)

type doctorServiceRPC struct {
	logger         *zap.Logger
	doctorServices usecase.DoctorServices
}

func DoctorsServiceRPC(logger *zap.Logger, doctorServicesUsecase usecase.DoctorServices) pb.DoctorsServiceServer {
	return &doctorServiceRPC{
		logger,
		doctorServicesUsecase,
	}
}

func (r doctorServiceRPC) CreateDoctorServices(ctx context.Context, in *pb.DoctorServices) (*pb.DoctorServices, error) {

	ctx, span := otlp.Start(ctx, serviceNameDoctorServiceDelivery, serviceNameDoctorServiceDeliveryRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("CreateDoctorServices").String(in.Id))
	defer span.End()

	tm, err := time.Parse("15:04", in.Duration)
	if err != nil {
		return nil, err
	}
	req := entity.DoctorServices{
		Id:               in.Id,
		Order:            in.DoctorServiceOrder,
		DoctorId:         in.DoctorId,
		SpecializationId: in.SpecializationId,
		OnlinePrice:      in.OnlinePrice,
		OfflinePrice:     in.OfflinePrice,
		Name:             in.Name,
		Duration:         tm,
	}
	resp, err := r.doctorServices.CreateDoctorServices(ctx, &req)
	if err != nil {
		return nil, err
	}

	return &pb.DoctorServices{
		Id:                 resp.Id,
		DoctorServiceOrder: resp.Order,
		DoctorId:           resp.DoctorId,
		SpecializationId:   resp.SpecializationId,
		OnlinePrice:        resp.OnlinePrice,
		OfflinePrice:       resp.OfflinePrice,
		Name:               resp.Name,
		Duration:           resp.Duration.Format("15:04"),
		CreatedAt:          resp.CreatedAt.String(),
		UpdatedAt:          resp.UpdatedAt.String(),
		DeletedAt:          resp.DeletedAt.String(),
	}, nil
}

func (r doctorServiceRPC) GetDoctorServiceByID(ctx context.Context, str *pb.GetReqStr) (*pb.DoctorServices, error) {
	ctx, span := otlp.Start(ctx, serviceNameDoctorServiceDelivery, serviceNameDoctorServiceDeliveryRepoPrefix+"Get")
	span.SetAttributes(attribute.Key("GetDoctorServiceByID").String(str.Value))
	defer span.End()
	ds, err := r.doctorServices.GetDoctorServiceByID(ctx, &entity.GetReqStr{
		Field:    str.Field,
		Value:    str.Value,
		IsActive: str.IsActive,
	})
	if err != nil {
		return nil, err
	}
	return &pb.DoctorServices{
		Id:                 ds.Id,
		DoctorServiceOrder: ds.Order,
		DoctorId:           ds.DoctorId,
		SpecializationId:   ds.SpecializationId,
		OnlinePrice:        ds.OnlinePrice,
		OfflinePrice:       ds.OfflinePrice,
		Name:               ds.Name,
		Duration:           ds.Duration.Format("15:04"),
		CreatedAt:          ds.CreatedAt.String(),
		UpdatedAt:          ds.UpdatedAt.String(),
		DeletedAt:          ds.DeletedAt.String(),
	}, nil
}

func (r doctorServiceRPC) GetAllDoctorServices(ctx context.Context, all *pb.GetAllDoctorServiceS) (*pb.ListDoctorServices, error) {
	ctx, span := otlp.Start(ctx, serviceNameDoctorServiceDelivery, serviceNameDoctorServiceDeliveryRepoPrefix+"Get all")
	span.SetAttributes(attribute.Key("GetAllDoctorServices").String(all.Value))
	defer span.End()
	doctorService, err := r.doctorServices.GetAllDoctorServices(ctx, &entity.GetAll{
		Page:     all.Page,
		Limit:    all.Limit,
		Field:    all.Field,
		Value:    all.Value,
		OrderBy:  all.OrderBy,
		IsActive: all.IsActive,
	})
	if err != nil {
		return nil, err
	}
	var listDoctorServices pb.ListDoctorServices
	for _, d := range doctorService.DoctorServices {
		listDoctorServices.DoctorServices = append(listDoctorServices.DoctorServices, &pb.DoctorServices{
			Id:                 d.Id,
			DoctorServiceOrder: d.Order,
			DoctorId:           d.DoctorId,
			SpecializationId:   d.SpecializationId,
			OnlinePrice:        d.OnlinePrice,
			OfflinePrice:       d.OfflinePrice,
			Name:               d.Name,
			Duration:           d.Duration.Format("15:04"),
			CreatedAt:          d.CreatedAt.String(),
			UpdatedAt:          d.UpdatedAt.String(),
			DeletedAt:          d.DeletedAt.String(),
		})
	}
	listDoctorServices.Count = int32(doctorService.Count)
	return &listDoctorServices, nil
}

func (r doctorServiceRPC) UpdateDoctorServices(ctx context.Context, services *pb.DoctorServices) (*pb.DoctorServices, error) {
	ctx, span := otlp.Start(ctx, serviceNameDoctorServiceDelivery, serviceNameDoctorServiceDeliveryRepoPrefix+"Update")
	span.SetAttributes(attribute.Key("UpdateDoctorServices").String(services.Id))
	defer span.End()

	tm, err := time.Parse("15:04", services.Duration)
	if err != nil {
		return nil, err
	}

	resp, err := r.doctorServices.UpdateDoctorServices(ctx, &entity.DoctorServices{
		Id:               services.Id,
		Order:            services.DoctorServiceOrder,
		DoctorId:         services.DoctorId,
		SpecializationId: services.SpecializationId,
		OnlinePrice:      services.OnlinePrice,
		OfflinePrice:     services.OfflinePrice,
		Name:             services.Name,
		Duration:         tm,
	})
	if err != nil {
		return nil, err
	}

	return &pb.DoctorServices{
		Id:                 resp.Id,
		DoctorServiceOrder: resp.Order,
		DoctorId:           resp.DoctorId,
		SpecializationId:   resp.SpecializationId,
		OnlinePrice:        resp.OnlinePrice,
		OfflinePrice:       resp.OfflinePrice,
		Name:               resp.Name,
		Duration:           resp.Duration.Format("15:04:05"),
		CreatedAt:          resp.CreatedAt.String(),
		UpdatedAt:          resp.UpdatedAt.String(),
		DeletedAt:          resp.DeletedAt.String(),
	}, nil
}

func (r doctorServiceRPC) DeleteDoctorService(ctx context.Context, str *pb.GetReqStr) (*pb.Status, error) {
	ctx, span := otlp.Start(ctx, serviceNameDoctorServiceDelivery, serviceNameDoctorServiceDeliveryRepoPrefix+"Delete")
	span.SetAttributes(attribute.Key("DeleteDoctorService").String(str.Value))
	defer span.End()
	status, err := r.doctorServices.DeleteDoctorService(ctx, &entity.GetReqStr{Field: str.Field, Value: str.Value, IsActive: str.IsActive})
	if err != nil {
		r.logger.Error("deleted department error", zap.Error(err))
		return nil, err
	}
	return &pb.Status{Status: status}, nil
}
