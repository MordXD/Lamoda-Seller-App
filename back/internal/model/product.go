package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// --- Вспомогательные JSON-типы для GORM ---

// Dimensions представляет размеры (длина, ширина, высота)
type Dimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// Scan реализует интерфейс Scanner для типа Dimensions
func (d *Dimensions) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &d)
}

// Value реализует интерфейс Valuer для типа Dimensions
func (d Dimensions) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// --- Основные модели ---

// Product представляет основную сущность товара.
type Product struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name              string         `gorm:"type:varchar(255);not null;index" json:"name"`
	Description       string         `gorm:"type:text" json:"description"`
	Brand             string         `gorm:"type:varchar(100);index" json:"brand"`
	Category          string         `gorm:"type:varchar(100);index" json:"category"`
	Subcategory       string         `gorm:"type:varchar(100);index" json:"subcategory"`
	SKU               string         `gorm:"type:varchar(100);uniqueIndex" json:"sku"` // Основной артикул товара
	Barcode           string         `gorm:"type:varchar(100);index" json:"barcode"`
	Price             float64        `gorm:"not null" json:"price"`
	CostPrice         float64        `json:"cost_price"`
	Currency          string         `gorm:"type:varchar(10);default:'RUB'" json:"currency"`
	TotalStock        int            `gorm:"default:0" json:"total_stock"` // Суммарный остаток по всем вариантам
	Rating            float64        `gorm:"default:0" json:"rating"`
	ReviewsCount      int            `gorm:"default:0" json:"reviews_count"`
	ReturnRate        float64        `gorm:"default:0" json:"return_rate"`
	Status            string         `gorm:"type:varchar(50);default:'draft';index" json:"status"` // active, inactive, draft
	SeasonalDemand    string         `gorm:"type:varchar(100)" json:"seasonal_demand"`
	IsBestseller      bool           `gorm:"default:false" json:"is_bestseller"`
	IsNew             bool           `gorm:"default:true" json:"is_new"`
	DiscountPercent   float64        `gorm:"default:0" json:"discount_percent"`
	Tags              pq.StringArray `gorm:"type:text[]" json:"tags"`
	Material          string         `gorm:"type:varchar(255)" json:"material"`
	CareInstructions  string         `gorm:"type:text" json:"care_instructions"`
	CountryOrigin     string         `gorm:"type:varchar(100)" json:"country_origin"`
	SupplierID        *uuid.UUID     `gorm:"type:uuid" json:"-"` // Ссылка на поставщика
	CreatedAt         time.Time      `json:"created_date"`
	UpdatedAt         time.Time      `json:"updated_date"`

	// --- Связи ---
	Supplier *Supplier        `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Images   []ProductImage   `gorm:"foreignKey:ProductID" json:"images,omitempty"`
	Variants []ProductVariant `gorm:"foreignKey:ProductID" json:"variants,omitempty"`

	// --- Поля, которые не хранятся в БД, а вычисляются ---
	MarginPercent   float64 `gorm:"-" json:"margin_percent"`
	MainImage       string  `gorm:"-" json:"main_image,omitempty"`
	AvailableSizes  []string `gorm:"-" json:"available_sizes,omitempty"`
	AvailableColors []string `gorm:"-" json:"available_colors,omitempty"`
	SalesCount30d   int     `gorm:"-" json:"sales_count_30d,omitempty"` // Пример вычисляемого поля
	Revenue30d      float64 `gorm:"-" json:"revenue_30d,omitempty"`     // Пример вычисляемого поля
}

// ProductVariant представляет вариант товара (SKU).
type ProductVariant struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ProductID  uuid.UUID  `gorm:"type:uuid;not null;index" json:"-"`
	SKU        string     `gorm:"type:varchar(100);not null;uniqueIndex" json:"sku"`
	Size       string     `gorm:"type:varchar(50);index" json:"size"`
	Color      string     `gorm:"type:varchar(100);index" json:"color"`
	ColorHex   string     `gorm:"type:varchar(20)" json:"color_hex"`
	Stock      int        `gorm:"default:0" json:"stock"`
	Reserved   int        `gorm:"default:0" json:"reserved"`
	Price      float64    `json:"price"` // Может отличаться от основного
	Weight     float64    `json:"weight"` // в граммах
	Dimensions Dimensions `gorm:"type:jsonb" json:"dimensions"`
}

// ProductImage представляет изображение товара.
type ProductImage struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;index" json:"-"`
	URL       string    `gorm:"type:varchar(2048);not null" json:"url"`
	AltText   string    `gorm:"type:varchar(255)" json:"alt"`
	IsMain    bool      `gorm:"default:false;index" json:"is_main"`
	Order     int       `gorm:"default:0" json:"order"`
}

// Supplier представляет поставщика.
type Supplier struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name    string    `gorm:"type:varchar(255);not null" json:"name"`
	Contact string    `gorm:"type:varchar(255)" json:"contact"`
}

// Category представляет иерархическую структуру категорий
type Category struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Subcategories []Category `json:"subcategories,omitempty"`
}

// SizeChart представляет размерную сетку для категории.
type SizeChart struct {
	Category string `json:"category"`
	Type     string `json:"type"` // clothing, shoes, accessories
	Sizes    []Size `json:"sizes"`
}

// Size представляет конкретный размер в сетке.
type Size struct {
	Size          string            `json:"size"`
	Measurements  map[string]string `json:"measurements"`
	International string            `json:"international"`
	US            string            `json:"us"`
}

// AfterFind хук GORM для вычисления полей после выборки из БД.
func (p *Product) AfterFind(tx *gorm.DB) (err error) {
	// Вычисляем маржинальность
	if p.Price > 0 {
		p.MarginPercent = ((p.Price - p.CostPrice) / p.Price) * 100
	}

	// Находим главное изображение
	for _, img := range p.Images {
		if img.IsMain {
			p.MainImage = img.URL
			break
		}
	}
    // Если главного нет, берем первое по порядку или просто первое
	if p.MainImage == "" && len(p.Images) > 0 {
		p.MainImage = p.Images[0].URL
	}

	// Собираем доступные размеры и цвета
	sizeSet := make(map[string]struct{})
	colorSet := make(map[string]struct{})
	for _, variant := range p.Variants {
		if variant.Stock > variant.Reserved {
			if variant.Size != "" {
				sizeSet[variant.Size] = struct{}{}
			}
			if variant.Color != "" {
				colorSet[variant.Color] = struct{}{}
			}
		}
	}
	for s := range sizeSet {
		p.AvailableSizes = append(p.AvailableSizes, s)
	}
	for c := range colorSet {
		p.AvailableColors = append(p.AvailableColors, c)
	}
	
	// TODO: Данные по продажам (SalesCount30d, Revenue30d) должны загружаться
	// отдельным запросом в репозитории/сервисе и заполняться здесь.
	// Для примера оставим их нулевыми.
	p.SalesCount30d = 45 // Заглушка
	p.Revenue30d = 400500 // Заглушка

	return
}