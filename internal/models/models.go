package models

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBModel struct {
	Order OrderModel
}

func InitDB(dataSourceName string) (*DBModel, error) {
	db, err := gorm.Open(sqlite.Open(dataSourceName), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Include ALL tables in AutoMigrate
	err = db.AutoMigrate(
		&Order{},
		&OrderItem{},
		&OrderItemTopping{},
		&OrderItemDietaryRequirement{},
		&OrderItemAllergy{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	dbModel := &DBModel{
		Order: OrderModel{DB: db},
	}

	return dbModel, nil
}
