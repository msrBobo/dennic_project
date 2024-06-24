package suit_tests

import (
	"booking_service/internal/entity/doctor_availability"
	repo "booking_service/internal/infrastructure/repository/postgresql"
	"booking_service/internal/pkg/config"
	db "booking_service/internal/pkg/postgres"
	"context"
	"github.com/google/uuid"
	"github.com/rickb777/date"
	"github.com/stretchr/testify/suite"
	"strconv"
	"testing"
	"time"
)

type DoctorAvailabilityTestSite struct {
	suite.Suite
	Repository  *repo.DoctorAvailability
	CleanUpFunc func()
}

func (s *DoctorAvailabilityTestSite) SetupSuite() {
	pgPool, _ := db.New(config.New())
	s.Repository = repo.NewDoctorAvailability(pgPool)
	s.CleanUpFunc = pgPool.Close
}

func (s *DoctorAvailabilityTestSite) TestUserCRUD() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(2))
	defer cancel()

	doctorDate, _ := date.AutoParse("1231-02-02")
	startTime, _ := time.Parse("2006-01-02 15:04:05", "2000-01-01 14:14:14")
	endTime, _ := time.Parse("2006-01-02 15:04:05", "2000-01-01 11:11:11")

	createReq := &doctor_availability.CreateDoctorAvailability{
		DepartmentId: uuid.New().String(),
		DoctorId:     uuid.New().String(),
		DoctorDate:   doctorDate,
		StartTime:    startTime,
		EndTime:      endTime,
		Status:       "available",
	}

	createRes, err := s.Repository.CreateDoctorAvailability(ctx, createReq)
	s.Suite.NoError(err)
	s.Suite.NotNil(createRes)
	s.Suite.Equal(createRes.DepartmentId, createReq.DepartmentId)
	s.Suite.Equal(createRes.DoctorId, createReq.DoctorId)
	s.Suite.Equal(createRes.DoctorDate, createReq.DoctorDate)
	s.Suite.Equal(createRes.StartTime, createReq.StartTime)
	s.Suite.Equal(createRes.EndTime, createReq.EndTime)
	s.Suite.Equal(createRes.Status, createReq.Status)

	getRes, err := s.Repository.GetDoctorAvailability(ctx, &doctor_availability.FieldValueReq{
		Field:        "id",
		Value:        strconv.Itoa(int(createRes.Id)),
		DeleteStatus: true,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(getRes)
	s.Suite.Equal(getRes.DepartmentId, createReq.DepartmentId)
	s.Suite.Equal(getRes.DoctorId, createReq.DoctorId)
	s.Suite.Equal(getRes.DoctorDate, createReq.DoctorDate)
	s.Suite.Equal(getRes.StartTime, createReq.StartTime)
	s.Suite.Equal(getRes.EndTime, createReq.EndTime)
	s.Suite.Equal(getRes.Status, createReq.Status)

	getAllRes, err := s.Repository.GetAllDoctorAvailability(ctx, &doctor_availability.GetAllReq{
		Page:         1,
		Limit:        5,
		DeleteStatus: true,
		OrderBy:      "id",
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(getAllRes)

	newDoctorDate, _ := date.AutoParse("1231-02-02")
	newStartTime, _ := time.Parse("2006-01-02 15:04:05", "2000-01-01 14:14:14")
	newEndTime, _ := time.Parse("2006-01-02 15:04:05", "2000-01-01 11:11:11")

	updateReq := &doctor_availability.UpdateDoctorAvailability{
		Field:        "id",
		Value:        strconv.Itoa(int(getRes.Id)),
		DepartmentId: uuid.New().String(),
		DoctorId:     uuid.New().String(),
		DoctorDate:   newDoctorDate,
		StartTime:    newStartTime,
		EndTime:      newEndTime,
		Status:       "unavailable",
	}

	updateRes, err := s.Repository.UpdateDoctorAvailability(ctx, updateReq)
	s.Suite.NoError(err)
	s.Suite.NotNil(updateRes)
	s.Suite.Equal(updateRes.DepartmentId, updateReq.DepartmentId)
	s.Suite.Equal(updateRes.DoctorId, updateReq.DoctorId)
	s.Suite.Equal(updateRes.DoctorDate, updateReq.DoctorDate)
	s.Suite.Equal(updateRes.StartTime, updateReq.StartTime)
	s.Suite.Equal(updateRes.EndTime, updateReq.EndTime)
	s.Suite.Equal(updateRes.Status, updateReq.Status)

	softDelRes, err := s.Repository.DeleteDoctorAvailability(ctx, &doctor_availability.FieldValueReq{
		Field:        "doctor_id",
		Value:        createReq.DoctorId,
		DeleteStatus: false,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(softDelRes)
	s.Suite.Equal(softDelRes.Status, true)

	hardDelRes, err := s.Repository.DeleteDoctorAvailability(ctx, &doctor_availability.FieldValueReq{
		Field:        "doctor_id",
		Value:        createReq.DoctorId,
		DeleteStatus: true,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(hardDelRes)
	s.Suite.Equal(hardDelRes.Status, true)

}

func (s *DoctorAvailabilityTestSite) TearDownSuite() {
	s.CleanUpFunc()
}

func TestDoctorAvailabilityTestSuite(t *testing.T) {
	suite.Run(t, new(DoctorAvailabilityTestSite))
}
