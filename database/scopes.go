package database

import "github.com/jinzhu/gorm"

// UserRelated is a scope that fetches records with a specified userID
func UserRelated(userID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	}
}
