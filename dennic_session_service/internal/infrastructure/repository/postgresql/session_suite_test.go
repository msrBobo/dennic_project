package postgresql

import (
	"context"
	"dennic_session_service/internal/entity"
	"dennic_session_service/internal/pkg/config"
	db "dennic_session_service/internal/pkg/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
	"time"
)

type DepartmentTestSuite struct {
	suite.Suite
	CleanUpFunc func()
	Repository  *SessionRepository
}

func (s *DepartmentTestSuite) SetupTest() {
	pgPool, err := db.New(config.New())
	if err != nil {
		log.Fatal(err)
		return
	}
	s.Repository = NewSessionRepository(pgPool)
	s.CleanUpFunc = pgPool.Close
}

func (s *DepartmentTestSuite) TestDepartmentCrud() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(2))
	defer cancel()

	session := &entity.SessionRequests{
		Id:           uuid.NewString(),
		IpAddress:    "Test password",
		UserId:       uuid.NewString(),
		FcmToken:     "Test Fcm Token",
		PlatformName: "Test Platform Name",
		PlatformType: "mobile",
	}
	respDep, err := s.Repository.CreateSession(ctx, session)
	s.Suite.NoError(err)
	s.Suite.NotNil(respDep)
	s.Suite.Equal(respDep.Id, session.Id)
	s.Suite.Equal(respDep.IpAddress, session.IpAddress)
	s.Suite.Equal(respDep.UserId, session.UserId)
	s.Suite.Equal(respDep.FcmToken, session.FcmToken)
	s.Suite.Equal(respDep.PlatformName, session.PlatformName)
	s.Suite.Equal(respDep.PlatformName, session.PlatformName)
	s.Suite.Equal(respDep.PlatformType, session.PlatformType)

	getSession, err := s.Repository.GetSessionById(ctx, &entity.StrReq{
		Id: session.Id,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(getSession)
	s.Suite.Equal(getSession.Id, session.Id)
	s.Suite.Equal(getSession.IpAddress, session.IpAddress)
	s.Suite.Equal(getSession.UserId, session.UserId)
	s.Suite.Equal(getSession.FcmToken, session.FcmToken)
	s.Suite.Equal(getSession.PlatformName, session.PlatformName)
	s.Suite.Equal(getSession.PlatformName, session.PlatformName)
	s.Suite.Equal(getSession.PlatformType, session.PlatformType)

	respAll, err := s.Repository.GetUserSessions(ctx, &entity.StrUserReq{UserId: session.UserId})
	s.Suite.NoError(err)
	s.Suite.NotNil(respAll)

	_, err = s.Repository.DeleteSessionById(ctx, &entity.StrReq{
		Id: session.Id,
	})
	s.Suite.NoError(err)

}

func (s *DepartmentTestSuite) TearDownTest() {
	s.CleanUpFunc()
}

func TestDepartmentTestSuite(t *testing.T) {
	suite.Run(t, new(DepartmentTestSuite))
}
