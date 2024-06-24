package entity

import (
	"time"
)

type GetReqStr struct {
	Field    string
	Value    string
	IsActive bool
}

type Status struct {
}

type ListDoctors struct {
	Doctors []Doctor
	Count   int64
}

type Doctor struct {
	Id              string
	Order           int32
	FirstName       string
	LastName        string
	ImageUrl        string
	Gender          string
	BirthDate       string
	PhoneNumber     string
	Email           string
	Address         string
	City            string
	Country         string
	Salary          float32
	Bio             string
	StartWorkDate   string
	EndWorkDate     string
	WorkYears       int32
	DepartmentId    string
	RoomNumber      int32
	Password        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       time.Time
	Specializations []DoctorSpec
}

type DoctorSpec struct {
	Id   string
	Name string
}

type DoctorAndDoctorHours struct {
	Id              string
	Order           int32
	FirstName       string
	LastName        string
	ImageUrl        string
	Gender          string
	BirthDate       string
	PhoneNumber     string
	Email           string
	Address         string
	City            string
	Country         string
	Salary          float32
	StartTime       string
	FinishTime      string
	DayOfWeek       string
	Bio             string
	StartWorkDate   string
	EndWorkDate     string
	WorkYears       int32
	DepartmentId    string
	RoomNumber      int32
	Password        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       time.Time
	Specializations []DoctorSpec
}

type ListDoctorsAndHours struct {
	Count   int64
	Doctors []DoctorAndDoctorHours
}

type GetReqStrDep struct {
	DepartmentId string
	IsActive     bool
	Page         int32
	Limit        int32
	Field        string
	Value        string
	OrderBy      string
}

type GetReqStrSpec struct {
	SpecializationId string
	IsActive         bool
	Page             int32
	Limit            int32
	Field            string
	Value            string
	OrderBy          string
}
