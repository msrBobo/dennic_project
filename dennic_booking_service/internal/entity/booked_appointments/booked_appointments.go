package booked_appointments

import (
	"time"

	"github.com/rickb777/date"
)

type Appointment struct {
	Id              int64
	DepartmentId    string
	DoctorId        string
	PatientId       string
	ServiceId       string
	AppointmentDate date.Date
	AppointmentTime time.Time
	Duration        int64
	Key             string
	ExpiresAt       time.Time
	PatientProblem  string
	Status          string
	PaymentType     string
	PaymentAmount   float64
	UserId          string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       time.Time
}

type AppointmentsType struct {
	Count        int64
	Appointments []*Appointment
}

type CreateAppointment struct {
	DepartmentId    string
	DoctorId        string
	PatientId       string
	ServiceId       string
	AppointmentDate date.Date
	AppointmentTime time.Time
	Duration        int64
	Key             string
	ExpiresAt       time.Time
	PatientProblem  string
	Status          string
	PaymentType     string
	PaymentAmount   float64
	UserId          string
}

type UpdateAppointment struct {
	Field           string
	Value           string
	DepartmentId    string
	DoctorId        string
	PatientId       string
	ServiceId       string
	AppointmentDate date.Date
	AppointmentTime time.Time
	Duration        int64
	Key             string
	ExpiresAt       time.Time
	PatientProblem  string
	Status          string
	PaymentType     string
	PaymentAmount   float64
	UserId          string
}

type GetAllAppointment struct {
	Page         uint64
	Limit        uint64
	DeleteStatus bool
	Field        string
	Value        string
	OrderBy      string
}

type GetFilteredRequest struct {
	Page         uint64
	Limit        uint64
	DeleteStatus bool
	Field        string
	Value        string
	OrderBy      string
	Status       string
}

type FieldValueReq struct {
	Field        string
	Value        string
	DeleteStatus bool
}

type StatusRes struct {
	Status bool
}
