package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// --- Вспомогательные JSON-типы для GORM ---

// JSONBValue - общий интерфейс для типов, хранимых как JSONB
type JSONBValue interface {
	Scan(value interface{}) error
	Value() (driver.Value, error)
}

// scanJSON a helper function to scan jsonb value
func scanJSON(target interface{}, value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, target)
}

// valueJSON a helper function to return value for jsonb
func valueJSON(source interface{}) (driver.Value, error) {
	return json.Marshal(source)
}

// CustomerInfo - снимок данных о покупателе на момент заказа
type CustomerInfo struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	IsRegular   bool      `json:"is_regular"`
	OrdersCount int       `json:"orders_count"`
}

func (ci *CustomerInfo) Scan(value interface{}) error { return scanJSON(ci, value) }
func (ci CustomerInfo) Value() (driver.Value, error)  { return valueJSON(ci) }

// DeliveryInfo - информация о доставке
type DeliveryInfo struct {
	Type           string    `json:"type"`
	Address        Address   `json:"address"`
	EstimatedDate  time.Time `json:"estimated_date"`
	Cost           float64   `json:"cost"`
	TrackingNumber string    `json:"tracking_number"`
}

func (di *DeliveryInfo) Scan(value interface{}) error { return scanJSON(di, value) }
func (di DeliveryInfo) Value() (driver.Value, error)  { return valueJSON(di) }

// Address - адрес доставки
type Address struct {
	City        string      `json:"city"`
	Street      string      `json:"street"`
	House       string      `json:"house"`
	Apartment   string      `json:"apartment"`
	PostalCode  string      `json:"postal_code"`
	Coordinates Coordinates `json:"coordinates"`
}

// Coordinates - геолокационные координаты
type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// PaymentInfo - информация об оплате
type PaymentInfo struct {
	Method           string  `json:"method"`
	Status           string  `json:"status"`
	Amount           float64 `json:"amount"`
	CommissionLamoda float64 `json:"commission_lamoda"`
	SellerAmount     float64 `json:"seller_amount"`
	TransactionID    string  `json:"transaction_id"`
}

func (pi *PaymentInfo) Scan(value interface{}) error { return scanJSON(pi, value) }
func (pi PaymentInfo) Value() (driver.Value, error)  { return valueJSON(pi) }

// TotalsInfo - итоговые суммы по заказу
type TotalsInfo struct {
	Subtotal float64 `json:"subtotal"`
	Discount float64 `json:"discount"`
	Delivery float64 `json:"delivery"`
	Total    float64 `json:"total"`
}

func (ti *TotalsInfo) Scan(value interface{}) error { return scanJSON(ti, value) }
func (ti TotalsInfo) Value() (driver.Value, error)  { return valueJSON(ti) }

// --- Основные модели ---

// Order представляет заказ
type Order struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;index"` // ID продавца, которому принадлежит заказ
	CustomerID  uuid.UUID `gorm:"type:uuid;index"` // ID покупателя для фильтрации
	OrderNumber string    `gorm:"type:varchar(100);uniqueIndex" json:"order_number"`
	Date        time.Time `gorm:"index" json:"date"`
	Status      string    `gorm:"type:varchar(50);index" json:"status"`
	Notes       string    `gorm:"type:text" json:"notes"`
	CreatedAt   time.Time `json:"created_date"`
	UpdatedAt   time.Time `json:"updated_date"`

	// Данные, хранимые в JSONB для сохранения истории
	Customer CustomerInfo `gorm:"type:jsonb" json:"customer"`
	Delivery DeliveryInfo `gorm:"type:jsonb" json:"delivery"`
	Payment  PaymentInfo  `gorm:"type:jsonb" json:"payment"`
	Totals   TotalsInfo   `gorm:"type:jsonb" json:"totals"`

	// Связи "один ко многим"
	StatusHistory []StatusHistory `gorm:"foreignKey:OrderID" json:"status_history"`
	Items         []OrderItem     `gorm:"foreignKey:OrderID" json:"items"`
}

// OrderItem представляет товарную позицию в заказе
type OrderItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null;index" json:"-"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	VariantID uuid.UUID `gorm:"type:uuid;not null" json:"variant_id"`
	Name      string    `gorm:"type:varchar(255)" json:"name"`
	Brand     string    `gorm:"type:varchar(100)" json:"brand"`
	SKU       string    `gorm:"type:varchar(100)" json:"sku"`
	Size      string    `gorm:"type:varchar(50)" json:"size"`
	Color     string    `gorm:"type:varchar(100)" json:"color"`
	Image     string    `gorm:"type:varchar(2048)" json:"image"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CostPrice float64   `json:"cost_price"`
	Discount  float64   `json:"discount"`
	Total     float64   `json:"total"`
}

// StatusHistory представляет запись в истории изменения статуса заказа
type StatusHistory struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"-"`
	OrderID uuid.UUID `gorm:"type:uuid;not null;index" json:"-"`
	Status  string    `gorm:"type:varchar(50)" json:"status"`
	Date    time.Time `json:"date"`
	Comment string    `gorm:"type:text" json:"comment"`
}

// --- Структуры для запросов/ответов, не являющиеся моделями БД ---

// UpdateOrderStatusRequest представляет тело запроса на смену статуса
type UpdateOrderStatusRequest struct {
	Status                string     `json:"status" binding:"required"`
	Comment               string     `json:"comment"`
	EstimatedDeliveryDate *time.Time `json:"estimated_delivery_date"`
}
