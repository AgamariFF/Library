package database

import "library/internal/models"

// Migrate создает таблицы на основе моделей
func Migrate() error {
	err := DB.AutoMigrate(&models.Book{}, &models.Genre{}, &models.User{})
	if err != nil {
		return err
	}
	return nil
}