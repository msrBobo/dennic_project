package suit_tests

import (
	"booking_service/internal/entity/patients"
	repo "booking_service/internal/infrastructure/repository/postgresql"
	"booking_service/internal/pkg/config"
	db "booking_service/internal/pkg/postgres"
	"context"
	"github.com/google/uuid"
	"github.com/rickb777/date"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type BookingPatientsTestSite struct {
	suite.Suite
	Repository  *repo.BookingPatients
	CleanUpFunc func()
}

func (s *BookingPatientsTestSite) SetupSuite() {
	pgPool, _ := db.New(config.New())
	s.Repository = repo.NewBookingPatients(pgPool)
	s.CleanUpFunc = pgPool.Close
}

func (s *BookingPatientsTestSite) TestUserCRUD() {
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

	createRes, err := s.Repository.CreatePatient(ctx, patient)
	s.Suite.NoError(err)
	s.Suite.NotNil(createRes)
	s.Suite.Equal(createRes.Id, patient.Id)
	s.Suite.Equal(createRes.FirstName, patient.FirstName)
	s.Suite.Equal(createRes.LastName, patient.LastName)
	s.Suite.Equal(createRes.BirthDate, brithDate)
	s.Suite.Equal(createRes.Gender, patient.Gender)
	s.Suite.Equal(createRes.BloodGroup, patient.BloodGroup)
	s.Suite.Equal(createRes.PhoneNumber, patient.PhoneNumber)
	s.Suite.Equal(createRes.City, patient.City)
	s.Suite.Equal(createRes.Country, patient.Country)
	s.Suite.Equal(createRes.Address, patient.Address)
	s.Suite.Equal(createRes.PatientProblem, patient.PatientProblem)

	getRes, err := s.Repository.GetPatient(ctx, &patients.FieldValueReq{
		Field:        "id",
		Value:        createRes.Id,
		DeleteStatus: true,
	})

	s.Suite.NoError(err)
	s.Suite.NotNil(getRes)
	s.Suite.Equal(getRes.Id, patient.Id)
	s.Suite.Equal(getRes.FirstName, patient.FirstName)
	s.Suite.Equal(getRes.LastName, patient.LastName)
	s.Suite.Equal(getRes.BirthDate, brithDate)
	s.Suite.Equal(getRes.Gender, patient.Gender)
	s.Suite.Equal(getRes.BloodGroup, patient.BloodGroup)
	s.Suite.Equal(getRes.PhoneNumber, patient.PhoneNumber)
	s.Suite.Equal(getRes.City, patient.City)
	s.Suite.Equal(getRes.Country, patient.Country)
	s.Suite.Equal(getRes.Address, patient.Address)
	s.Suite.Equal(getRes.PatientProblem, patient.PatientProblem)

	getAllRes, err := s.Repository.GetAllPatiens(ctx, &patients.GetAllPatients{
		Page:         1,
		Limit:        10,
		DeleteStatus: true,
		OrderBy:      "last_name",
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(getAllRes)

	newDate, _ := date.AutoParse("2005-12-15")
	newPatient := &patients.UpdatePatient{
		Field:          "id",
		Value:          createRes.Id,
		FirstName:      "Sherzod",
		LastName:       "Erkinov",
		BirthDate:      newDate,
		Gender:         "male",
		BloodGroup:     "B-",
		City:           "Tashkent",
		Country:        "Namangan",
		Address:        "New Address",
		PatientProblem: "Yes Problem",
	}

	updateRes, err := s.Repository.UpdatePatient(ctx, newPatient)
	s.Suite.NoError(err)
	s.Suite.NotNil(updateRes)
	s.Suite.Equal(updateRes.FirstName, newPatient.FirstName)
	s.Suite.Equal(updateRes.LastName, newPatient.LastName)
	s.Suite.Equal(updateRes.BirthDate, newDate)
	s.Suite.Equal(updateRes.Gender, newPatient.Gender)
	s.Suite.Equal(updateRes.BloodGroup, newPatient.BloodGroup)
	s.Suite.Equal(updateRes.City, newPatient.City)
	s.Suite.Equal(updateRes.Country, newPatient.Country)
	s.Suite.Equal(updateRes.Address, newPatient.Address)
	s.Suite.Equal(updateRes.PatientProblem, newPatient.PatientProblem)

	upPhone, err := s.Repository.UpdatePhonePatient(ctx, &patients.UpdatePhoneNumber{
		Field:       "id",
		Value:       updateRes.Id,
		PhoneNumber: "+998996542345",
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(upPhone)
	s.Suite.Equal(upPhone.Status, true)

	softDeleteRes, err := s.Repository.DeletePatient(ctx, &patients.FieldValueReq{
		Field:        "id",
		Value:        updateRes.Id,
		DeleteStatus: false,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(softDeleteRes)
	s.Suite.Equal(softDeleteRes.Status, true)

	hardDeleteRes, err := s.Repository.DeletePatient(ctx, &patients.FieldValueReq{
		Field:        "id",
		Value:        updateRes.Id,
		DeleteStatus: true,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(hardDeleteRes)
	s.Suite.Equal(hardDeleteRes.Status, true)
}

func (s *BookingPatientsTestSite) TearDownSuite() {
	s.CleanUpFunc()
}

func TestBookingPatientsTestSuite(t *testing.T) {
	suite.Run(t, new(BookingPatientsTestSite))
}
