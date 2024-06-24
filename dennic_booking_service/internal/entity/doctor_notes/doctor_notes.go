package doctor_notes

import (
	"time"
)

type DoctorNote struct {
	Id            int64
	AppointmentId int64
	DoctorId      string
	PatientId     string
	Prescription  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
}

type DoctorNotesType struct {
	Count       int64
	DoctorNotes []*DoctorNote
}

type CreatedDoctorNote struct {
	AppointmentId int64
	DoctorId      string
	PatientId     string
	Prescription  string
}

type UpdateDoctorNoteReq struct {
	Field         string
	Value         string
	AppointmentId int64
	DoctorId      string
	PatientId     string
	Prescription  string
}

type GetAllNotes struct {
	Page         uint64
	Limit        uint64
	DeleteStatus bool
	Field        string
	Value        string
	OrderBy      string
}

type FieldValueReq struct {
	Field        string
	Value        string
	DeleteStatus bool
}

type StatusRes struct {
	Status bool
}
