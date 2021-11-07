package database

import (
	"github.com/jinzhu/gorm"
	"github.com/mcaubrey/go_rest_api/internal/services/comment"
)

// MigrateDB - migrates our database and creates the comment table
func MigrateDB(db *gorm.DB) error {
	if result := db.AutoMigrate(comment.Comment{}); result.Error != nil {
		return result.Error
	}
	return nil
}
