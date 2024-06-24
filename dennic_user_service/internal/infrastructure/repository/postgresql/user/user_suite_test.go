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

type UserReposisitoryTestSuite struct {
	suite.Suite
	repo      *userRepo
	CleanUpFunc        func()
}

func (s *UserReposisitoryTestSuite) SetupSuite() {
	pgPool, _ := db.New(config.New())
	s.repo = NewUserRepo(pgPool)
	s.CleanUpFunc = pgPool.Close
}

// test func
func (s *UserReposisitoryTestSuite) TestUserCRUD() {

	ctx := context.Background()

	// struct for create user
	user := entity.User{
		UserOrder:    101,
		FirstName:    "firstname",
		LastName:     "lastname",
		BirthDate:    "2000-08-30",
		PhoneNumber:  "+998994767316",
		Password:     "testpassword",
		Gender:       "male",
		RefreshToken: "testrefreshtoken",
		CreatedAt:    time.Now().UTC(),
	}
	// uuid generating
	user.Id = uuid.New().String()

	updUser := entity.User{
		Id:          user.Id,
		FirstName:   "updfirstname",
		LastName:    "updlastname",
		BirthDate:   "2000-07-20",
		PhoneNumber: "+998934767316",
		Password:    "updtestpassword",
		Gender:      "male",
		UpdatedAt:   time.Now(),
	}

	// check create user method
	err := s.repo.Create(ctx, &user)
	s.Suite.NoError(err)
	req := entity.FieldValueReq{
		Field: "id",
		Value: user.Id,
		DeleteStatus: false,
	}

	// check get user method
	getUser, err := s.repo.Get(ctx, &req)
	s.Suite.NoError(err)
	s.Suite.NotNil(getUser)
	s.Suite.Equal(getUser.Id, user.Id)
	s.Suite.Equal(getUser.FirstName, user.FirstName)
	s.Suite.Equal(getUser.PhoneNumber, user.PhoneNumber)
	// check update user method
	err = s.repo.Update(ctx, &updUser)
	s.Suite.NoError(err)
	updGetUser, err := s.repo.Get(ctx, &req)
	s.Suite.NoError(err)
	s.Suite.NotNil(updGetUser)
	s.Suite.Equal(updGetUser.Id, updUser.Id)
	s.Suite.Equal(updGetUser.FirstName, updUser.FirstName)

	// check getAllUsers method
	getAllReq := entity.GetAllReq{
		Page:         1,
		Limit:        5,
		DeleteStatus: false,
		Field:        "first_name",
		Value:        updUser.FirstName,
		OrderBy:      "first_name",
	}
	getAllAdmins, err := s.repo.List(ctx, &getAllReq)
	s.Suite.NoError(err)
	s.Suite.NotNil(getAllAdmins)

	// check CheckField user method
	CheckFieldReq := entity.CheckFieldReq{
		Value:    user.Id,
		Field:    "id",
	}
	result, err := s.repo.CheckField(ctx, &CheckFieldReq)
	s.Suite.NoError(err)
	s.Suite.NotNil(result)
	s.Suite.Equal(result.Status, true)


	// check ChangePassword user method
	change_password_req := entity.ChangeUserPasswordReq{
		PhoneNumber: user.PhoneNumber,
		Password:    "new_password",
	}
	resp_change_password, err := s.repo.ChangePassword(ctx, &change_password_req)
	s.Suite.NoError(err)
	s.Suite.NotNil(resp_change_password)
	s.Suite.Equal(resp_change_password.Status, true)

	// check UpdateRefreshToken user method
	req_update_refresh_token := entity.UpdateRefreshTokenReq{
		Id:           updUser.Id,
		RefreshToken: "new_refresh_token",
	}
	resp_update_refresh_token, err := s.repo.UpdateRefreshToken(ctx, &req_update_refresh_token)
	s.Suite.NoError(err)
	s.Suite.NotNil(resp_update_refresh_token)
	s.Suite.Equal(resp_update_refresh_token.Status, true)

	//check delete user method
	DeleteAdminReq := entity.FieldValueReq{
		Field:        "id",
		Value:        user.Id,
		DeleteStatus: false,
	}
	status, err := s.repo.Delete(ctx, &DeleteAdminReq)
	s.Suite.NoError(err)
	s.Suite.Equal(status.Status, true)

}

func TestExampleUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserReposisitoryTestSuite))
}
