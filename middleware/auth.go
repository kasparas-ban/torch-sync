package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"
	"torch/torch-sync/config"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/gofiber/fiber/v2"
)

var client clerk.Client

func InitClerk(env config.EnvVars) {
	var err error
	client, err = clerk.NewClient(env.CLERK_SECRET_KEY)
	clerk.WithSessionV2(
		client,
		clerk.WithLeeway(5*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func AuthMiddleware(c *fiber.Ctx) error {
	sessionToken := c.Get("Authorization")
	sessionToken = strings.TrimPrefix(sessionToken, "Bearer ")
	_, err := client.VerifyToken(sessionToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON("unauthorized")
	}

	return c.Next()
}
