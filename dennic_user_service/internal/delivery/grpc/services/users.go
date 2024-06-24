package services

import (
	"context"
	pb "dennic_user_service/genproto/user_service"
	"dennic_user_service/internal/entity"
	"dennic_user_service/internal/pkg/minio"
	"dennic_user_service/internal/pkg/otlp"
	"dennic_user_service/internal/usecase"
	"dennic_user_service/internal/usecase/event"
	"time"

	"go.uber.org/zap"
)

const (
	UserServiceName = "userService"
	UserSpanName    = "userUsecase"
)

type userRPC struct {
	logger         *zap.Logger
	user           usecase.UserStorageI
	brokerProducer event.BrokerProducer
}

func NewUserRPC(logger *zap.Logger, user usecase.UserStorageI,
	brokerProducer event.BrokerProducer) pb.UserServiceServer {
	return &userRPC{
		logger:         logger,
		user:           user,
		brokerProducer: brokerProducer,
	}
}

func (u userRPC) Create(ctx context.Context, user *pb.User) (*pb.User, error) {

	ctx, span := otlp.Start(ctx, UserServiceName, UserSpanName+"Create")
	defer span.End()

	reqImageUrl := minio.RemoveImageUrl(user.ImageUrl)
	req := entity.User{
		Id:           user.Id,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		BirthDate:    user.BirthDate,
		PhoneNumber:  user.PhoneNumber,
		Password:     user.Password,
		Gender:       user.Gender,
		RefreshToken: user.RefreshToken,
		ImageUrl:     reqImageUrl,
	}
	UserId, err := u.user.Create(ctx, &req)
	if err != nil {
		return nil, err
	}

	resp, err := u.user.Get(ctx, &entity.FieldValueReq{
		Field:        "id",
		Value:        UserId,
		DeleteStatus: false,
	})
	if err != nil {
		return nil, err
	}

	return &pb.User{
		Id:           resp.Id,
		UserOrder:    resp.UserOrder,
		FirstName:    resp.FirstName,
		LastName:     resp.LastName,
		BirthDate:    resp.BirthDate,
		PhoneNumber:  resp.PhoneNumber,
		Password:     resp.Password,
		Gender:       resp.Gender,
		RefreshToken: resp.RefreshToken,
		CreatedAt:    resp.CreatedAt.String(),
	}, nil
}

func (u userRPC) Get(ctx context.Context, req *pb.GetUserReq) (*pb.User, error) {

	ctx, span := otlp.Start(ctx, UserServiceName, UserSpanName+"Get")
	defer span.End()
	var (
		respImageUrl string
	)
	resp, err := u.user.Get(ctx, &entity.FieldValueReq{
		Field:        req.Field,
		Value:        req.Value,
		DeleteStatus: req.IsActive,
	})

	if err != nil {
		return nil, err
	}
	if resp.ImageUrl != "" {
		respImageUrl = minio.AddImageUrl(resp.ImageUrl, cfg.MinioService.Bucket.User)
	}
	response := &pb.User{
		Id:           resp.Id,
		UserOrder:    resp.UserOrder,
		FirstName:    resp.FirstName,
		LastName:     resp.LastName,
		BirthDate:    resp.BirthDate,
		PhoneNumber:  resp.PhoneNumber,
		Password:     resp.Password,
		Gender:       resp.Gender,
		RefreshToken: resp.RefreshToken,
		ImageUrl:     respImageUrl,
		CreatedAt:    resp.CreatedAt.String(),
		UpdatedAt:    resp.UpdatedAt.String(),
		DeletedAt:    resp.DeletedAt.String(),
	}

	if response.UpdatedAt == "0001-01-01 00:00:00 +0000 UTC" {
		response.UpdatedAt = ""
	}
	if response.DeletedAt == "0001-01-01 00:00:00 +0000 UTC" {
		response.DeletedAt = ""
	}

	return response, nil
}

func (u userRPC) ListUsers(ctx context.Context, req *pb.ListUsersReq) (*pb.ListUsersResp, error) {

	ctx, span := otlp.Start(ctx, UserServiceName, UserSpanName+"ListUsers")
	defer span.End()
	resp, err := u.user.List(ctx, &entity.GetAllReq{
		Page:         req.Page,
		Limit:        req.Limit,
		DeleteStatus: req.IsActive,
		Field:        req.Field,
		Value:        req.Value,
		OrderBy:      req.OrderBy,
	})

	if err != nil {
		return nil, err
	}

	var (
		users        pb.ListUsersResp
		respImageUrl string
	)

	for _, in := range resp {
		if in.ImageUrl != "" {
			respImageUrl = minio.AddImageUrl(in.ImageUrl, cfg.MinioService.Bucket.User)
		}
		if in.ImageUrl == "" {
			respImageUrl = minio.RemoveImageUrl(in.ImageUrl)
		}
		user := &pb.User{
			Id:           in.Id,
			UserOrder:    in.UserOrder,
			FirstName:    in.FirstName,
			LastName:     in.LastName,
			BirthDate:    in.BirthDate,
			PhoneNumber:  in.PhoneNumber,
			Password:     in.Password,
			Gender:       in.Gender,
			RefreshToken: in.RefreshToken,
			ImageUrl:     respImageUrl,
			CreatedAt:    in.CreatedAt.String(),
			UpdatedAt:    in.UpdatedAt.String(),
			DeletedAt:    in.DeletedAt.String(),
		}

		if in.UpdatedAt.String() == "0001-01-01 00:00:00 +0000 UTC" {
			user.UpdatedAt = ""
		}
		if in.DeletedAt.String() == "0001-01-01 00:00:00 +0000 UTC" {
			user.DeletedAt = ""
		}
		users.Users = append(users.Users, user)
		users.Count = uint64(in.Count)
	}

	return &users, nil
}

func (u userRPC) Update(ctx context.Context, user *pb.User) (*pb.User, error) {

	ctx, span := otlp.Start(ctx, UserServiceName, UserSpanName+"Update")
	defer span.End()
	reqImageUrl := minio.RemoveImageUrl(user.ImageUrl)
	req := entity.User{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		BirthDate: user.BirthDate,
		Gender:    user.Gender,
		ImageUrl:  reqImageUrl,
		UpdatedAt: time.Now().Add(time.Hour * 5),
	}

	err := u.user.Update(ctx, &req)

	if err != nil {
		return nil, err
	}
	var (
		respImageUrl string
	)
	resp, err := u.user.Get(ctx, &entity.FieldValueReq{
		Field:        "id",
		Value:        user.Id,
		DeleteStatus: false,
	})
	if err != nil {
		return nil, err
	}
	if resp.ImageUrl != "" {
		respImageUrl = minio.AddImageUrl(resp.ImageUrl, cfg.MinioService.Bucket.User)
	}
	response := &pb.User{
		Id:           resp.Id,
		UserOrder:    resp.UserOrder,
		FirstName:    resp.FirstName,
		LastName:     resp.LastName,
		BirthDate:    resp.BirthDate,
		PhoneNumber:  resp.LastName,
		Password:     resp.Password,
		Gender:       resp.Gender,
		RefreshToken: resp.RefreshToken,
		ImageUrl:     respImageUrl,
		CreatedAt:    resp.CreatedAt.String(),
		UpdatedAt:    resp.UpdatedAt.String(),
	}
	return response, nil
}

func (u userRPC) Delete(ctx context.Context, req *pb.DeleteUserReq) (resp *pb.CheckDeleteUserResp, err error) {

	ctx, span := otlp.Start(ctx, UserServiceName, UserSpanName+"Delete")
	defer span.End()
	status, err := u.user.Delete(ctx, &entity.FieldValueReq{
		Field:        req.Field,
		Value:        req.Value,
		DeleteStatus: req.IsActive,
	})
	if err != nil {
		return nil, err
	}
	resp = &pb.CheckDeleteUserResp{
		Status: status.Status,
	}

	return resp, nil
}

func (u userRPC) CheckField(ctx context.Context, req *pb.CheckFieldUserReq) (*pb.CheckFieldUserResp, error) {

	ctx, span := otlp.Start(ctx, UserServiceName, UserSpanName+"CheckField")
	defer span.End()
	reqUser := entity.CheckFieldReq{
		Value: req.Value,
		Field: req.Field,
	}

	resp, err := u.user.CheckField(ctx, &reqUser)
	if err != nil {
		return nil, err
	}
	response := &pb.CheckFieldUserResp{
		Status: resp.Status,
	}

	return response, nil
}

func (u userRPC) ChangePassword(ctx context.Context, phone *pb.ChangeUserPasswordReq) (resp *pb.ChangeUserPasswordResp, err error) {

	ctx, span := otlp.Start(ctx, UserServiceName, UserSpanName+"ChangePassword")
	defer span.End()
	req := entity.ChangeUserPasswordReq{
		PhoneNumber: phone.PhoneNumber,
		Password:    phone.Password,
	}
	status, err := u.user.ChangePassword(ctx, &req)
	if err != nil {
		return nil, err
	}
	resp = &pb.ChangeUserPasswordResp{
		Status: status.Status,
	}

	return resp, nil
}

func (u userRPC) UpdateRefreshToken(ctx context.Context, id *pb.UpdateRefreshTokenUserReq) (resp *pb.UpdateRefreshTokenUserResp, err error) {

	ctx, span := otlp.Start(ctx, UserServiceName, UserSpanName+"UpdateRefreshToken")
	defer span.End()
	req := entity.UpdateRefreshTokenReq{
		Id:           id.Id,
		RefreshToken: id.RefreshToken,
	}
	status, err := u.user.UpdateRefreshToken(ctx, &req)
	if err != nil {
		return nil, err
	}

	resp = &pb.UpdateRefreshTokenUserResp{
		Status: status.Status,
	}

	return resp, nil
}
