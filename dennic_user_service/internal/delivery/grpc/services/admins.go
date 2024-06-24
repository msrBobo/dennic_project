package services

import (
	"context"
	pb "dennic_user_service/genproto/user_service"
	"dennic_user_service/internal/entity"
	"dennic_user_service/internal/pkg/config"
	"dennic_user_service/internal/pkg/minio"
	"dennic_user_service/internal/pkg/otlp"
	"dennic_user_service/internal/usecase"
	"dennic_user_service/internal/usecase/event"
	"time"

	"go.uber.org/zap"
)

const (
	AdminServiceName = "adminService"
	AdinSpanName     = "adminUsecase"
)

type adminRPC struct {
	logger         *zap.Logger
	admin          usecase.AdminStorageI
	brokerProducer event.BrokerProducer
}

var cfg = config.New()

func NewAdminRPC(logger *zap.Logger, admin usecase.AdminStorageI,
	brokerProducer event.BrokerProducer) pb.AdminServiceServer {
	return &adminRPC{
		logger:         logger,
		admin:          admin,
		brokerProducer: brokerProducer,
	}
}

func (a adminRPC) Create(ctx context.Context, admin *pb.Admin) (*pb.Admin, error) {

	ctx, span := otlp.Start(ctx, AdminServiceName, AdinSpanName+"Create")
	defer span.End()

	reqImageUrl := minio.RemoveImageUrl(admin.ImageUrl)
	req := entity.Admin{
		Id:            admin.Id,
		AdminOrder:    admin.AdminOrder,
		Role:          admin.Role,
		FirstName:     admin.FirstName,
		LastName:      admin.LastName,
		BirthDate:     admin.BirthDate,
		PhoneNumber:   admin.PhoneNumber,
		Email:         admin.Email,
		Password:      admin.Password,
		Gender:        admin.Gender,
		Salary:        admin.Salary,
		Biography:     admin.Biography,
		StartWorkYear: admin.StartWorkYear,
		EndWorkYear:   admin.EndWorkYear,
		WorkYears:     admin.WorkYears,
		RefreshToken:  admin.RefreshToken,
		ImageUrl:      reqImageUrl,
	}
	AdminId, err := a.admin.Create(ctx, &req)
	if err != nil {
		a.logger.Error("Create admin error", zap.Error(err))
		return nil, err
	}
	resp, err := a.admin.Get(ctx, &entity.FieldValueReq{
		Field:        "id",
		Value:        AdminId,
		DeleteStatus: false,
	})
	if err != nil {
		a.logger.Error("Create admin error", zap.Error(err))
		return nil, err
	}
	return &pb.Admin{
		Id:            resp.Id,
		AdminOrder:    resp.AdminOrder,
		Role:          resp.Role,
		FirstName:     resp.FirstName,
		LastName:      resp.LastName,
		BirthDate:     resp.BirthDate,
		PhoneNumber:   resp.PhoneNumber,
		Email:         resp.Email,
		Password:      resp.Password,
		Gender:        resp.Gender,
		Salary:        resp.Salary,
		Biography:     resp.Biography,
		StartWorkYear: resp.StartWorkYear,
		EndWorkYear:   resp.EndWorkYear,
		WorkYears:     resp.WorkYears,
		RefreshToken:  resp.RefreshToken,
		CreatedAt:     resp.CreatedAt.String(),
	}, nil
}

func (a adminRPC) Get(ctx context.Context, req *pb.GetAdminReq) (*pb.Admin, error) {

	ctx, span := otlp.Start(ctx, AdminServiceName, AdinSpanName+"Get")
	defer span.End()
	resp, err := a.admin.Get(ctx, &entity.FieldValueReq{
		Field:        req.Field,
		Value:        req.Value,
		DeleteStatus: req.IsActive,
	})

	var (
		respImageUrl string
	)
	if err != nil {
		a.logger.Error("get admin error", zap.Error(err))
		return nil, err
	}
	if resp.ImageUrl != "" {
		respImageUrl = minio.AddImageUrl(resp.ImageUrl, cfg.MinioService.Bucket.User)
	}
	response := &pb.Admin{
		Id:            resp.Id,
		AdminOrder:    resp.AdminOrder,
		Role:          resp.Role,
		FirstName:     resp.FirstName,
		LastName:      resp.LastName,
		BirthDate:     resp.BirthDate,
		PhoneNumber:   resp.PhoneNumber,
		Email:         resp.Email,
		Password:      resp.Password,
		Gender:        resp.Gender,
		Salary:        resp.Salary,
		Biography:     resp.Biography,
		StartWorkYear: resp.StartWorkYear,
		EndWorkYear:   resp.EndWorkYear,
		WorkYears:     resp.WorkYears,
		RefreshToken:  resp.RefreshToken,
		ImageUrl:      respImageUrl,
		CreatedAt:     resp.CreatedAt.String(),
		UpdatedAt:     resp.UpdatedAt.String(),
		DeletedAt:     resp.DeletedAt.String(),
	}

	if response.UpdatedAt == "0001-01-01 00:00:00 +0000 UTC" {
		response.UpdatedAt = ""
	}
	if response.DeletedAt == "0001-01-01 00:00:00 +0000 UTC" {
		response.DeletedAt = ""
	}

	return response, nil
}

func (a adminRPC) ListAdmins(ctx context.Context, req *pb.ListAdminsReq) (*pb.ListAdminsResp, error) {

	ctx, span := otlp.Start(ctx, AdminServiceName, AdinSpanName+"ListAdmins")
	defer span.End()
	resp, err := a.admin.List(ctx, &entity.GetAllReq{
		Page:         req.Page,
		Limit:        req.Limit,
		DeleteStatus: req.IsActive,
		Field:        req.Field,
		Value:        req.Value,
		OrderBy:      req.OrderBy,
	})

	if err != nil {
		a.logger.Error("get all admin error", zap.Error(err))
		return nil, err
	}

	var admins pb.ListAdminsResp
	var respImageUrl string

	for _, in := range resp {
		if in.ImageUrl != "" {
			respImageUrl = minio.AddImageUrl(in.ImageUrl, cfg.MinioService.Bucket.User)
		}
		if in.ImageUrl == "" {
			respImageUrl = minio.RemoveImageUrl(in.ImageUrl)
		}
		admin := &pb.Admin{
			Id:            in.Id,
			AdminOrder:    in.AdminOrder,
			Role:          in.Role,
			FirstName:     in.FirstName,
			LastName:      in.LastName,
			BirthDate:     in.BirthDate,
			PhoneNumber:   in.PhoneNumber,
			Email:         in.Email,
			Password:      in.Password,
			Gender:        in.Gender,
			Salary:        in.Salary,
			Biography:     in.Biography,
			StartWorkYear: in.StartWorkYear,
			EndWorkYear:   in.EndWorkYear,
			WorkYears:     in.WorkYears,
			RefreshToken:  in.RefreshToken,
			ImageUrl:      respImageUrl,
			CreatedAt:     in.CreatedAt.String(),
			UpdatedAt:     in.UpdatedAt.String(),
			DeletedAt:     in.DeletedAt.String(),
		}
		if in.UpdatedAt.String() == "0001-01-01 00:00:00 +0000 UTC" {
			admin.UpdatedAt = ""
		}
		if in.DeletedAt.String() == "0001-01-01 00:00:00 +0000 UTC" {
			admin.DeletedAt = ""
		}
		admins.Admins = append(admins.Admins, admin)
		admins.Count = uint64(in.Count)
	}

	return &admins, nil
}

func (a adminRPC) Update(ctx context.Context, admin *pb.Admin) (*pb.Admin, error) {

	ctx, span := otlp.Start(ctx, AdminServiceName, AdinSpanName+"Update")
	defer span.End()
	reqImageUrl := minio.RemoveImageUrl(admin.ImageUrl)
	req := entity.Admin{
		Id:            admin.Id,
		FirstName:     admin.FirstName,
		LastName:      admin.LastName,
		BirthDate:     admin.BirthDate,
		Gender:        admin.Gender,
		Salary:        admin.Salary,
		Biography:     admin.Biography,
		StartWorkYear: admin.StartWorkYear,
		EndWorkYear:   admin.EndWorkYear,
		WorkYears:     admin.WorkYears,
		ImageUrl:      reqImageUrl,
		UpdatedAt:     time.Now().Add(time.Hour * 5),
	}

	err := a.admin.Update(ctx, &req)

	if err != nil {
		a.logger.Error("update admin error", zap.Error(err))
		return nil, err
	}

	resp, err := a.admin.Get(ctx, &entity.FieldValueReq{
		Field:        "id",
		Value:        admin.Id,
		DeleteStatus: false,
	})

	var (
		respImageUrl string
	)
	if err != nil {
		a.logger.Error("Create admin error", zap.Error(err))
		return nil, err
	}
	if resp.ImageUrl != "" {
		respImageUrl = minio.AddImageUrl(resp.ImageUrl, cfg.MinioService.Bucket.User)
	}
	responer := &pb.Admin{
		Id:            resp.Id,
		AdminOrder:    resp.AdminOrder,
		Role:          resp.Role,
		FirstName:     resp.FirstName,
		LastName:      resp.LastName,
		BirthDate:     resp.BirthDate,
		PhoneNumber:   resp.PhoneNumber,
		Email:         resp.Email,
		Password:      resp.Password,
		Gender:        resp.Gender,
		Salary:        resp.Salary,
		Biography:     resp.Biography,
		StartWorkYear: resp.StartWorkYear,
		EndWorkYear:   resp.EndWorkYear,
		WorkYears:     resp.WorkYears,
		RefreshToken:  resp.RefreshToken,
		ImageUrl:      respImageUrl,
		CreatedAt:     resp.CreatedAt.String(),
		UpdatedAt:     resp.UpdatedAt.String(),
	}

	return responer, nil
}

func (a adminRPC) Delete(ctx context.Context, req *pb.DeleteAdminReq) (resp *pb.CheckAdminDeleteResp, err error) {

	ctx, span := otlp.Start(ctx, AdminServiceName, AdinSpanName+"Delete")
	defer span.End()
	status, err := a.admin.Delete(ctx, &entity.FieldValueReq{
		Field:        req.Field,
		Value:        req.Value,
		DeleteStatus: req.IsActive,
	})
	if err != nil {
		a.logger.Error("delete admin error", zap.Error(err))
		return nil, err
	}

	resp = &pb.CheckAdminDeleteResp{
		Status: status.Status,
	}

	return resp, nil
}

func (a adminRPC) CheckField(ctx context.Context, req *pb.CheckAdminFieldReq) (*pb.CheckAdminFieldResp, error) {

	ctx, span := otlp.Start(ctx, AdminServiceName, AdinSpanName+"CheckField")
	defer span.End()
	reqAdmin := entity.CheckFieldReq{
		Value: req.Value,
		Field: req.Field,
	}

	resp, err := a.admin.CheckField(ctx, &reqAdmin)
	if err != nil {
		a.logger.Error("delete admin error", zap.Error(err))
		return nil, err
	}
	response := &pb.CheckAdminFieldResp{
		Status: resp.Status,
	}

	return response, nil
}

func (a adminRPC) ChangePassword(ctx context.Context, phone *pb.ChangeAdminPasswordReq) (resp *pb.ChangeAdminPasswordResp, err error) {

	ctx, span := otlp.Start(ctx, AdminServiceName, AdinSpanName+"ChangePassword")
	defer span.End()
	req := entity.ChangeAdminPasswordReq{
		Email:       phone.Email,
		PhoneNumber: phone.PhoneNumber,
		Password:    phone.Password,
	}
	status, err := a.admin.ChangePassword(ctx, &req)
	if err != nil {
		a.logger.Error("delete admin error", zap.Error(err))
		return nil, err
	}
	resp = &pb.ChangeAdminPasswordResp{
		Status: status.Status,
	}

	return resp, nil
}

func (a adminRPC) UpdateRefreshToken(ctx context.Context, id *pb.UpdateRefreshTokenAdminReq) (resp *pb.UpdateRefreshTokenAdminResp, err error) {

	ctx, span := otlp.Start(ctx, AdminServiceName, AdinSpanName+"UpdateRefreshToken")
	defer span.End()
	req := entity.UpdateRefreshTokenReq{
		Id:           id.Id,
		RefreshToken: id.RefreshToken,
	}
	status, err := a.admin.UpdateRefreshToken(ctx, &req)
	if err != nil {
		a.logger.Error("delete admin error", zap.Error(err))
		return nil, err
	}

	resp = &pb.UpdateRefreshTokenAdminResp{
		Status: status.Status,
	}

	return resp, nil
}
