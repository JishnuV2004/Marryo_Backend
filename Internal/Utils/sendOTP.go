package utils

import (
	"log"
	"net/smtp"

	"os"

	"github.com/joho/godotenv"
)

func SendOTPEmail(to, otp string) error {

	if err := godotenv.Load(".././.env"); err != nil {
		log.Println("No .env file found",err)
	}

	from := os.Getenv("E_MAIL")
	password := os.Getenv("APP_PASS")

	msg := []byte(
		"Subject: Marryo Email Verification\n\n" +
			"Your OTP is: " + otp + "\n\n" +
			"This OTP is valid for 5 minutes.",
	)

	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

	return smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		from,
		[]string{to},
		msg,
	)
}