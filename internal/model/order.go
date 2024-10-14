package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Order struct {
	ID        int         `json:"id" gorm:"column:id"`
	UserID    int         `json:"user_id" gorm:"column:user_id"`
	ShopID    int         `json:"shop_id" gorm:"column:shop_id"`
	Status    string      `json:"status" gorm:"column:status"`
	Detail    OrderDetail `json:"detail" gorm:"type:jsonb;column:detail"`
	CreatedAt time.Time   `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time   `json:"updated_at" gorm:"column:updated_at"`
}

func (Order) TableName() string {
	return "orders"
}

type OrderDetail struct {
	UserDetail      UserDetail      `json:"user_detail"`
	WarehouseDetail WarehouseDetail `json:"warehouse_detail"`
}

// Implement the Valuer interface for Detail
func (d OrderDetail) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Implement the Scanner interface for Detail
func (d *OrderDetail) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan Detail")
	}
	return json.Unmarshal(bytes, d)
}
