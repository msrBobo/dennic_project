package model_admin_service

import (
	"net/mail"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
)

type CreateAdminReq struct {
	Id            string  `json:"-"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	Role          string  `json:"role"`
	BrithDate     string  `json:"birth_date"`
	PhoneNumber   string  `json:"phone_number"`
	Email         string  `json:"email"`
	Password      string  `json:"password"`
	Gender        string  `json:"gender"`
	Salary        float32 `json:"salary"`
	Biography     string  `json:"biography"`
	StartWorkYear string  `json:"start_work_year"`
	EndWorkYear   string  `json:"-"`
	WorkYears     uint64  `json:"work_years"`
	RefreshToken  string  `json:"-"`
	ImageUrl      string  `json:"-"`
}

type CreateAdminResp struct {
	Id            string  `json:"id"`
	AdminOrder    int64   `json:"admin_order"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	Role          string  `json:"role"`
	BrithDate     string  `json:"birth_date"`
	PhoneNumber   string  `json:"phone_number"`
	Email         string  `json:"email"`
	Password      string  `json:"-"`
	Gender        string  `json:"gender"`
	Salary        float32 `json:"salary"`
	Biography     string  `json:"biography"`
	StartWorkYear string  `json:"start_work_year"`
	EndWorkYear   string  `json:"-"`
	WorkYears     uint64  `json:"work_years"`
	ImageUrl      string  `json:"-"`
	CreatedAt     string  `json:"created_at"`
	AccessToken   string  `json:"access_token"`
	ReflashToken  string  `json:"reflash_token"`
}

type GetAdminResp struct {
	Id            string  `json:"id"`
	AdminOrder    uint64  `json:"user_order"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	Role          string  `json:"role"`
	BrithDate     string  `json:"birth_date"`
	PhoneNumber   string  `json:"phone_number"`
	Email         string  `json:"email"`
	Password      string  `json:"-"`
	Gender        string  `json:"gender"`
	Salary        float32 `json:"salary"`
	Biography     string  `json:"biography"`
	StartWorkYear string  `json:"start_work_year"`
	EndWorkYear   string  `json:"end_work_year"`
	WorkYears     uint64  `json:"work_years"`
	RefreshToken  string  `json:"-"`
	ImageUrl      string  `json:"image_url"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	DeletedAt     string  `json:"deleted_at"`
}

type ListAdminrResp struct {
	Count  uint64         `json:"count"`
	Admins []GetAdminResp `json:"admins"`
}

type UpdAdminReq struct {
	Id            string  `json:"id"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	BrithDate     string  `json:"birth_date"`
	Gender        string  `json:"gender"`
	Salary        float32 `json:"salary"`
	Biography     string  `json:"biography"`
	StartWorkYear string  `json:"start_work_year"`
	EndWorkYear   string  `json:"end_work_year"`
	WorkYears     uint64  `json:"work_years"`
	ImageUrl      string  `json:"image_url"`
}

type UpdAdminResp struct {
	Id            string  `json:"id"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	BrithDate     string  `json:"birth_date"`
	Gender        string  `json:"gender"`
	Salary        float32 `json:"salary"`
	Biography     string  `json:"biography"`
	StartWorkYear string  `json:"start_work_year"`
	EndWorkYear   string  `json:"end_work_year"`
	WorkYears     uint64  `json:"work_years"`
	ImageUrl      string  `json:"image_url"`
	UpdatedAt     string  `json:"updated_at"`
}

type GetAllReq struct {
	Page    uint64 `json:"page"`
	Limit   uint64 `json:"limit"`
	Field   string `json:"field"`
	Value   string `json:"value"`
	OrderBy string `json:"order_by"`
}

type GetAdminReq struct {
	Field    string `json:"field"`
	Value    string `json:"value"`
	IsActive bool   `json:"is_active"`
}

type DeleteAdminReq struct {
	Field        string `json:"field"`
	Value        string `json:"value"`
	DeleteStatus bool   `json:"delete_status"`
}

type CheckAdminFieldResp struct {
	Status bool `json:"status"`
}

type Admin struct {
	Id            string
	AdminOrder    int64
	Role          string
	FirstName     string
	LastName      string
	BirthDate     string
	PhoneNumber   string
	Email         string
	Password      string
	Gender        string
	Salary        float32
	Biography     string
	StartWorkYear string
	EndWorkYear   string
	WorkYears     uint64
	RefreshToken  string
	ImageUrl      string
	Count         int64
	CreatedAt     string
	UpdatedAt     string
	DeletedAt     string
}

// User info Validate
func (u *CreateAdminReq) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.PhoneNumber, validation.Required, validation.Length(13, 13), validation.Match(regexp.MustCompile("^\\+[0-9]")).Error("Phone number is not valid")),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 32), validation.Match(regexp.MustCompile("^[a-zA-Z0-9!@#$%^&*()-_=+]")).Error("Password is not valid")),
		validation.Field(&u.FirstName, validation.Required, validation.Length(3, 50), validation.Match(regexp.MustCompile("^[A-Z][a-zA-Z']*([\\\\s-][A-Z][a-zA-Z']*)*$")).Error("First name is not valid")),
		validation.Field(&u.LastName, validation.Required, validation.Length(3, 50), validation.Match(regexp.MustCompile("^[A-Z][a-zA-Z']*([\\\\s-][A-Z][a-zA-Z']*)*$")).Error("Last name is not valid")),
	)
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
