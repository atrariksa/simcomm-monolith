package model

import (
	"errors"
	"fmt"
	"strings"
)

type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

func (lr *LoginRequest) Validate() error {
	var errMessage string
	errTemplate := "%s is not valid;"
	if lr.Identifier == "" {
		errMessage += fmt.Sprintf(errTemplate, "email or phone")
	}

	//TODO: validate with more proper validataion

	if lr.Password == "" {
		errMessage += fmt.Sprintf(errTemplate, "password")
	}
	if errMessage != "" {
		return errors.New(errMessage)
	}
	return nil
}

type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (sur *SignUpRequest) Validate() error {
	var errMessage string
	errTemplate := "%s is not valid;"
	if sur.Name == "" {
		errMessage += fmt.Sprintf(errTemplate, "name")
	}
	if sur.Email == "" {
		errMessage += fmt.Sprintf(errTemplate, "email")
	}
	if sur.Phone == "" {
		errMessage += fmt.Sprintf(errTemplate, "phone")
	}
	if sur.Password == "" {
		errMessage += fmt.Sprintf(errTemplate, "password")
	}
	if strings.ToLower(sur.Role) != "customer" {
		errMessage += fmt.Sprintf(errTemplate, "role")
	}

	//TODO: validate with more proper validataion

	if errMessage != "" {
		return errors.New(errMessage)
	}
	return nil
}
