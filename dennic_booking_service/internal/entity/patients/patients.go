package patients

import (
	"github.com/rickb777/date"
	"time"
)

// datient

type Patient struct {
	Id             string
	FirstName      string
	LastName       string
	BirthDate      date.Date
	Gender         string
	BloodGroup     string
	PhoneNumber    string
	City           string
	Country        string
	Address        string
	PatientProblem string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
}

type PatientsType struct {
	Count    int64
	Patients []*Patient
}

type CreatedPatient struct {
	Id             string
	FirstName      string
	LastName       string
	BirthDate      date.Date
	Gender         string
	BloodGroup     string
	PhoneNumber    string
	City           string
	Country        string
	Address        string
	PatientProblem string
}

type UpdatePatient struct {
	Field          string
	Value          string
	FirstName      string
	LastName       string
	BirthDate      date.Date
	Gender         string
	BloodGroup     string
	City           string
	Country        string
	Address        string
	PatientProblem string
}

type UpdatePhoneNumber struct {
	Field       string
	Value       string
	PhoneNumber string
}

type FieldValueReq struct {
	Field        string
	Value        string
	DeleteStatus bool
}

type GetAllPatients struct {
	Page         uint64
	Limit        uint64
	DeleteStatus bool
	Field        string
	Value        string
	OrderBy      string
}

type StatusRes struct {
	Status bool
}
