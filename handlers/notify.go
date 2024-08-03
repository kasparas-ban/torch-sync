package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/gomail.v2"
)

type NotifyReq struct {
	EmailAddress string
}

func EmailNotifyHandler(c *fiber.Ctx) error {
	var emailData NotifyReq
	if err := c.BodyParser(&emailData); err != nil {
		fmt.Printf("Error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email"})
	}

	email := os.Getenv("EMAIL")
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	hostAddress := os.Getenv("EMAIL_HOST_ADDRESS")
	hostPortStr := os.Getenv("EMAIL_PORT")
	hostPort, err := strconv.Atoi(hostPortStr)
	if err != nil {
		log.Fatal(err)
	}

	emailBody := fmt.Sprintf("Notify me when new features get added.<br><br>Email address: %v", emailData.EmailAddress)

	// Form the email
	m := gomail.NewMessage()
	m.SetHeader("From", "Torch App <"+email+">")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Torch: Notify for updates")
	m.SetBody("text/html", emailBody)

	d := gomail.NewDialer(
		hostAddress,
		hostPort,
		email,
		emailPassword,
	)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.SendStatus(http.StatusOK)
}
