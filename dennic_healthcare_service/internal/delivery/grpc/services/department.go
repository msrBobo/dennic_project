package services

import (
	pb "Healthcare_Evrone/genproto/healthcare-service"
	"Healthcare_Evrone/internal/entity"
	"Healthcare_Evrone/internal/pkg/config"
	"Healthcare_Evrone/internal/pkg/minio"
	"Healthcare_Evrone/internal/pkg/otlp"
	"Healthcare_Evrone/internal/usecase"
	"context"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"time"
)

type departmentRPC struct {
	logger     *zap.Logger
	department usecase.DepartmentsUsecase
}

const (
	serviceNameDepartmentDelivery           = "DepartmentDelivery"
	serviceNameDepartmentDeliveryRepoPrefix = "DepartmentDelivery"
)

var cfg = config.New()

func DepartmentRPC(logget *zap.Logger, departmentUsecase usecase.DepartmentsUsecase) pb.DepartmentServiceServer {
	return &departmentRPC{
		logget,
		departmentUsecase,
	}

}

func (r departmentRPC) CreateDepartment(ctx context.Context, dep *pb.Department) (*pb.Department, error) {

	ctx, span := otlp.Start(ctx, serviceNameDepartmentDelivery, serviceNameDepartmentDeliveryRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("CreateDepartment").String(dep.Id))
	defer span.End()

	reqImageUrl := minio.RemoveImageUrl(dep.ImageUrl)

	req := entity.Department{
		Id:               dep.Id,
		Order:            dep.Order,
		Name:             dep.Name,
		Description:      dep.Description,
		ImageUrl:         reqImageUrl,
		FloorNumber:      dep.FloorNumber,
		ShortDescription: dep.ShortDescription,
	}

	resp, err := r.department.CreateDepartment(ctx, &req)

	if err != nil {
		return nil, err
	}
	respImageUrl := minio.AddImageUrl(resp.ImageUrl, cfg.MinioService.Bucket.Department)

	return &pb.Department{
		Id:               resp.Id,
		Order:            resp.Order,
		Name:             resp.Name,
		Description:      resp.Description,
		ImageUrl:         respImageUrl,
		FloorNumber:      resp.FloorNumber,
		ShortDescription: resp.ShortDescription,
		CreatedAt:        resp.CreatedAt.String(),
		UpdatedAt:        resp.UpdatedAt.String(),
		DeletedAt:        resp.DeletedAt.String(),
	}, nil

}

func (r departmentRPC) GetDepartmentById(ctx context.Context, get *pb.GetReqStrDepartment) (*pb.Department, error) {
	ctx, span := otlp.Start(ctx, serviceNameDepartmentDelivery, serviceNameDepartmentDeliveryRepoPrefix+"Get")
	span.SetAttributes(attribute.Key("GetDepartmentById").String(get.Value))
	defer span.End()
	resp, err := r.department.GetDepartmentById(ctx, &entity.GetReqStr{Field: get.Field, Value: get.Value, IsActive: get.IsActive})
	if err != nil {
		return nil, err
	}
	respImageUrl := minio.AddImageUrl(resp.ImageUrl, cfg.MinioService.Bucket.Department)
	return &pb.Department{
		Id:               resp.Id,
		Order:            resp.Order,
		Name:             resp.Name,
		Description:      resp.Description,
		ImageUrl:         respImageUrl,
		FloorNumber:      resp.FloorNumber,
		ShortDescription: resp.ShortDescription,
		CreatedAt:        resp.CreatedAt.String(),
		UpdatedAt:        resp.UpdatedAt.String(),
		DeletedAt:        resp.DeletedAt.String(),
	}, nil
}

func (r departmentRPC) GetAllDepartments(ctx context.Context, get *pb.GetAllDepartment) (*pb.ListDepartments, error) {

	ctx, span := otlp.Start(ctx, serviceNameDepartmentDelivery, serviceNameDepartmentDeliveryRepoPrefix+"Get all")
	span.SetAttributes(attribute.Key("GetAllDepartments").String(get.Value))
	defer span.End()

	resp, err := r.department.GetAllDepartments(ctx, &entity.GetAll{
		Page:     get.Page,
		Limit:    get.Limit,
		Field:    get.Field,
		Value:    get.Value,
		OrderBy:  get.OrderBy,
		IsActive: get.IsActive,
	})
	if err != nil {
		return nil, err
	}

	var departments pb.ListDepartments
	for _, dep := range resp.Departments {
		respImageUrl := minio.AddImageUrl(dep.ImageUrl, cfg.MinioService.Bucket.Department)
		departments.Departments = append(departments.Departments, &pb.Department{
			Id:               dep.Id,
			Order:            dep.Order,
			Name:             dep.Name,
			Description:      dep.Description,
			ImageUrl:         respImageUrl,
			FloorNumber:      dep.FloorNumber,
			ShortDescription: dep.ShortDescription,
			CreatedAt:        dep.CreatedAt.String(),
			UpdatedAt:        dep.UpdatedAt.String(),
			DeletedAt:        dep.DeletedAt.String(),
		})
	}
	departments.Count = resp.Count
	return &departments, nil
}

func (r departmentRPC) UpdateDepartment(ctx context.Context, update *pb.Department) (*pb.Department, error) {

	ctx, span := otlp.Start(ctx, serviceNameDepartmentDelivery, serviceNameDepartmentDeliveryRepoPrefix+"Update")
	span.SetAttributes(attribute.Key("UpdateDepartment").String(update.Id))
	defer span.End()
	reqImageUrl := minio.RemoveImageUrl(update.ImageUrl)
	req := entity.Department{
		Id:               update.Id,
		Order:            update.Order,
		Name:             update.Name,
		Description:      update.Description,
		ImageUrl:         reqImageUrl,
		FloorNumber:      update.FloorNumber,
		ShortDescription: update.ShortDescription,
		CreatedAt:        time.Time{},
		UpdatedAt:        time.Time{},
		DeletedAt:        time.Time{},
	}
	resp, err := r.department.UpdateDepartment(ctx, &req)
	if err != nil {
		return nil, err
	}
	respImageUrl := minio.AddImageUrl(resp.ImageUrl, cfg.MinioService.Bucket.Department)
	return &pb.Department{
		Id:               resp.Id,
		Order:            resp.Order,
		Name:             resp.Name,
		Description:      resp.Description,
		ImageUrl:         respImageUrl,
		FloorNumber:      resp.FloorNumber,
		ShortDescription: resp.ShortDescription,
		CreatedAt:        resp.CreatedAt.String(),
		UpdatedAt:        resp.UpdatedAt.String(),
		DeletedAt:        resp.DeletedAt.String(),
	}, nil
}

func (r departmentRPC) DeleteDepartment(ctx context.Context, del *pb.GetReqStrDepartment) (*pb.StatusDepartment, error) {
	ctx, span := otlp.Start(ctx, serviceNameDepartmentDelivery, serviceNameDepartmentDeliveryRepoPrefix+"Delete")
	span.SetAttributes(attribute.Key("DeleteDepartment").String(del.Value))
	defer span.End()
	status, err := r.department.DeleteDepartment(ctx, &entity.GetReqStr{Field: del.Field, Value: del.Value, IsActive: del.IsActive})
	if err != nil {
		r.logger.Error("deleted department error", zap.Error(err))
		return nil, err
	}
	fmt.Println(status)
	return &pb.StatusDepartment{Status: status}, nil
}
