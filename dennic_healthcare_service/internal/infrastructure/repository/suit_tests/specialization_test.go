package suit_tests

import (
	"Healthcare_Evrone/internal/entity"
	repo "Healthcare_Evrone/internal/infrastructure/repository/postgresql"
	"Healthcare_Evrone/internal/pkg/config"
	db "Healthcare_Evrone/internal/pkg/postgres"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
	"time"
)

type SpecializationTestSuite struct {
	suite.Suite
	CleanUpFunc          func()
	Repository           *repo.Specialization
	RepositoryDepartment *repo.DepartMent
}

func (s *SpecializationTestSuite) SetupTest() {
	pgPool, err := db.New(config.New())
	if err != nil {
		log.Fatal(err)
		return
	}
	s.Repository = repo.NewSpecializationRepo(pgPool)
	s.RepositoryDepartment = repo.NewDepartmentRepo(pgPool)
	s.CleanUpFunc = pgPool.Close
}

func (s *SpecializationTestSuite) TestSpecializationCrud() {
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

	special := &entity.Specialization{
		ID:           uuid.NewString(),
		Name:         "Test Name",
		Description:  "Test Description",
		DepartmentId: department.Id,
	}

	respSpec, err := s.Repository.CreateSpecialization(ctx, special)
	s.Suite.NoError(err)
	s.Suite.NotNil(respSpec)
	s.Suite.NotNil(respSpec.CreatedAt)
	s.Suite.Equal(respSpec.ID, special.ID)
	s.Suite.Equal(respSpec.Name, special.Name)
	s.Suite.Equal(respSpec.Description, special.Description)
	s.Suite.Equal(respSpec.DepartmentId, special.DepartmentId)

	getSpecialization, err := s.Repository.GetSpecializationById(ctx, &entity.GetReqStr{
		Field:    "id",
		Value:    special.ID,
		IsActive: false,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(getSpecialization)
	s.Suite.Equal(getSpecialization.ID, special.ID)
	s.Suite.Equal(getSpecialization.Name, special.Name)
	s.Suite.Equal(getSpecialization.Description, special.Description)
	s.Suite.Equal(getSpecialization.DepartmentId, special.DepartmentId)

	respAll, err := s.Repository.GetAllSpecializations(ctx, &entity.GetAllSpecializations{
		Page:         1,
		Limit:        10,
		Field:        "",
		Value:        "",
		OrderBy:      "",
		DepartmentId: "",
		IsActive:     false,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(respAll)

	newUpName := "Update Name"
	newUpDescription := "Update Description"
	fmt.Println(special.ID)
	updatedSpecialization, err := s.Repository.UpdateSpecialization(ctx, &entity.Specialization{
		ID:           special.ID,
		Name:         newUpName,
		Description:  newUpDescription,
		DepartmentId: special.DepartmentId,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(updatedSpecialization)
	s.Suite.NotNil(updatedSpecialization.UpdatedAt)
	s.Suite.Equal(updatedSpecialization.ID, special.ID)
	s.Suite.Equal(newUpName, updatedSpecialization.Name)
	s.Suite.Equal(newUpDescription, updatedSpecialization.Description)

	deleteSpec, err := s.Repository.DeleteSpecialization(ctx, &entity.GetReqStr{
		Field:    "id",
		Value:    special.ID,
		IsActive: true,
	})
	s.Suite.NotNil(deleteSpec)
	s.Suite.NoError(err)

	deleteDep, err := s.RepositoryDepartment.DeleteDepartment(ctx, &entity.GetReqStr{
		Field:    "id",
		Value:    department.Id,
		IsActive: true,
	})
	s.Suite.NotNil(deleteDep)
	s.Suite.NoError(err)
}

func (s *SpecializationTestSuite) TearDownTest() {
	s.CleanUpFunc()
}

func TestSpecializationTestSuite(t *testing.T) {
	suite.Run(t, new(SpecializationTestSuite))
}
