package initializers

import "github.com/rathink4/rest-crud-go/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
