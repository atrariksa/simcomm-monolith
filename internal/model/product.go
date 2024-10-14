package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Product struct {
	ID        int           `json:"id" gorm:"column:id"`
	Code      string        `json:"code" gorm:"column:code"`
	Name      string        `json:"name" gorm:"column:name"`
	Status    string        `json:"status" gorm:"column:status"`
	Detail    ProductDetail `json:"detail" gorm:"type:jsonb;column:detail"`
	CreatedAt time.Time     `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time     `json:"updated_at" gorm:"column:updated_at"`
}

func (Product) TableName() string {
	return "products"
}

type ProductDetail struct {
	Weight   int    `json:"weight"`
	ImageURL string `json:"image_url"`
}

// Implement the Valuer interface for Detail
func (d ProductDetail) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Implement the Scanner interface for Detail
func (d *ProductDetail) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan Detail")
	}
	return json.Unmarshal(bytes, d)
}
