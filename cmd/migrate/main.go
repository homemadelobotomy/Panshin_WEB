package main

import (
	"lab/internal/app/ds"
	"lab/internal/app/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&ds.SolarPanel{},
		&ds.SolarPanelRequest{},
		&ds.RequestPanels{},
		&ds.User{},
	)
	if err != nil {
		panic("cant migrate db")
	}

}
