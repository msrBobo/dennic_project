package suit_tests

import (
	"booking_service/internal/entity/archive"
	"booking_service/internal/entity/doctor_availability"
	"booking_service/internal/entity/doctor_notes"
	"booking_service/internal/entity/patients"
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

type DoctorNotesTestSite struct {
	suite.Suite
	Repository         *repo.DoctorNotes
	DoctorAvailability *repo.DoctorAvailability
	Archive            *repo.BookingArchive
	Patient            *repo.BookingPatients
	CleanUpFunc        func()
}

func (s *DoctorNotesTestSite) SetupSuite() {
	pgPool, _ := db.New(config.New())
	s.Repository = repo.NewDoctorNotes(pgPool)
	s.DoctorAvailability = repo.NewDoctorAvailability(pgPool)
	s.Archive = repo.NewBookingArchive(pgPool)
	s.Patient = repo.NewBookingPatients(pgPool)
	s.CleanUpFunc = pgPool.Close
}

func (s *DoctorNotesTestSite) TestUserCRUD() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(2))
	defer cancel()

	brithDate := date.Today()
	patient := &patients.CreatedPatient{
		Id:             uuid.New().String(),
		FirstName:      "Husan",
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

	createPatientRes, err := s.Patient.CreatePatient(ctx, patient)
	s.Suite.NoError(err)
	s.Suite.NotNil(createPatientRes)
	s.Suite.Equal(createPatientRes.Id, patient.Id)
	s.Suite.Equal(createPatientRes.FirstName, patient.FirstName)
	s.Suite.Equal(createPatientRes.LastName, patient.LastName)
	s.Suite.Equal(createPatientRes.BirthDate, brithDate)
	s.Suite.Equal(createPatientRes.Gender, patient.Gender)
	s.Suite.Equal(createPatientRes.BloodGroup, patient.BloodGroup)
	s.Suite.Equal(createPatientRes.PhoneNumber, patient.PhoneNumber)
	s.Suite.Equal(createPatientRes.City, patient.City)
	s.Suite.Equal(createPatientRes.Country, patient.Country)
	s.Suite.Equal(createPatientRes.Address, patient.Address)
	s.Suite.Equal(createPatientRes.PatientProblem, patient.PatientProblem)

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

	createArchiveRes, err := s.Archive.CreateArchive(ctx, createArchiveReq)
	s.Suite.NoError(err)
	s.Suite.NotNil(createArchiveRes)
	s.Suite.Equal(createArchiveRes.DoctorAvailabilityId, createArchiveReq.DoctorAvailabilityId)
	s.Suite.Equal(createArchiveRes.StartTime, createArchiveReq.StartTime)
	s.Suite.Equal(createArchiveRes.EndTime, createArchiveReq.EndTime)
	s.Suite.Equal(createArchiveRes.PatientProblem, createArchiveReq.PatientProblem)
	s.Suite.Equal(createArchiveRes.Status, createArchiveReq.Status)
	s.Suite.Equal(createArchiveRes.PaymentType, createArchiveReq.PaymentType)
	s.Suite.Equal(createArchiveRes.PaymentAmount, createArchiveReq.PaymentAmount)

	createNoteReq := &doctor_notes.CreatedDoctorNote{
		AppointmentId: createArchiveRes.Id,
		DoctorId:      uuid.New().String(),
		PatientId:     createPatientRes.Id,
		Prescription:  "Test Text",
	}

	createNoteRes, err := s.Repository.CreateDoctorNotes(ctx, createNoteReq)
	s.Suite.NoError(err)
	s.Suite.NotNil(createNoteRes)
	s.Suite.Equal(createNoteRes.AppointmentId, createNoteReq.AppointmentId)
	s.Suite.Equal(createNoteRes.DoctorId, createNoteReq.DoctorId)
	s.Suite.Equal(createNoteRes.PatientId, createNoteReq.PatientId)
	s.Suite.Equal(createNoteRes.Prescription, createNoteReq.Prescription)

	getRes, err := s.Repository.GetDoctorNotes(ctx, &doctor_notes.FieldValueReq{
		Field:        "id",
		Value:        strconv.Itoa(int(createNoteRes.Id)),
		DeleteStatus: true,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(getRes)
	s.Suite.Equal(getRes.Id, createNoteRes.Id)
	s.Suite.Equal(getRes.AppointmentId, createNoteRes.AppointmentId)
	s.Suite.Equal(getRes.DoctorId, createNoteRes.DoctorId)
	s.Suite.Equal(getRes.PatientId, createNoteRes.PatientId)
	s.Suite.Equal(getRes.Prescription, createNoteRes.Prescription)

	getAllRes, err := s.Repository.GetAllDoctorNotes(ctx, &doctor_notes.GetAllNotes{
		Page:         1,
		Limit:        5,
		DeleteStatus: true,
		OrderBy:      "id",
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(getAllRes)

	updateReq := &doctor_notes.UpdateDoctorNoteReq{
		Field:         "id",
		Value:         strconv.Itoa(int(getRes.Id)),
		AppointmentId: createArchiveRes.Id,
		DoctorId:      uuid.New().String(),
		PatientId:     createPatientRes.Id,
		Prescription:  "Update Text",
	}

	updateRes, err := s.Repository.UpdateDoctorNotes(ctx, updateReq)
	s.Suite.NoError(err)
	s.Suite.NotNil(updateRes)
	s.Suite.Equal(updateRes.Id, getRes.Id)
	s.Suite.Equal(updateRes.AppointmentId, updateReq.AppointmentId)
	s.Suite.Equal(updateRes.DoctorId, updateReq.DoctorId)
	s.Suite.Equal(updateRes.PatientId, updateReq.PatientId)
	s.Suite.Equal(updateRes.Prescription, updateReq.Prescription)

	//

	softDeleteNote, err := s.Repository.DeleteDoctorNotes(ctx, &doctor_notes.FieldValueReq{
		Field:        "id",
		Value:        strconv.Itoa(int(createNoteRes.Id)),
		DeleteStatus: false,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(softDeleteNote)
	s.Suite.Equal(softDeleteNote.Status, true)

	hardDeleteNote, err := s.Repository.DeleteDoctorNotes(ctx, &doctor_notes.FieldValueReq{
		Field:        "id",
		Value:        strconv.Itoa(int(createNoteRes.Id)),
		DeleteStatus: true,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(hardDeleteNote)
	s.Suite.Equal(hardDeleteNote.Status, true)

	hardDeleteRes, err := s.Patient.DeletePatient(ctx, &patients.FieldValueReq{
		Field:        "id",
		Value:        createPatientRes.Id,
		DeleteStatus: true,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(hardDeleteRes)
	s.Suite.Equal(hardDeleteRes.Status, true)

	hardDeleteArchiveRes, err := s.Archive.DeleteArchive(ctx, &archive.FieldValueReq{
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

func (s *DoctorNotesTestSite) TearDownSuite() {
	s.CleanUpFunc()
}

func TestDoctorNotesTestSuite(t *testing.T) {
	suite.Run(t, new(DoctorNotesTestSite))
}
