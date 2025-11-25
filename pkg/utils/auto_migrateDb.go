package utils

import (
	"go/hioto/pkg/model"

	"gorm.io/gorm"
)

func AutoMigrateDb(db *gorm.DB) {
	db.AutoMigrate(&model.Registration{})
	db.AutoMigrate(&model.RuleDevice{})
}
