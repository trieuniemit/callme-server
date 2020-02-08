package migrate

import (
	"github.com/jinzhu/gorm"
)

func Migrate(db *gorm.DB) {
	// db.Debug().AutoMigrate(&models.Error{}, &models.Todo{})
}
