package db

import (
	"testing"
	"time"

	"github.com/google/uuid"
)


func uuid7String() string {
	uuid, _ := uuid.NewV7()
	return uuid.String()
}

func TestWB_UpdateSupplier(t *testing.T) {
	tests := []struct {
		name   string
	}{
		{
			name: "Добавление бота",
		},
	}
	connDB()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Supplier{
				Products: &[]Product{
						{
							UUID: uuid7String(), 
							SupplierArticle: "P_STAR", 
							NmId: 181614365, 
							Barcode: "2038769615287", 
							Category: "Спортивные аксессуары",
							Subject: "Полупальцы",
							Brand: "MariStar",
							TechSize: "S",
							Stock: []Stock{
								{
									UUID: uuid7String(),
									LastChangeDate: time.Now(), 
									WarehouseName: "Электросталь",
									Quantity: 1, 
									InWayToClient: 2, 
									InWayFromClient: 3,
									QuantityFull: 3,
								},
								{
									UUID: uuid7String(),
									LastChangeDate: time.Now(), 
									WarehouseName: "Электросталь",
									Quantity: 1, 
									InWayToClient: 2, 
									InWayFromClient: 3,
									QuantityFull: 3,
								},
							},
						},
					},
			}
			connDB()
			s.UpdateSupplier()
		})
	}
}