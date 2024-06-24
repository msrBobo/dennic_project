package services

import (
	"context"
	pb "dennic_session_service/genproto/session_service"
	"dennic_session_service/internal/entity"
	"dennic_session_service/internal/pkg/otlp"
	"dennic_session_service/internal/usecase"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

const (
	serviceNameSessionDelivery           = "Delivery_session_service"
	serviceNameSessionDeliveryRepoPrefix = "Delivery_session_service"
)

type sessionRPC struct {
	logger  *zap.Logger
	session usecase.SessionUsecase
}

func SessionRPC(logger *zap.Logger, session usecase.SessionUsecase) pb.SessionServiceServer {
	return &sessionRPC{
		logger:  logger,
		session: session,
	}
}

func (s sessionRPC) CreateSession(ctx context.Context, requests *pb.SessionRequests) (*pb.Session, error) {
	ctx, span := otlp.Start(ctx, serviceNameSessionDelivery, serviceNameSessionDeliveryRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("CreateSession").String(requests.Id))
	defer span.End()

	req := entity.SessionRequests{
		Id:           requests.Id,
		IpAddress:    requests.IpAddress,
		UserId:       requests.UserId,
		FcmToken:     requests.FcmToken,
		PlatformName: requests.PlatformName,
		PlatformType: requests.PlatformType,
	}
	resp, err := s.session.CreateSession(ctx, &req)
	if err != nil {
		s.logger.Error("CreateSession", zap.Error(err))
	}
	return &pb.Session{
		Id:           resp.Id,
		Order:        resp.Order,
		IpAddress:    resp.IpAddress,
		UserId:       resp.UserId,
		FcmToken:     resp.FcmToken,
		PlatformName: resp.PlatformName,
		PlatformType: resp.PlatformType,
		LoginAt:      resp.LoginAt.String(),
		CreatedAt:    resp.CreatedAt.String(),
		UpdatedAt:    resp.UpdatedAt.String(),
		DeletedAt:    resp.DeletedAt.String(),
	}, nil
}

func (s sessionRPC) GetSessionById(ctx context.Context, req *pb.StrReq) (*pb.Session, error) {

	ctx, span := otlp.Start(ctx, serviceNameSessionDelivery, serviceNameSessionDeliveryRepoPrefix+"Get")
	span.SetAttributes(attribute.Key("GetSessionById").String(req.Id))
	defer span.End()

	resp, err := s.session.GetSessionById(ctx, &entity.StrReq{Id: req.Id})
	if err != nil {
		s.logger.Error("GetSessionById", zap.Error(err))
	}
	return &pb.Session{
		Id:           resp.Id,
		Order:        resp.Order,
		IpAddress:    resp.IpAddress,
		UserId:       resp.UserId,
		FcmToken:     resp.FcmToken,
		PlatformName: resp.PlatformName,
		PlatformType: resp.PlatformType,
		LoginAt:      resp.LoginAt.String(),
		CreatedAt:    resp.CreatedAt.String(),
		UpdatedAt:    resp.UpdatedAt.String(),
		DeletedAt:    resp.DeletedAt.String(),
	}, nil
}

func (s sessionRPC) DeleteSessionById(ctx context.Context, req *pb.StrReq) (*pb.Empty, error) {
	ctx, span := otlp.Start(ctx, serviceNameSessionDelivery, serviceNameSessionDeliveryRepoPrefix+"Delete")
	span.SetAttributes(attribute.Key("DeleteSessionById").String(req.Id))
	defer span.End()
	_, err := s.session.DeleteSessionById(ctx, &entity.StrReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (s sessionRPC) DeleteSessionByUserId(ctx context.Context, req *pb.StrUserReq) (*pb.Empty, error) {
	ctx, span := otlp.Start(ctx, serviceNameSessionDelivery, serviceNameSessionDeliveryRepoPrefix+"Delete")
	span.SetAttributes(attribute.Key("DeleteSessionByUserId").String(req.UserId))
	defer span.End()
	_, err := s.session.DeleteSessionByUserId(ctx, &entity.StrUserReq{UserId: req.UserId})
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (s sessionRPC) GetUserSessions(ctx context.Context, req *pb.StrUserReq) (*pb.UserSessionsList, error) {
	ctx, span := otlp.Start(ctx, serviceNameSessionDelivery, serviceNameSessionDeliveryRepoPrefix+"Get")
	span.SetAttributes(attribute.Key("GetUserSessions").String(req.UserId))
	defer span.End()
	resp, err := s.session.GetUserSessions(ctx, &entity.StrUserReq{UserId: req.UserId, IsActive: req.IsActive})
	if err != nil {
		s.logger.Error("GetUserSessions", zap.Error(err))
	}
	var respSessions pb.UserSessionsList
	for _, session := range resp.UserSessions {
		respSessions.UserSessions = append(respSessions.UserSessions, &pb.Session{
			Id:           session.Id,
			Order:        session.Order,
			IpAddress:    session.IpAddress,
			UserId:       session.UserId,
			FcmToken:     session.FcmToken,
			PlatformName: session.PlatformName,
			PlatformType: session.PlatformType,
			LoginAt:      session.LoginAt.String(),
			CreatedAt:    session.CreatedAt.String(),
			UpdatedAt:    session.UpdatedAt.String(),
			DeletedAt:    session.DeletedAt.String(),
		})
	}
	respSessions.Count = resp.Count

	return &respSessions, nil
}

func (s sessionRPC) HasUserSession(ctx context.Context, req *pb.StrUserReq) (*pb.SessionExistsResponse, error) {
	ctx, span := otlp.Start(ctx, serviceNameSessionDelivery, serviceNameSessionDeliveryRepoPrefix+"Has user")
	span.SetAttributes(attribute.Key("HasUserSession").String(req.UserId))
	defer span.End()
	resp, err := s.GetUserSessions(ctx, req)
	if err != nil {
		s.logger.Error("HasUserSession", zap.Error(err))
		return nil, err
	}
	if len(resp.UserSessions) == 0 {
		return &pb.SessionExistsResponse{IsExists: false}, nil
	}
	return &pb.SessionExistsResponse{IsExists: true}, nil
}
