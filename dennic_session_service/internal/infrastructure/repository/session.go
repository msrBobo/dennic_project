package repository

import (
	"context"
	pb "dennic_session_service/internal/entity"
)

type SessionRepo interface {
	CreateSession(context.Context, *pb.SessionRequests) (*pb.Session, error)
	GetSessionById(context.Context, *pb.StrReq) (*pb.Session, error)
	DeleteSessionById(context.Context, *pb.StrReq) (*pb.Empty, error)
	DeleteSessionByUserId(context.Context, *pb.StrUserReq) (*pb.Empty, error)
	GetUserSessions(context.Context, *pb.StrUserReq) (*pb.UserSessionsList, error)
}
