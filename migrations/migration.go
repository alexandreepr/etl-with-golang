package migrations

import (
	"etl-with-golang/infra/database"
	"etl-with-golang/models"
)

func Migrate() {
	var migrationModels = []interface{}{&models.Register{}}
	err := database.DB.AutoMigrate(migrationModels...)
	if err != nil {
		return
	}
}
