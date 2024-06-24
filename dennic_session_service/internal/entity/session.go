package entity

import "time"

type Session struct {
	Id           string
	Order        int32
	IpAddress    string
	UserId       string
	FcmToken     string
	PlatformName string
	PlatformType string
	LoginAt      time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}

type SessionRequests struct {
	Id           string
	IpAddress    string
	UserId       string
	FcmToken     string
	PlatformName string
	PlatformType string
}

type StrReq struct {
	Id string
}

type Empty struct {
}

type StrUserReq struct {
	UserId   string
	IsActive bool
}

type UserSessionsList struct {
	UserSessions []*Session
	Count        int32
}

type SessionExistsResponse struct {
	IsExists bool
}
