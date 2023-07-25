package initializers

import "github.com/f0rSaaaa/JWTAuthenticationGO/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
