package entity

import "time"

type Reasons struct {
	Id               string
	Name             string
	SpecializationId string
	ImageUrl         string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        time.Time
}

type ListReasons struct {
	Reasons []Reasons
	Count   int32
}

type GetReqStrReasons struct {
	Id       string
	IsActive bool
}

type StatusReasons struct {
	Status bool
}

type GetAllReas struct {
	Page     int32
	Limit    int32
	IsActive bool
	Field    string
	Value    string
	OrderBy  string
}
