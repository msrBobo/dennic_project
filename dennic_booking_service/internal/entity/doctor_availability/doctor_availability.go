package doctor_availability

import (
	"github.com/rickb777/date"
	"time"
)

type DoctorAvailability struct {
	Id           int64
	DepartmentId string
	DoctorId     string
	DoctorDate   date.Date
	StartTime    time.Time
	EndTime      time.Time
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}

type DoctorTimeStatusReq struct {
	DepartmentID string
	DoctorID     string
	DoctorDate   string
	StartTime    string
}

type DoctorTimeFieldValueByIdReq struct {
	Date     string
	DoctorID string
}

type Hours struct {
	Hours []string
}

type DoctorTimeUpdateStatusReq struct {
	DepartmentID string
	DoctorID     string
	DoctorDate   string
	StartTime    string
	Status       string
}

type TimeStatus struct {
	Status string
}

type DoctorAvailabilityType struct {
	Count               int64
	DoctorAvailabilitys []*DoctorAvailability
}

type CreateDoctorAvailability struct {
	DepartmentId string
	DoctorId     string
	DoctorDate   date.Date
	StartTime    time.Time
	EndTime      time.Time
	Status       string
}

type UpdateDoctorAvailability struct {
	Field        string
	Value        string
	DepartmentId string
	DoctorId     string
	DoctorDate   date.Date
	StartTime    time.Time
	EndTime      time.Time
	Status       string
}

type FieldValueReq struct {
	Field        string
	Value        string
	DeleteStatus bool
}

type GetAllReq struct {
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
