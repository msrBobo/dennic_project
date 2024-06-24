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

type DoctorWorkingHoursTestSuite struct {
	suite.Suite
	CleanUpFunc          func()
	Repository           *repo.Dwh
	RepositoryDoctor     *repo.DocTor
	RepositoryDepartment *repo.DepartMent
}

func (s *DoctorWorkingHoursTestSuite) SetupTest() {
	pgPool, err := db.New(config.New())
	if err != nil {
		log.Fatal(err)
		return
	}
	s.Repository = repo.NewDoctorWorkingHoursRepo(pgPool)
	s.RepositoryDoctor = repo.NewDoctorRepo(pgPool)
	s.RepositoryDepartment = repo.NewDepartmentRepo(pgPool)
	s.CleanUpFunc = pgPool.Close
}

func (s *DoctorWorkingHoursTestSuite) TestDoctorWorkingHoursCrud() {
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
	respDoctor, err := s.RepositoryDoctor.CreateDoctor(ctx, doctor)
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
	respDoctorService, err := s.Repository.CreateDoctorWorkingHours(ctx, dwh)
	s.Suite.NoError(err)
	s.Suite.NotNil(respDoctorService)
	s.Suite.NotNil(respDoctorService.CreatedAt)
	s.Suite.Equal(respDoctorService.DoctorId, dwh.DoctorId)
	s.Suite.Equal(respDoctorService.DayOfWeek, dwh.DayOfWeek)
	s.Suite.Equal(respDoctorService.StartTime, dwh.StartTime)
	s.Suite.Equal(respDoctorService.StartTime, dwh.StartTime)

	getDoctorWorkingHours, err := s.Repository.GetDoctorWorkingHoursById(ctx, &entity.GetRequest{
		Field:    "id",
		Value:    cast.ToString(respDoctorService.Id),
		IsActive: false,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(getDoctorWorkingHours)
	s.Suite.Equal(respDoctorService.DoctorId, dwh.DoctorId)
	s.Suite.Equal(respDoctorService.DayOfWeek, dwh.DayOfWeek)
	s.Suite.Equal(respDoctorService.StartTime, dwh.StartTime)
	s.Suite.Equal(respDoctorService.StartTime, dwh.StartTime)

	respAll, err := s.Repository.GetAllDoctorWorkingHours(ctx, &entity.GetAll{
		Page:     1,
		Limit:    10,
		Field:    "",
		Value:    "",
		OrderBy:  "",
		IsActive: false,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(respAll)

	newUpDayOfWeek := "Sunday"

	updatedDoctorWorkingHours, err := s.Repository.UpdateDoctorWorkingHours(ctx, &entity.DoctorWorkingHours{
		Id:         dwh.Id,
		DoctorId:   dwh.DoctorId,
		DayOfWeek:  newUpDayOfWeek,
		StartTime:  "12-12-12 12:12:12",
		FinishTime: "12-12-12 12:12:12",
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(updatedDoctorWorkingHours)
	s.Suite.NotNil(updatedDoctorWorkingHours.UpdatedAt)
	s.Suite.Equal(respDoctorService.DoctorId, updatedDoctorWorkingHours.DoctorId)
	s.Suite.Equal(newUpDayOfWeek, updatedDoctorWorkingHours.DayOfWeek)
	s.Suite.Equal(respDoctorService.StartTime, updatedDoctorWorkingHours.StartTime)
	s.Suite.Equal(respDoctorService.StartTime, updatedDoctorWorkingHours.StartTime)

	deleteDoctorWorkingHours, err := s.Repository.DeleteDoctorWorkingHours(ctx, &entity.GetReqStr{
		Field:    "id",
		Value:    cast.ToString(dwh.Id),
		IsActive: true,
	})
	s.Suite.NotNil(deleteDoctorWorkingHours)
	s.Suite.NoError(err)

	deleteDoctor, err := s.RepositoryDoctor.DeleteDoctor(ctx, &entity.GetReqStr{
		Field:    "id",
		Value:    doctor.Id,
		IsActive: true,
	})
	s.Suite.NotNil(deleteDoctor)
	s.Suite.NoError(err)

	deleteDep, err := s.RepositoryDepartment.DeleteDepartment(ctx, &entity.GetReqStr{
		Field:    "id",
		Value:    department.Id,
		IsActive: false,
	})
	s.Suite.NotNil(deleteDep)
	s.Suite.NoError(err)

}

func (s *DoctorWorkingHoursTestSuite) TearDownTest() {
	s.CleanUpFunc()
}

func TestDoctorWorkingHoursTestSuite(t *testing.T) {
	suite.Run(t, new(DoctorWorkingHoursTestSuite))
}
