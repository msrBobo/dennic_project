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

type DepartmentTestSuite struct {
	suite.Suite
	CleanUpFunc func()
	Repository  *repo.DepartMent
}

func (s *DepartmentTestSuite) SetupTest() {
	pgPool, err := db.New(config.New())
	if err != nil {
		log.Fatal(err)
		return
	}
	s.Repository = repo.NewDepartmentRepo(pgPool)
	s.CleanUpFunc = pgPool.Close
}

func (s *DepartmentTestSuite) TestDepartmentCrud() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(2))
	defer cancel()

	department := &entity.Department{
		Id:               uuid.NewString(),
		Name:             "Test name",
		Description:      "Test description",
		ImageUrl:         "Test imageUrl",
		FloorNumber:      1,
		ShortDescription: "Test shortDescription",
	}
	respDep, err := s.Repository.CreateDepartment(ctx, department)
	s.Suite.NoError(err)
	s.Suite.NotNil(respDep)
	s.Suite.Equal(respDep.Id, department.Id)
	s.Suite.Equal(respDep.Name, department.Name)
	s.Suite.Equal(respDep.Description, department.Description)
	s.Suite.Equal(respDep.ImageUrl, department.ImageUrl)
	s.Suite.Equal(respDep.FloorNumber, department.FloorNumber)
	s.Suite.Equal(respDep.ShortDescription, department.ShortDescription)

	getDepartment, err := s.Repository.GetDepartmentById(ctx, &entity.GetReqStr{
		Field:    "id",
		Value:    department.Id,
		IsActive: false,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(getDepartment)
	s.Suite.Equal(getDepartment.Id, department.Id)
	s.Suite.Equal(getDepartment.Name, department.Name)
	s.Suite.Equal(getDepartment.Description, department.Description)
	s.Suite.Equal(getDepartment.ImageUrl, department.ImageUrl)
	s.Suite.Equal(getDepartment.FloorNumber, department.FloorNumber)
	s.Suite.Equal(getDepartment.ShortDescription, department.ShortDescription)

	respAll, err := s.Repository.GetAllDepartments(ctx, &entity.GetAll{
		Page:     1,
		Limit:    10,
		Field:    "",
		Value:    "",
		OrderBy:  "",
		IsActive: false,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(respAll)

	newUpName := "Update Name"
	newUpDescription := "Update Description"
	newUpImageUrl := "Update image url"

	updatedDepartment, err := s.Repository.UpdateDepartment(ctx, &entity.Department{
		Id:               department.Id,
		Name:             newUpName,
		Description:      newUpDescription,
		ImageUrl:         newUpImageUrl,
		FloorNumber:      1,
		ShortDescription: department.ShortDescription,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(updatedDepartment)
	s.Suite.NotNil(updatedDepartment.UpdatedAt)
	s.Suite.Equal(updatedDepartment.Id, department.Id)
	s.Suite.Equal(newUpName, updatedDepartment.Name)
	s.Suite.Equal(newUpDescription, updatedDepartment.Description)
	s.Suite.Equal(newUpImageUrl, updatedDepartment.ImageUrl)
	s.Suite.Equal(department.ShortDescription, updatedDepartment.ShortDescription)

	deleteDep, err := s.Repository.DeleteDepartment(ctx, &entity.GetReqStr{
		Field:    "id",
		Value:    department.Id,
		IsActive: true,
	})
	s.Suite.NotNil(deleteDep)
	s.Suite.NoError(err)
}

func (s *DepartmentTestSuite) TearDownTest() {
	s.CleanUpFunc()
}

func TestDepartmentTestSuite(t *testing.T) {
	suite.Run(t, new(DepartmentTestSuite))
}
