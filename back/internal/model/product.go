package model

import (
	"time"

	"github.com/lib/pq" // Импорт для поддержки массивов в PostgreSQL
)

// Product представляет основную сущность товара в каталоге.
type Product struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null;index" json:"name"` // Индекс по имени для поиска
	Price       float64   `gorm:"not null" json:"price"`               // Примечание: для фин. данных лучше использовать decimal
	OldPrice    float64   `json:"old_price,omitempty"`
	ImageURL    string    `gorm:"size:2048" json:"image_url"`
	ShortDesc   string    `gorm:"size:512" json:"short_desc"`
	FullDesc    string    `gorm:"type:text" json:"full_desc"`
	BrandID     int       `gorm:"index" json:"brand_id"`
	CategoryID  int       `gorm:"index" json:"category_id"`
	Rating      float64   `gorm:"default:0" json:"rating"`
	RatingCount int       `gorm:"default:0" json:"rating_count"`
	InStock     int       `gorm:"not null;default:0" json:"in_stock"` // Общее кол-во на складе для всех вариантов
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// --- Связи и сложные типы ---
	
	// Tags хранится как массив строк в PostgreSQL
	Tags pq.StringArray `gorm:"type:text[]" json:"tags"`

	// Один-ко-многим связи
	Variants     []ProductVariant `gorm:"foreignKey:ProductID" json:"variants"`
	PriceHistory []PricePoint     `gorm:"foreignKey:ProductID" json:"price_history"`
	Sales        []ProductSales   `gorm:"foreignKey:ProductID" json:"sales"`
}

// ProductVariant представляет вариант товара (например, по цвету или размеру).
type ProductVariant struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	ProductID int    `gorm:"not null;index" json:"product_id"`
	Color     string `gorm:"size:100;index" json:"color"`
	Size      string `gorm:"size:50;index" json:"size"`
	SKU       string `gorm:"size:255;not null;uniqueIndex" json:"sku"` // Артикул, должен быть уникальным
	InStock   int    `gorm:"not null;default:0" json:"in_stock"`
}

// PricePoint представляет точку в истории изменения цены товара.
type PricePoint struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	ProductID int       `gorm:"not null;index" json:"product_id"`
	Date      time.Time `gorm:"not null" json:"date"`
	Price     float64   `gorm:"not null" json:"price"`
}

// ProductSales представляет данные о продажах товара за определенную дату.
type ProductSales struct {
	ID         int       `gorm:"primaryKey" json:"id"`
	ProductID  int       `gorm:"not null;index:idx_product_date,unique" json:"product_id"` // Составной индекс
	Date       time.Time `gorm:"type:date;not null;index:idx_product_date,unique" json:"date"`  // Храним только дату, без времени
	SalesCount int       `gorm:"not null" json:"sales_count"`
}