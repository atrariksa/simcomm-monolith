package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Shop struct {
	ID        int        `json:"id" gorm:"column:id"`
	UserID    int        `json:"user_id" gorm:"column:user_id"`
	Name      string     `json:"name" gorm:"column:name"`
	Status    string     `json:"status" gorm:"column:status"`
	Detail    ShopDetail `json:"detail" gorm:"type:jsonb;column:detail"`
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"column:updated_at"`
}

func (Shop) TableName() string {
	return "shops"
}

type ShopDetail struct {
	Contact   Contact   `json:"contact"`
	Addresses []Address `json:"addresses"`
	ImageURL  string    `json:"image_url"`
}

// Implement the Valuer interface for Detail
func (d ShopDetail) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Implement the Scanner interface for Detail
func (d *ShopDetail) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan Detail")
	}
	return json.Unmarshal(bytes, d)
}
