package util

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

const ErrUserAlreadyExists = "user already exists"
const ErrUserNotFound = "user not found"
const ErrInternalServerError = "internal server error"

const DateFormatYYYYMMDD = "2006-01-02"
const DateFormatYYYYMMDDTHHmmss = "2006-01-02T15:04:05"

var TimeNow = func() time.Time {
	return time.Now()
}

func ToDateTimeYYYYMMDD(dateString string) (dt time.Time, err error) {
	return time.Parse(DateFormatYYYYMMDD, dateString)
}

func ToDateTimeYYYYMMDDTHHmmss(dateString string) (dt time.Time, err error) {
	return time.Parse(DateFormatYYYYMMDDTHHmmss, dateString)
}

func HashPassword(input string) (string, error) {
	password := []byte(input)

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Comparing the password with the hash
	err = bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ValidatePassword(givenPlainTextPassword string, storedHashedPassword string) error {
	password := []byte(givenPlainTextPassword)
	hashedPassword := []byte(storedHashedPassword)
	// Comparing the password with the hash
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
