package entity

import "time"

type DoctorWorkingHours struct {
	Id         int32
	DoctorId   string
	DayOfWeek  string
	StartTime  string
	FinishTime string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
}

type GetRequest struct {
	Field     string
	Value     string
	IsActive  bool
	DayOfWeek string
}

type ListDoctorWorkingHours struct {
	DoctorWhs []DoctorWorkingHours
	Count     int32
}
