package entity

import "time"

type DoctorServices struct {
	Id               string
	Order            int32
	DoctorId         string
	SpecializationId string
	OnlinePrice      float32
	OfflinePrice     float32
	Name             string
	Duration         time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        time.Time
}

type ListDoctorServices struct {
	DoctorServices []DoctorServices
	Count          int64
}
