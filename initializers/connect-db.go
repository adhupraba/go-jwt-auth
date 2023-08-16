package initializers

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDb() {
	var err error
	dsn := os.Getenv("DB_DSN")
	DB, err = gorm.Open(postgres.Open(dsn))

	if err != nil {
		panic("Error connecting to database")
	}
}
