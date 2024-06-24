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

type ReasonsTestSuite struct {
	suite.Suite
	CleanUpFunc              func()
	Repository               *repo.Reasons
	RepositorySpecialization *repo.Specialization
	RepositoryDepartment     *repo.DepartMent
}

func (s *ReasonsTestSuite) SetupTest() {
	pgPool, err := db.New(config.New())
	if err != nil {
		log.Fatal(err)
		return
	}
	s.Repository = repo.NewReasonsRepo(pgPool)
	s.RepositorySpecialization = repo.NewSpecializationRepo(pgPool)
	s.RepositoryDepartment = repo.NewDepartmentRepo(pgPool)
	s.CleanUpFunc = pgPool.Close
}

func (s *ReasonsTestSuite) TestReasonsCrud() {
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

	respSpec, err := s.RepositorySpecialization.CreateSpecialization(ctx, special)
	s.Suite.NoError(err)
	s.Suite.NotNil(respSpec)
	s.Suite.NotNil(respSpec.CreatedAt)
	s.Suite.Equal(respSpec.ID, special.ID)
	s.Suite.Equal(respSpec.Name, special.Name)
	s.Suite.Equal(respSpec.Description, special.Description)
	s.Suite.Equal(respSpec.DepartmentId, special.DepartmentId)

	reason := &entity.Reasons{
		Id:               uuid.NewString(),
		Name:             "Test Name",
		SpecializationId: special.ID,
	}

	respReas, err := s.Repository.CreateReasons(ctx, reason)
	s.Suite.NoError(err)
	s.Suite.NotNil(respReas)
	s.Suite.Equal(respReas.Id, reason.Id)
	s.Suite.Equal(respReas.Name, reason.Name)
	s.Suite.Equal(respReas.SpecializationId, reason.SpecializationId)
	s.Suite.NotNil(respReas.CreatedAt)

	getReas, err := s.Repository.GetReasonsById(ctx, &entity.GetReqStr{
		Field:    "id",
		Value:    reason.Id,
		IsActive: true,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(getReas)
	s.Suite.Equal(getReas.Id, reason.Id)
	s.Suite.Equal(getReas.Name, reason.Name)
	s.Suite.Equal(getReas.SpecializationId, reason.SpecializationId)
	s.Suite.Equal(getReas.CreatedAt, reason.CreatedAt)

	getAllReason, err := s.Repository.GetAllReasons(ctx, &entity.GetAllReas{
		Page:     1,
		Limit:    10,
		IsActive: true,
		Field:    "",
		Value:    "",
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(getAllReason)

	updateReason, err := s.Repository.UpdateReasons(ctx, &entity.Reasons{
		Name:             "Updated Name",
		SpecializationId: special.ID,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(updateReason)
	s.Suite.NotNil(updateReason.UpdatedAt)

	deleteReas, err := s.Repository.DeleteReasons(ctx, &entity.GetReqStr{
		Field:    "id",
		Value:    reason.Id,
		IsActive: true,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(deleteReas)

	deleteSpec, err := s.RepositorySpecialization.DeleteSpecialization(ctx, &entity.GetReqStr{
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

func (s *ReasonsTestSuite) TearDownTest() {
	s.CleanUpFunc()
}

func TestReasonsTestSuite(t *testing.T) {
	suite.Run(t, new(ReasonsTestSuite))
}
