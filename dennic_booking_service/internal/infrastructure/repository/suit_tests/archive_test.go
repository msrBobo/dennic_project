package suit_tests

import (
	"booking_service/internal/entity/archive"
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

type BookingArchiveTestSite struct {
	suite.Suite
	Repository         *repo.BookingArchive
	DoctorAvailability *repo.DoctorAvailability
	CleanUpFunc        func()
}

func (s *BookingArchiveTestSite) SetupSuite() {
	pgPool, _ := db.New(config.New())
	s.Repository = repo.NewBookingArchive(pgPool)
	s.DoctorAvailability = repo.NewDoctorAvailability(pgPool)
	s.CleanUpFunc = pgPool.Close
}

func (s *BookingArchiveTestSite) TestUserCRUD() {
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

	createRes, err := s.DoctorAvailability.CreateDoctorAvailability(ctx, createReq)
	s.Suite.NoError(err)
	s.Suite.NotNil(createRes)
	s.Suite.Equal(createRes.DepartmentId, createReq.DepartmentId)
	s.Suite.Equal(createRes.DoctorId, createReq.DoctorId)
	s.Suite.Equal(createRes.DoctorDate, createReq.DoctorDate)
	s.Suite.Equal(createRes.StartTime, createReq.StartTime)
	s.Suite.Equal(createRes.EndTime, createReq.EndTime)
	s.Suite.Equal(createRes.Status, createReq.Status)

	createArchiveReq := &archive.CreatedArchive{
		DoctorAvailabilityId: createRes.Id,
		StartTime:            startTime,
		EndTime:              endTime,
		PatientProblem:       "No Problem",
		Status:               "attended",
		PaymentType:          "card",
		PaymentAmount:        10,
	}

	createArchiveRes, err := s.Repository.CreateArchive(ctx, createArchiveReq)
	s.Suite.NoError(err)
	s.Suite.NotNil(createArchiveRes)
	s.Suite.Equal(createArchiveRes.DoctorAvailabilityId, createArchiveReq.DoctorAvailabilityId)
	s.Suite.Equal(createArchiveRes.StartTime, createArchiveReq.StartTime)
	s.Suite.Equal(createArchiveRes.EndTime, createArchiveReq.EndTime)
	s.Suite.Equal(createArchiveRes.PatientProblem, createArchiveReq.PatientProblem)
	s.Suite.Equal(createArchiveRes.Status, createArchiveReq.Status)
	s.Suite.Equal(createArchiveRes.PaymentType, createArchiveReq.PaymentType)
	s.Suite.Equal(createArchiveRes.PaymentAmount, createArchiveReq.PaymentAmount)

	getRes, err := s.Repository.GetArchive(ctx, &archive.FieldValueReq{
		Field:        "id",
		Value:        strconv.Itoa(int(createArchiveRes.Id)),
		DeleteStatus: true,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(getRes)
	s.Suite.Equal(getRes.Id, createArchiveRes.Id)
	s.Suite.Equal(getRes.DoctorAvailabilityId, createArchiveRes.DoctorAvailabilityId)
	s.Suite.Equal(getRes.StartTime, createArchiveRes.StartTime)
	s.Suite.Equal(getRes.EndTime, createArchiveRes.EndTime)
	s.Suite.Equal(getRes.PatientProblem, createArchiveRes.PatientProblem)
	s.Suite.Equal(getRes.Status, createArchiveRes.Status)
	s.Suite.Equal(getRes.PaymentType, createArchiveRes.PaymentType)
	s.Suite.Equal(getRes.PaymentAmount, createArchiveRes.PaymentAmount)

	getAllRes, err := s.Repository.GetAllArchive(ctx, &archive.GetAllArchives{
		Page:         1,
		Limit:        5,
		DeleteStatus: true,
		OrderBy:      "status",
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(getAllRes)

	newStartTime, _ := time.Parse("2006-01-02 15:04:05", "2000-01-01 14:14:14")
	newEndTime, _ := time.Parse("2006-01-02 15:04:05", "2000-01-01 11:11:11")

	updateReq := &archive.UpdateArchive{
		Field:                "id",
		Value:                strconv.Itoa(int(getRes.Id)),
		DoctorAvailabilityId: createRes.Id,
		StartTime:            newStartTime,
		EndTime:              newEndTime,
		PatientProblem:       "Yes Problem",
		Status:               "cancelled",
		PaymentType:          "cash",
		PaymentAmount:        12,
	}
	updateRes, err := s.Repository.UpdateArchive(ctx, updateReq)
	s.Suite.NoError(err)
	s.Suite.NotNil(updateRes)
	s.Suite.Equal(updateRes.Id, getRes.Id)
	s.Suite.Equal(updateRes.DoctorAvailabilityId, updateReq.DoctorAvailabilityId)
	s.Suite.Equal(updateRes.StartTime, updateReq.StartTime)
	s.Suite.Equal(updateRes.EndTime, updateReq.EndTime)
	s.Suite.Equal(updateRes.PatientProblem, updateReq.PatientProblem)
	s.Suite.Equal(updateRes.Status, updateReq.Status)
	s.Suite.Equal(updateRes.PaymentType, updateReq.PaymentType)
	s.Suite.Equal(updateRes.PaymentAmount, updateReq.PaymentAmount)

	//

	softDeleteArchiveRes, err := s.Repository.DeleteArchive(ctx, &archive.FieldValueReq{
		Field:        "id",
		Value:        strconv.Itoa(int(createArchiveRes.Id)),
		DeleteStatus: true,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(softDeleteArchiveRes)
	s.Suite.Equal(softDeleteArchiveRes.Status, true)

	hardDeleteArchiveRes, err := s.Repository.DeleteArchive(ctx, &archive.FieldValueReq{
		Field:        "id",
		Value:        strconv.Itoa(int(createArchiveRes.Id)),
		DeleteStatus: true,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(hardDeleteArchiveRes)
	s.Suite.Equal(hardDeleteArchiveRes.Status, true)

	hardDelRes, err := s.DoctorAvailability.DeleteDoctorAvailability(ctx, &doctor_availability.FieldValueReq{
		Field:        "doctor_id",
		Value:        createReq.DoctorId,
		DeleteStatus: true,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(hardDelRes)
	s.Suite.Equal(hardDelRes.Status, true)
}

func (s *BookingArchiveTestSite) TearDownSuite() {
	s.CleanUpFunc()
}

func TestUBookingArchiveTestSuite(t *testing.T) {
	suite.Run(t, new(BookingArchiveTestSite))
}
