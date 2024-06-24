package archive

import "time"

type Archive struct {
	Id                   int64
	DoctorAvailabilityId int64
	StartTime            time.Time
	EndTime              time.Time
	PatientProblem       string
	Status               string
	PaymentType          string
	PaymentAmount        float64
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            time.Time
}

type ArchivesType struct {
	Count    int64
	Archives []*Archive
}

type CreatedArchive struct {
	DoctorAvailabilityId int64
	StartTime            time.Time
	EndTime              time.Time
	PatientProblem       string
	Status               string
	PaymentType          string
	PaymentAmount        float64
}

type UpdateArchive struct {
	Field                string
	Value                string
	DoctorAvailabilityId int64
	StartTime            time.Time
	EndTime              time.Time
	PatientProblem       string
	Status               string
	PaymentType          string
	PaymentAmount        float64
}

type GetAllArchives struct {
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
