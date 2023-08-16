package initializers

import "github.com/adhupraba/go-jwt-auth/models"

func MigrateDatabase() {
	DB.AutoMigrate(&models.User{})
}
