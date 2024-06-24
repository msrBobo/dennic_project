package models

type FieldValueReq struct {
	Field        string `json:"field"`
	Value        string `json:"value"`
	DeleteStatus bool   `json:"-"`
}

type ListReq struct {
	Page         string `json:"page"`
	Limit        string `json:"limit"`
	OrderBy      string `json:"order_by"`
	Field        string `json:"-"`
	Value        string `json:"value"`
	DeleteStatus bool   `json:"-"`
}

type StatusRes struct {
	Status bool `json:"status"`
}

type AccessToken struct {
	Token string `json:"token"`
}

type Date struct {
	Date string `json:"date"`
	Week string `json:"week"`
}

type SupportReq struct {
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Message     string `json:"message"`
}

type SupportRes struct {
	Message string `json:"message"`
}
