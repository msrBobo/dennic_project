package entity

import "time"

type Specialization struct {
	ID           string
	Order        int32
	Name         string
	Description  string
	DepartmentId string
	ImageUrl     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}

type ListSpecializations struct {
	Specializations []Specialization
	Count           int32
}

type GetAllSpecializations struct {
	Page         int64
	Limit        int64
	Field        string
	Value        string
	OrderBy      string
	DepartmentId string
	IsActive     bool
}
