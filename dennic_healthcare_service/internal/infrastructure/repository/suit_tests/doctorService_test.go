package suit_tests

import (
	"Healthcare_Evrone/internal/entity"
	repo "Healthcare_Evrone/internal/infrastructure/repository/postgresql"
	"Healthcare_Evrone/internal/pkg/config"
	db "Healthcare_Evrone/internal/pkg/postgres"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
	"time"
)

type DoctorServiceTestSuite struct {
	suite.Suite
	CleanUpFunc              func()
	Repository               *repo.Ds
	RepositoryDoctor         *repo.DocTor
	RepositorySpecialization *repo.Specialization
	RepositoryDepartment     *repo.DepartMent
}

func (s *DoctorServiceTestSuite) SetupTest() {
	pgPool, err := db.New(config.New())
	if err != nil {
		log.Fatal(err)
		return
	}
	s.Repository = repo.NewDoctorServicesRepo(pgPool)
	s.RepositoryDoctor = repo.NewDoctorRepo(pgPool)
	s.RepositorySpecialization = repo.NewSpecializationRepo(pgPool)
	s.RepositoryDepartment = repo.NewDepartmentRepo(pgPool)
	s.CleanUpFunc = pgPool.Close
}

func (s *DoctorServiceTestSuite) TestDoctorServiceCrud() {
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

	special := &entity.Specialization{
		ID:           uuid.NewString(),
		Name:         "Test Name",
		Description:  "Test Description",
		DepartmentId: department.Id,
	}

	respSpec, err := s.RepositorySpecialization.CreateSpecialization(ctx, special)
	s.Suite.NoError(err)
	s.Suite.NotNil(respSpec)
	s.Suite.NotNil(respSpec.CreatedAt)
	s.Suite.Equal(respSpec.ID, special.ID)
	s.Suite.Equal(respSpec.Name, special.Name)
	s.Suite.Equal(respSpec.Description, special.Description)
	s.Suite.Equal(respSpec.DepartmentId, special.DepartmentId)

	tm, err := time.Parse("15:04", "12:12")
	s.Suite.NoError(err)
	doctorServices := &entity.DoctorServices{
		Id:               uuid.NewString(),
		DoctorId:         doctor.Id,
		SpecializationId: special.ID,
		OnlinePrice:      1,
		OfflinePrice:     1,
		Name:             "Test Name",
		Duration:         tm,
	}
	respDoctorService, err := s.Repository.CreateDoctorServices(ctx, doctorServices)
	s.Suite.NoError(err)
	s.Suite.NotNil(respDoctorService)
	s.Suite.NotNil(respDoctorService.CreatedAt)
	s.Suite.Equal(respDoctorService.DoctorId, doctorServices.DoctorId)
	s.Suite.Equal(respDoctorService.Id, doctorServices.Id)
	s.Suite.Equal(respDoctorService.SpecializationId, doctorServices.SpecializationId)
	s.Suite.Equal(respDoctorService.OnlinePrice, doctorServices.OnlinePrice)
	s.Suite.Equal(respDoctorService.OfflinePrice, doctorServices.OfflinePrice)
	s.Suite.Equal(respDoctorService.Name, doctorServices.Name)
	s.Suite.Equal(respDoctorService.Duration, doctorServices.Duration)

	getDoctorService, err := s.Repository.GetDoctorServiceByID(ctx, &entity.GetReqStr{
		Field:    "id",
		Value:    doctorServices.Id,
		IsActive: false,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(getDoctorService)
	s.Suite.Equal(respDoctorService.DoctorId, doctorServices.DoctorId)
	s.Suite.Equal(respDoctorService.Id, doctorServices.Id)
	s.Suite.Equal(respDoctorService.SpecializationId, doctorServices.SpecializationId)
	s.Suite.Equal(respDoctorService.OnlinePrice, doctorServices.OnlinePrice)
	s.Suite.Equal(respDoctorService.OfflinePrice, doctorServices.OfflinePrice)
	s.Suite.Equal(respDoctorService.Name, doctorServices.Name)
	s.Suite.Equal(respDoctorService.Duration, doctorServices.Duration)

	respAll, err := s.Repository.GetAllDoctorServices(ctx, &entity.GetAll{
		Page:     1,
		Limit:    10,
		Field:    "",
		Value:    "",
		OrderBy:  "",
		IsActive: false,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(respAll)

	newUpOnlinePrice := 2
	newUpOfflinePrice := 2

	updatedDoctorService, err := s.Repository.UpdateDoctorServices(ctx, &entity.DoctorServices{
		Id:               doctorServices.Id,
		DoctorId:         doctorServices.DoctorId,
		SpecializationId: doctorServices.SpecializationId,
		OnlinePrice:      float32(newUpOnlinePrice),
		OfflinePrice:     float32(newUpOfflinePrice),
		Name:             doctorServices.Name,
		Duration:         doctorServices.Duration,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(updatedDoctorService)
	s.Suite.NotNil(updatedDoctorService.UpdatedAt)
	s.Suite.Equal(doctorServices.DoctorId, updatedDoctorService.DoctorId)
	s.Suite.Equal(doctorServices.SpecializationId, updatedDoctorService.SpecializationId)
	s.Suite.Equal(float32(newUpOnlinePrice), updatedDoctorService.OfflinePrice)
	s.Suite.Equal(float32(newUpOfflinePrice), updatedDoctorService.OfflinePrice)
	s.Suite.Equal(doctorServices.Name, updatedDoctorService.Name)
	s.Suite.Equal(doctorServices.Duration, updatedDoctorService.Duration)

	deleteDoctorService, err := s.Repository.DeleteDoctorService(ctx, &entity.GetReqStr{
		Field:    "id",
		Value:    doctorServices.Id,
		IsActive: true,
	})
	s.Suite.NotNil(deleteDoctorService)
	s.Suite.NoError(err)

	deleteSpec, err := s.RepositorySpecialization.DeleteSpecialization(ctx, &entity.GetReqStr{
		Field:    "id",
		Value:    special.ID,
		IsActive: true,
	})
	s.Suite.NotNil(deleteSpec)
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

func (s *DoctorServiceTestSuite) TearDownTest() {
	s.CleanUpFunc()
}

func TestDoctorServiceTestSuite(t *testing.T) {
	suite.Run(t, new(DoctorServiceTestSuite))
}
