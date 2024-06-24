package models

type ListReq struct {
	IsActive     bool   `json:"is_active"`
	Page         string `json:"page"`
	Limit        string `json:"limit"`
	OrderBy      string `json:"order_by"`
	Value        string `json:"id"`
	DeleteStatus bool   `json:"-"`
}

type StatusRes struct {
	Status bool `json:"status"`
}

type AccessToken struct {
	Token string `json:"token"`
}
