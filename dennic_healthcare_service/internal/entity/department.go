package entity

import "time"

type Department struct {
	Id               string
	Order            int32
	Name             string
	Description      string
	ImageUrl         string
	FloorNumber      int32
	ShortDescription string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        time.Time
}

type GetAll struct {
	Page     int64
	Limit    int64
	Field    string
	Value    string
	OrderBy  string
	IsActive bool
}

type ListDepartments struct {
	Count       int64
	Departments []Department
}
