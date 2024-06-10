package common

import (
	"log"
	"net/smtp"
	"os"
	"strings"
)

func SendEmail(email string, token string) {

	LoadEnvFile()

	//Mailtrap account config
	username := os.Getenv("MAILTRAP_USERNAME")
	password := os.Getenv("MAILTRAP_PASSWORD")
	smtpHost := os.Getenv("MAILTRAP_SMT_HOST")
	smtpPort := os.Getenv("MAILTRAP_SMT_PORT")

	// Message data
	from := "thisisthebot@gmail.com"
	to := []string{email}
	subject := "Reset Password"
	htmlContent := "<p>To reset your password, visit the following link:</p><p><a href=\"http://localhost:8080/api/v1/reset-password?token=" + token + "\">Reset Password</a></p>"

	// MIME header
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := "Subject: " + subject + "\n" + mime + "\n" + "\n" + htmlContent
	msg = strings.Replace(msg, "\n", "\r\n", -1)

	fmt.Println("Hello from orther user")

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(smtpHost+":"+smtpPort,
		smtp.PlainAuth("", username, password, smtpHost),
		from, to, []byte(msg))

	if err != nil {
		log.Fatal(err)
	}
}
