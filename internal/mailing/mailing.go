package mailing

import (
	"library/internal/models"
	"library/logger"

	"gorm.io/gorm"
)

func SendNewBookEmail(book models.Book, db *gorm.DB) {
	emails, err := GetSubscribers(db)
	if err != nil {
		logger.ErrorLog.Println("Failed to get subscribers: ", err)
		return
	}

	for _, email := range emails {
		go sendEmail(email)
	}
}

func GetSubscribers(db *gorm.DB) ([]string, error) {
	var emails []string
	err := db.Model(&models.User{}).Where("mailing - ?", true).Pluck("email", &emails).Error
	return emails, err
}
