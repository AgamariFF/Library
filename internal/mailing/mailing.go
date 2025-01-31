package mailing

import (
	"bytes"
	"library/internal/models"
	"library/logger"
	"os"
	"strconv"
	"text/template"

	"gopkg.in/gomail.v2"

	"gorm.io/gorm"
)

type EmailData struct {
	Title           string
	Author          string
	Genres          string
	Description     string
	BookLink        string
	UnsubscribeLink string
}

func SendNewBookEmail(book models.Book, db *gorm.DB) {
	genres := book.Genres[0].Name
	if len(book.Genres) > 1 {
		for i, genre := range book.Genres {
			if i == 0 {
				continue
			}
			genres = genres + ", " + genre.Name
		}
	}

	emailBook := EmailData{
		Title:           book.Title,
		Author:          book.Author,
		Genres:          genres,
		Description:     book.Description,
		BookLink:        "http://localhost:8080/getBook?bookId=" + strconv.Itoa(int(book.ID)),
		UnsubscribeLink: "",
	}

	emails, err := GetSubscribers(db)
	if err != nil {
		logger.ErrorLog.Println("Failed to get subscribers: ", err)
		return
	}
	logger.InfoLog.Println("Geting subscribers for mailing succesfully")

	html, err := GenerateEmailNewBookBody(emailBook)
	if err != nil {
		logger.ErrorLog.Println("Failed to create html body to send email about new book: ", err)
	}
	logger.InfoLog.Println("Generate html body for mailing succesfully")

	SendEmail(emails, "Новая книга доступна!", html)
}

func GetSubscribers(db *gorm.DB) ([]string, error) {
	var emails []string
	err := db.Model(&models.User{}).Where("mailing = ?", true).Pluck("email", &emails).Error
	return emails, err
}

func SendEmail(to []string, subject, body string) {
	from := os.Getenv("SMTP_Name")
	password := os.Getenv("SMTP_Password")
	if (from == "") || (password == "") {
		logger.ErrorLog.Println(`"SMTP_Name" or "SMTP_Password" empty in .env`)
		return
	}

	logger.InfoLog.Println("SMTP_Name:"+from, "\tSMTP_Password:"+password)

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", from)
	mailer.SetHeader("To", to...)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)
	dialer := gomail.NewDialer("smtp.mail.ru", 465, from, password)
	dialer.SSL = true

	if err := dialer.DialAndSend(mailer); err != nil {
		logger.ErrorLog.Printf("Failed to send email: %+v", err)
		return
	}

	logger.InfoLog.Println("Email sent to: ", to)
}

func GenerateEmailNewBookBody(book EmailData) (string, error) {
	html, err := os.ReadFile("HTML/NewBook.html")
	if err != nil {
		return "", err
	}

	emailTemplate := string(html)

	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		return "", err
	}
	logger.InfoLog.Println(book)
	var body bytes.Buffer
	if err := tmpl.Execute(&body, book); err != nil {
		return "", err
	}

	return body.String(), nil
}
