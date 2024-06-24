package postgresql

import (
	"context"
	"dennic_user_service/internal/entity"
	"dennic_user_service/internal/pkg/config"
	db "dennic_user_service/internal/pkg/postgres"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/stretchr/testify/suite"
)

type AdminReposisitoryTestSuite struct {
	suite.Suite
	repo        *adminRepo
	CleanUpFunc func()
}

func (s *AdminReposisitoryTestSuite) SetupSuite() {
	pgPool, _ := db.New(config.New())
	s.repo = NewAdminRepo(pgPool)
	s.CleanUpFunc = pgPool.Close
}

// test func
func (s *AdminReposisitoryTestSuite) TestAdminCRUD() {

	ctx := context.Background()

	// struct for create admin
	admin := entity.Admin{
		AdminOrder:    777,
		Role:          "admin",
		FirstName:     "testdata",
		LastName:      "testdata",
		BirthDate:     "2000-08-30",
		PhoneNumber:   "testdata",
		Email:         "testdata",
		Password:      "testdata",
		Gender:        "male",
		Salary:        777.7,
		Biography:     "testdata",
		StartWorkYear: "2000-08-30",
		EndWorkYear:   "2000-08-30",
		WorkYears:     777,
		RefreshToken:  "testdata",
		CreatedAt:     time.Now().UTC(),
	}
	// uuid generating
	admin.Id = uuid.New().String()

	updAdmin := entity.Admin{
		Id:            admin.Id,
		AdminOrder:    888,
		Role:          "superadmin",
		FirstName:     "updtestdata",
		LastName:      "updtestdata",
		BirthDate:     "2000-07-20",
		PhoneNumber:   "updtestdata",
		Email:         "updtestdata",
		Password:      "updtestdata",
		Gender:        "male",
		Salary:        888,
		Biography:     "updtestdata",
		StartWorkYear: "2000-08-30",
		EndWorkYear:   "2000-08-30",
		WorkYears:     888,
		RefreshToken:  "updtestdata",
		UpdatedAt:     time.Now(),
	}
	_ = updAdmin.UpdatedAt
	// check create admin method
	err := s.repo.Create(ctx, &admin)
	s.Suite.NoError(err)
	req := entity.FieldValueReq{
		Field: "id",
		Value: admin.Id,
		DeleteStatus: false,
	}

	// check get admin method
	getAdmin, err := s.repo.Get(ctx, &req)
	s.Suite.NoError(err)
	s.Suite.NotNil(getAdmin)
	s.Suite.Equal(getAdmin.Id, admin.Id)
	s.Suite.Equal(getAdmin.FirstName, admin.FirstName)
	s.Suite.Equal(getAdmin.PhoneNumber, admin.PhoneNumber)

	// check update admin method
	err = s.repo.Update(ctx, &updAdmin)
	s.Suite.NoError(err)
	updGetAdmin, err := s.repo.Get(ctx, &req)
	s.Suite.NoError(err)
	s.Suite.NotNil(updGetAdmin)
	s.Suite.Equal(updGetAdmin.Id, updAdmin.Id)
	s.Suite.Equal(updGetAdmin.FirstName, updAdmin.FirstName)

	// // check getAllAdmins method
	getAllReq := entity.GetAllReq{
		Page:         1,
		Limit:        5,
		DeleteStatus: false,
		Field:        "first_name",
		Value:        updAdmin.FirstName,
		OrderBy:      "first_name",
	}
	getAllAdmins, err := s.repo.List(ctx, &getAllReq)
	s.Suite.NoError(err)
	s.Suite.NotNil(getAllAdmins)

	// // check CheckField admin method
	CheckFieldReq := entity.CheckFieldReq{
		Field:    "id",
		Value:    admin.Id,
	}
	result, err := s.repo.CheckField(ctx, &CheckFieldReq)
	s.Suite.NoError(err)
	s.Suite.NotNil(updGetAdmin)
	s.Suite.Equal(result.Status, true)

	// check ChangePassword user method for PhoneNumber change
	change_password_req := entity.ChangeAdminPasswordReq{
		PhoneNumber: admin.PhoneNumber,
		Password:    "new_password",
	}

	resp_change_password, err := s.repo.ChangePassword(ctx, &change_password_req)
	s.Suite.NoError(err)
	s.Suite.NotNil(resp_change_password)
	s.Suite.Equal(resp_change_password.Status, true)

	// check ChangePassword user method for Email change
	change_password_req_2 := entity.ChangeAdminPasswordReq{
		Email:    admin.Email,
		Password: "new_password",
	}
	resp_change_password_2, err := s.repo.ChangePassword(ctx, &change_password_req_2)
	s.Suite.NoError(err)
	s.Suite.NotNil(resp_change_password_2)
	s.Suite.Equal(resp_change_password_2.Status, true)

	// // check UpdateRefreshToken admin method
	req_update_refresh_token := entity.UpdateRefreshTokenReq{
		Id:           updAdmin.Id,
		RefreshToken: "new_refresh_token",
	}
	resp_update_refresh_token, err := s.repo.UpdateRefreshToken(ctx, &req_update_refresh_token)
	s.Suite.NoError(err)
	s.Suite.NotNil(resp_update_refresh_token)
	s.Suite.Equal(resp_update_refresh_token.Status, true)

	// // check delete admin method
	DeleteAdminReq := entity.FieldValueReq{
		Field:        "id",
		Value:        admin.Id,
		DeleteStatus: false,
	}
	status, err := s.repo.Delete(ctx, &DeleteAdminReq)
	s.Suite.NoError(err)
	s.Suite.Equal(status.Status, true)
}

func TestExampleAdminTestSuite(t *testing.T) {
	suite.Run(t, new(AdminReposisitoryTestSuite))
}
