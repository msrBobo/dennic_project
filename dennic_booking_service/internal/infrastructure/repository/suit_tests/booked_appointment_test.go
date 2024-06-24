package suit_tests

import (
	"booking_service/internal/entity/booked_appointments"
	"booking_service/internal/entity/patients"
	repo "booking_service/internal/infrastructure/repository/postgresql"
	"booking_service/internal/pkg/config"
	db "booking_service/internal/pkg/postgres"
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rickb777/date"
	"github.com/stretchr/testify/suite"
)

type BookingAppointmentTestSite struct {
	suite.Suite
	Repository  *repo.BookingAppointment
	Patient     *repo.BookingPatients
	CleanUpFunc func()
}

func (s *BookingAppointmentTestSite) SetupSuite() {
	pgPool, _ := db.New(config.New())
	s.Repository = repo.NewBookingAppointment(pgPool)
	s.Patient = repo.NewBookingPatients(pgPool)
	s.CleanUpFunc = pgPool.Close
}

func (s *BookingAppointmentTestSite) TestUserCRUD() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(2))
	defer cancel()

	brithDate := date.Today()
	patient := &patients.CreatedPatient{
		Id:             uuid.New().String(),
		FirstName:      "Husanboy",
		LastName:       "Gofurov",
		BirthDate:      brithDate,
		Gender:         "male",
		BloodGroup:     "A+",
		PhoneNumber:    "+998950230605",
		City:           "Andijon",
		Country:        "Uzbekistan",
		Address:        "Shahrixon",
		PatientProblem: "Now Problem",
	}
	patientRes, err := s.Patient.CreatePatient(ctx, patient)
	s.Suite.NoError(err)
	s.Suite.NotNil(patientRes)
	s.Suite.Equal(patientRes.Id, patient.Id)
	s.Suite.Equal(patientRes.FirstName, patient.FirstName)
	s.Suite.Equal(patientRes.LastName, patient.LastName)
	s.Suite.Equal(patientRes.BirthDate, brithDate)
	s.Suite.Equal(patientRes.Gender, patient.Gender)
	s.Suite.Equal(patientRes.BloodGroup, patient.BloodGroup)
	s.Suite.Equal(patientRes.PhoneNumber, patient.PhoneNumber)
	s.Suite.Equal(patientRes.City, patient.City)
	s.Suite.Equal(patientRes.Country, patient.Country)
	s.Suite.Equal(patientRes.Address, patient.Address)
	s.Suite.Equal(patientRes.PatientProblem, patient.PatientProblem)

	appDate, _ := date.AutoParse("1221-12-12")
	appTime, _ := time.Parse("2006-01-02 15:04:05", "2000-01-01 12:12:12")
	axpTime, _ := time.Parse("2006-01-02 15:04:05", "2000-01-01 01:01:01")

	createReq := booked_appointments.CreateAppointment{
		DepartmentId:    uuid.New().String(),
		DoctorId:        uuid.New().String(),
		PatientId:       patientRes.Id,
		AppointmentDate: appDate,
		AppointmentTime: appTime,
		Duration:        11,
		Key:             "ABC",
		ExpiresAt:       axpTime,
		Status:          "waiting",
		PaymentType:     "cash",
		PaymentAmount:   100000,
		PatientProblem:  "Now Problem",
	}

	createRes, err := s.Repository.CreateAppointment(ctx, &createReq)
	s.Suite.NoError(err)
	s.Suite.NotNil(createRes)
	s.Suite.Equal(createRes.DoctorId, createReq.DoctorId)
	s.Suite.Equal(createRes.PatientId, createReq.PatientId)
	s.Suite.Equal(createRes.AppointmentDate, createReq.AppointmentDate)
	s.Suite.Equal(createRes.AppointmentTime, createReq.AppointmentTime)
	s.Suite.Equal(createRes.Duration, createReq.Duration)
	s.Suite.Equal(createRes.Key, createReq.Key)
	s.Suite.Equal(createRes.ExpiresAt, createReq.ExpiresAt)
	s.Suite.Equal(createRes.Status, createReq.Status)

	getRes, err := s.Repository.GetAppointment(ctx, &booked_appointments.FieldValueReq{
		Field:        "patient_id",
		Value:        createRes.PatientId,
		DeleteStatus: true,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(getRes)
	s.Suite.Equal(getRes.Id, createRes.Id)
	s.Suite.Equal(getRes.DoctorId, createRes.DoctorId)
	s.Suite.Equal(getRes.PatientId, createRes.PatientId)
	s.Suite.Equal(getRes.AppointmentDate, createRes.AppointmentDate)
	s.Suite.Equal(getRes.AppointmentTime, createRes.AppointmentTime)
	s.Suite.Equal(getRes.Duration, createRes.Duration)
	s.Suite.Equal(getRes.Key, createRes.Key)
	s.Suite.Equal(getRes.ExpiresAt, createRes.ExpiresAt)
	s.Suite.Equal(getRes.Status, createRes.Status)

	getAllRes, err := s.Repository.GetAllAppointment(ctx, &booked_appointments.GetAllAppointment{
		Page:         1,
		Limit:        10,
		DeleteStatus: true,
		OrderBy:      "id",
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(getAllRes)

	newAppDate, _ := date.AutoParse("1231-02-02")
	newAppTime, _ := time.Parse("2006-01-02 15:04:05", "2000-01-01 14:14:14")
	newAxpTime, _ := time.Parse("2006-01-02 15:04:05", "2000-01-01 11:11:11")

	upReq := &booked_appointments.UpdateAppointment{
		Field:           "id",
		Value:           strconv.Itoa(int(getRes.Id)),
		AppointmentDate: newAppDate,
		AppointmentTime: newAppTime,
		Duration:        12,
		Key:             "New",
		ExpiresAt:       newAxpTime,
		Status:          "no_show",
	}
	upRes, err := s.Repository.UpdateAppointment(ctx, upReq)
	s.Suite.NoError(err)
	s.Suite.NotNil(upRes)
	s.Suite.Equal(upRes.Id, getRes.Id)
	s.Suite.Equal(upRes.AppointmentDate, newAppDate)
	s.Suite.Equal(upRes.AppointmentTime, newAppTime)
	s.Suite.Equal(upRes.Duration, upReq.Duration)
	s.Suite.Equal(upRes.Key, upReq.Key)
	s.Suite.Equal(upRes.ExpiresAt, newAxpTime)
	s.Suite.Equal(upRes.Status, upReq.Status)

	softDelRes, err := s.Repository.DeleteAppointment(ctx, &booked_appointments.FieldValueReq{
		Field:        "id",
		Value:        strconv.Itoa(int(getRes.Id)),
		DeleteStatus: false,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(softDelRes)
	s.Suite.Equal(softDelRes.Status, true)

	hardDelRes, err := s.Repository.DeleteAppointment(ctx, &booked_appointments.FieldValueReq{
		Field:        "id",
		Value:        strconv.Itoa(int(getRes.Id)),
		DeleteStatus: false,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(hardDelRes)
	s.Suite.Equal(hardDelRes.Status, true)

	hardDeleteRes, err := s.Patient.DeletePatient(ctx, &patients.FieldValueReq{
		Field:        "id",
		Value:        patient.Id,
		DeleteStatus: true,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(hardDeleteRes)
	s.Suite.Equal(hardDeleteRes.Status, true)
}

func (s *BookingAppointmentTestSite) TearDownSuite() {
	s.CleanUpFunc()
}

func TestBookingAppointmentTestSuite(t *testing.T) {
	suite.Run(t, new(BookingAppointmentTestSite))
}
