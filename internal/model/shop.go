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
	Location  string     `json:"location" gorm:"column:location"`
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

type ShopProduct struct {
	ID        int                `json:"id" gorm:"column:id"`
	ProductID int                `json:"product_id" gorm:"column:product_id"`
	ShopID    int                `json:"shop_id" gorm:"column:shop_id"`
	Status    string             `json:"status" gorm:"column:status"`
	Stock     int                `json:"stock" gorm:"column:stock"`
	Detail    ShopProductDetails `json:"detail" gorm:"type:jsonb;column:detail"`
	CreatedAt time.Time          `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time          `json:"updated_at" gorm:"column:updated_at"`
}

func (ShopProduct) TableName() string {
	return "shop_products"
}

type ShopProductDetails struct {
	ShopProductDetails []ShopProductDetail `json:"shop_product_details"`
}

type ShopProductDetail struct {
	WarehouseID     int    `json:"warehouse_id"`
	WarehouseStatus string `json:"warehouse_status"`
	Stock           int    `json:"stock"`
}

// Implement the Valuer interface for Detail
func (d *ShopProductDetails) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Implement the Scanner interface for Detail
func (d *ShopProductDetails) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan Detail")
	}
	return json.Unmarshal(bytes, d)
}

type TransferProduct struct {
	ID                     int                   `json:"id" gorm:"column:id"`
	ShopProductID          int                   `json:"shop_product_id" gorm:"column:shop_product_id"`
	StockToTransfer        int                   `json:"stock_to_transfer" gorm:"column:stock_to_transfer"`
	WarehouseIDSource      int                   `json:"warehouse_id_source" gorm:"column:warehouse_id_source"`
	WarehouseIDDestination int                   `json:"warehouse_id_destination" gorm:"column:warehouse_id_destination"`
	Status                 string                `json:"status" gorm:"column:status"`
	Detail                 TransferProductDetail `json:"detail" gorm:"type:jsonb;column:detail"`
	CreatedAt              time.Time             `json:"created_at" gorm:"column:created_at"`
	UpdatedAt              time.Time             `json:"updated_at" gorm:"column:updated_at"`
}

func (TransferProduct) TableName() string {
	return "transferred_products"
}

type TransferProductDetail struct {
	Histories []TransferProductHostory `json:"histories"`
}

type TransferProductHostory struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// Implement the Valuer interface for Detail
func (d *TransferProductDetail) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Implement the Scanner interface for Detail
func (d *TransferProductDetail) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan Detail")
	}
	return json.Unmarshal(bytes, d)
}
