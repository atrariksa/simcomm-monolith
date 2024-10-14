package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Warehouse struct {
	ID        int             `json:"id" gorm:"column:id"`
	ShopID    int             `json:"shop_id" gorm:"column:shop_id"`
	Name      string          `json:"name" gorm:"column:name"`
	Location  string          `json:"location" gorm:"column:location"`
	Status    string          `json:"status" gorm:"column:status"`
	Detail    WarehouseDetail `json:"detail" gorm:"type:jsonb"`
	CreatedAt time.Time       `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"column:updated_at"`
}

func (Warehouse) TableName() string {
	return "warehouses"
}

type WarehouseDetail struct {
	Contact   Contact   `json:"contact"`
	Addresses []Address `json:"addresses"`
	ImageURL  string    `json:"image_url"`
}

type WarehouseStoredProduct struct {
	ID              int    `json:"id" gorm:"column:id"`
	ShopProductID   int    `json:"shop_product_id" gorm:"column:shop_product_id"`
	ShopProductName string `json:"shop_product_name" gorm:"column:shop_product_name"`
	Stock           int    `json:"stock" gorm:"column:stock"`
}

// Implement the Valuer interface for Detail
func (d WarehouseDetail) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Implement the Scanner interface for Detail
func (d *WarehouseDetail) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan Detail")
	}
	return json.Unmarshal(bytes, d)
}
