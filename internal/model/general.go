package model

type Contact struct {
	Phone string `json:"phone"`
	Email string `json:"email"`
}

type Address struct {
	Street     string `json:"street"`
	PostalCode string `json:"postal_code"`
}
