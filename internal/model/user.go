package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type User struct {
	ID         int        `json:"id" gorm:"column:id"`
	Name       string     `json:"name" gorm:"column:name"`
	Email      string     `json:"email" gorm:"column:email"`
	Phone      string     `json:"phone" gorm:"column:phone"`
	Passsword  string     `json:"password" gorm:"column:password"`
	UserDetail UserDetail `json:"detail" gorm:"type:jsonb;column:detail"`
	CreatedAt  time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "users"
}

type UserDetail struct {
	Roles []string `json:"roles"`
}

// Implement the Valuer interface for Detail
func (d UserDetail) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Implement the Scanner interface for Detail
func (d *UserDetail) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan Detail")
	}
	return json.Unmarshal(bytes, d)
}
