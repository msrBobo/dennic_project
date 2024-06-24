package suit_tests

import (
	"Healthcare_Evrone/internal/entity"
	repo "Healthcare_Evrone/internal/infrastructure/repository/postgresql"
	"Healthcare_Evrone/internal/pkg/config"
	db "Healthcare_Evrone/internal/pkg/postgres"
	"context"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
	"time"
)

type DoctorTestSuite struct {
	suite.Suite
	CleanUpFunc          func()
	Repository           *repo.DocTor
	RepositoryDepartment *repo.DepartMent
	RepositoryDWH        *repo.Dwh
}

func (s *DoctorTestSuite) SetupTest() {
	pgPool, err := db.New(config.New())
	if err != nil {
		log.Fatal(err)
		return
	}
	s.Repository = repo.NewDoctorRepo(pgPool)
	s.RepositoryDepartment = repo.NewDepartmentRepo(pgPool)
	s.RepositoryDWH = repo.NewDoctorWorkingHoursRepo(pgPool)
	s.CleanUpFunc = pgPool.Close
}

func (s *DoctorTestSuite) TestDoctorCrud() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(2))
	defer cancel()

	department := &entity.Department{
		Id:          uuid.NewString(),
		Name:        "Test name",
		Description: "Test description",
		ImageUrl:    "Test imageUrl",
		FloorNumber: 1,
	}

	respDep, err := s.RepositoryDepartment.CreateDepartment(ctx, department)
	s.Suite.NoError(err)
	s.Suite.NotNil(respDep)
	s.Suite.Equal(respDep.Id, department.Id)
	s.Suite.Equal(respDep.Name, department.Name)
	s.Suite.Equal(respDep.Description, department.Description)
	s.Suite.Equal(respDep.ImageUrl, department.ImageUrl)
	s.Suite.Equal(respDep.FloorNumber, department.FloorNumber)

	doctor := &entity.Doctor{
		Id:            uuid.NewString(),
		FirstName:     "Test first name",
		LastName:      "Test last name",
		Gender:        "male",
		BirthDate:     "12-12-12",
		PhoneNumber:   "Testphonenumber",
		Email:         "Test email",
		Address:       "Test address",
		City:          "Test city",
		Country:       "Test country",
		Salary:        1.1,
		Bio:           "Test bio",
		StartWorkDate: "12-12-12",
		EndWorkDate:   "12-12-12",
		WorkYears:     3,
		DepartmentId:  department.Id,
		RoomNumber:    1,
		Password:      "Test password",
	}
	respDoctor, err := s.Repository.CreateDoctor(ctx, doctor)
	s.Suite.NoError(err)
	s.Suite.NotNil(respDoctor)
	s.Suite.NotNil(respDoctor.CreatedAt)
	s.Suite.Equal(respDoctor.Id, doctor.Id)
	s.Suite.Equal(respDoctor.FirstName, doctor.FirstName)
	s.Suite.Equal(respDoctor.LastName, doctor.LastName)
	s.Suite.Equal(respDoctor.Gender, doctor.Gender)
	s.Suite.Equal(respDoctor.BirthDate, doctor.BirthDate)
	s.Suite.Equal(respDoctor.PhoneNumber, doctor.PhoneNumber)
	s.Suite.Equal(respDoctor.Email, doctor.Email)
	s.Suite.Equal(respDoctor.Address, doctor.Address)
	s.Suite.Equal(respDoctor.City, doctor.City)
	s.Suite.Equal(respDoctor.Country, doctor.Country)
	s.Suite.Equal(respDoctor.Salary, doctor.Salary)
	s.Suite.Equal(respDoctor.Bio, doctor.Bio)
	s.Suite.Equal(respDoctor.StartWorkDate, doctor.StartWorkDate)
	s.Suite.Equal(respDoctor.EndWorkDate, doctor.EndWorkDate)
	s.Suite.Equal(respDoctor.WorkYears, doctor.WorkYears)
	s.Suite.Equal(respDoctor.DepartmentId, doctor.DepartmentId)
	s.Suite.Equal(respDoctor.RoomNumber, doctor.RoomNumber)
	s.Suite.Equal(respDoctor.Password, doctor.Password)

	dwh := &entity.DoctorWorkingHours{
		DoctorId:   doctor.Id,
		DayOfWeek:  "Monday",
		StartTime:  "12-12-12 12:12:12",
		FinishTime: "12-12-12 12:12:12",
	}
	respDoctorService, err := s.RepositoryDWH.CreateDoctorWorkingHours(ctx, dwh)
	s.Suite.NoError(err)
	s.Suite.NotNil(respDoctorService)
	s.Suite.NotNil(respDoctorService.CreatedAt)
	s.Suite.Equal(respDoctorService.DoctorId, dwh.DoctorId)
	s.Suite.Equal(respDoctorService.DayOfWeek, dwh.DayOfWeek)
	s.Suite.Equal(respDoctorService.StartTime, dwh.StartTime)
	s.Suite.Equal(respDoctorService.StartTime, dwh.StartTime)

	getDoctor, err := s.Repository.GetDoctorById(ctx, &entity.GetReqStr{
		Field:    "id",
		Value:    doctor.Id,
		IsActive: false,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(getDoctor)
	s.Suite.Equal(getDoctor.Id, doctor.Id)
	s.Suite.Equal(getDoctor.Id, doctor.Id)
	s.Suite.Equal(getDoctor.FirstName, doctor.FirstName)
	s.Suite.Equal(getDoctor.LastName, doctor.LastName)
	s.Suite.Equal(getDoctor.Gender, doctor.Gender)
	s.Suite.Equal(getDoctor.PhoneNumber, doctor.PhoneNumber)
	s.Suite.Equal(getDoctor.Email, doctor.Email)
	s.Suite.Equal(getDoctor.Address, doctor.Address)
	s.Suite.Equal(getDoctor.City, doctor.City)
	s.Suite.Equal(getDoctor.Country, doctor.Country)
	s.Suite.Equal(getDoctor.Salary, doctor.Salary)
	s.Suite.Equal(getDoctor.Bio, doctor.Bio)
	s.Suite.Equal(getDoctor.StartWorkDate, doctor.StartWorkDate)
	s.Suite.Equal(getDoctor.WorkYears, doctor.WorkYears)
	s.Suite.Equal(getDoctor.DepartmentId, doctor.DepartmentId)
	s.Suite.Equal(getDoctor.RoomNumber, doctor.RoomNumber)
	s.Suite.Equal(getDoctor.Password, doctor.Password)

	respAll, err := s.Repository.GetAllDoctors(ctx, &entity.GetAll{
		Page:     1,
		Limit:    10,
		Field:    "",
		Value:    "",
		OrderBy:  "",
		IsActive: false,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(respAll)

	newUpFirsName := "UpdateFirstname"
	newUpLastName := "UpdateLastname"
	newUpEmail := "Update Email"
	newUpAddress := "Update Address"
	newUpCity := "Update City"
	newUpCountry := "Update Country"
	newUpSalary := 3.1
	newUpBio := "Update Bio"

	updatedDoctor, err := s.Repository.UpdateDoctor(ctx, &entity.Doctor{
		Id:            doctor.Id,
		FirstName:     newUpFirsName,
		LastName:      newUpLastName,
		Gender:        "male",
		BirthDate:     "12-12-12",
		PhoneNumber:   "Testphonenumber",
		Email:         newUpEmail,
		Address:       newUpAddress,
		City:          newUpCity,
		Country:       newUpCountry,
		Salary:        float32(newUpSalary),
		Bio:           newUpBio,
		StartWorkDate: "12-12-12",
		EndWorkDate:   "12-12-12",
		WorkYears:     3,
		DepartmentId:  department.Id,
		RoomNumber:    1,
		Password:      "Test password",
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(updatedDoctor)
	s.Suite.Equal(updatedDoctor.Id, doctor.Id)
	s.Suite.Equal(updatedDoctor.Id, doctor.Id)
	s.Suite.Equal(updatedDoctor.FirstName, newUpFirsName)
	s.Suite.Equal(updatedDoctor.LastName, newUpLastName)
	s.Suite.Equal(updatedDoctor.Gender, doctor.Gender)
	s.Suite.Equal(updatedDoctor.BirthDate, doctor.BirthDate)
	s.Suite.Equal(updatedDoctor.PhoneNumber, doctor.PhoneNumber)
	s.Suite.Equal(updatedDoctor.Email, newUpEmail)
	s.Suite.Equal(updatedDoctor.Address, newUpAddress)
	s.Suite.Equal(updatedDoctor.City, newUpCity)
	s.Suite.Equal(updatedDoctor.Country, newUpCountry)
	s.Suite.Equal(updatedDoctor.Salary, float32(newUpSalary))
	s.Suite.Equal(updatedDoctor.Bio, newUpBio)
	s.Suite.Equal(updatedDoctor.StartWorkDate, doctor.StartWorkDate)
	s.Suite.Equal(updatedDoctor.WorkYears, doctor.WorkYears)
	s.Suite.Equal(updatedDoctor.DepartmentId, doctor.DepartmentId)
	s.Suite.Equal(updatedDoctor.RoomNumber, doctor.RoomNumber)
	s.Suite.Equal(updatedDoctor.Password, doctor.Password)

	resp, err := s.Repository.ListDoctorsByDepartmentId(ctx, &entity.GetReqStrDep{
		DepartmentId: doctor.DepartmentId,
		IsActive:     true,
		Page:         1,
		Limit:        10,
		Field:        "",
		Value:        "",
		OrderBy:      "",
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(resp)

	deleteDoctor, err := s.Repository.DeleteDoctor(ctx, &entity.GetReqStr{
		Field:    "id",
		Value:    doctor.Id,
		IsActive: true,
	})
	s.Suite.NotNil(deleteDoctor)
	s.Suite.NoError(err)

	deleteDep, err := s.RepositoryDepartment.DeleteDepartment(ctx, &entity.GetReqStr{
		Field:    "id",
		Value:    department.Id,
		IsActive: true,
	})
	s.Suite.NotNil(deleteDep)
	s.Suite.NoError(err)

	deleteDoctorWorkingHours, err := s.RepositoryDWH.DeleteDoctorWorkingHours(ctx, &entity.GetReqStr{
		Field:    "id",
		Value:    cast.ToString(dwh.Id),
		IsActive: true,
	})
	s.Suite.NotNil(deleteDoctorWorkingHours)
	s.Suite.NoError(err)
}

func (s *DoctorTestSuite) TearDownTest() {
	s.CleanUpFunc()
}

func TestDoctorTestSuite(t *testing.T) {
	suite.Run(t, new(DoctorTestSuite))
}
