package db

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	UUID string  `gorm:"primaryKey"` // уникальный номер
	SupplierArticle string // Артикул продавца
	NmId int64 // Артикул WB
	Barcode string `gorm:"uniqueIndex"`// Баркод
	Category string // Категория
	Subject string // Предмет
	Brand string // Бренд
	TechSize string // Размер
	Price decimal.Decimal `gorm:"-"` // Цена
	Discount bool // Скидка
	Stock []Stock `gorm:"foreignKey:ProductID"`
	CreatedAt time.Time 
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Stock struct {
	UUID string `gorm:"primaryKey"`
	ProductID string
	LastChangeDate time.Time // Дата и время обновления информации в сервисе.
	WarehouseName string // Название склада
	Quantity int64 // Количество, доступное для продажи (сколько можно добавить в корзину)
	InWayToClient int64 // В пути к клиенту
	InWayFromClient int64 // В пути от клиента
	QuantityFull int64 // Полное (непроданное) количество, которое числится за складом (= quantity + в пути)
	CreatedAt time.Time 
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Supplier struct {
	Products *[]Product
}


type WB interface {
	UpdateSupplier()
	UpdateStocks()
}

func (s *Supplier) UpdateSupplier() {
	// var products = s.Products
	// conDB.Migrator().CreateTable(&Stock{})
	// if !conDB.Migrator().HasTable(products) {
	// 	conDB.Migrator().CreateTable(products)
	// 	conDB.CreateInBatches(products, 100)
	// }
	// Processing records in batches of 100
	// conDB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Supplier{})
	// result := conDB.Create(&s) // pass pointer of data to Create
	// fmt.Print(result)
 
}