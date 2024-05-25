package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"torch/torch-sync/config"

	"torch/torch-sync/storage"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/gofiber/fiber/v2"
)

var client clerk.Client

var UserID_metadata = "user_id"

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
	sessClaims, err := client.VerifyToken(sessionToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON("unauthorized")
	}

	err = saveClerkIDContext(c, sessClaims)
	if err != nil {
		return c.Status(http.StatusUnauthorized).SendString("unauthorized")
	}
	return c.Next()
}

func saveClerkIDContext(c *fiber.Ctx, claims *clerk.SessionClaims) error {
	// Read clerkID
	user, err := client.Users().Read(claims.Claims.Subject)
	if err != nil {
		return err
	}

	userID := user.PrivateMetadata.(map[string]interface{})[UserID_metadata]
	if userID == nil || userID == "" {
		// Add userID to the Clerk metadata
		userID, err = AddUserID(c, user)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("user not found")
		}
	}

	// Add userID to the current context
	c.Locals(UserID_metadata, userID)

	return nil
}

func AddUserID(c *fiber.Ctx, user *clerk.User) (string, error) {
	userInfo, err := storage.GetUserByClerkID(user.ID)
	if err != nil {
		return "", err
	}
	if userInfo.UserID == "" {
		return "", errors.New("failed to get User ID")
	}

	currPrivMetadata := make(map[string]interface{})
	currPrivMetadata[UserID_metadata] = userInfo.UserID

	var keyValStrings []string
	for key, value := range currPrivMetadata {
		keyValStrings = append(keyValStrings, fmt.Sprintf("\"%s\": \"%s\"", key, value))
	}
	finalStr := fmt.Sprintf("{%s}", strings.Join(keyValStrings, ","))
	_, err = client.Users().Update(user.ID, &clerk.UpdateUser{
		PrivateMetadata: finalStr,
	})
	if err != nil {
		return "", err
	}

	return userInfo.UserID, nil
}

func GetUserID(c *fiber.Ctx) (string, error) {
	userID, ok := c.Locals(UserID_metadata).(string)
	if !ok || userID == "" {
		return "", errors.New("user ID not found")
	}

	return userID, nil
}
