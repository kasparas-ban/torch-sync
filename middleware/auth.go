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
var ClerkUserIDContext = "clerk_id"

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
	err := VerifyToken(c, sessionToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).SendString("unauthorized")
	}

	return c.Next()
}

func VerifyToken(c *fiber.Ctx, token string) error {
	sessClaims, err := client.VerifyToken(token)
	if err != nil {
		return errors.New("unauthorized")
	}

	userID, err := saveClerkIDContext(c, sessClaims)
	if err != nil {
		return errors.New("unauthorized")
	}

	// Add userID to the current context
	c.Locals(UserID_metadata, userID)

	return nil
}

func saveClerkIDContext(c *fiber.Ctx, claims *clerk.SessionClaims) (string, error) {
	// Read clerkID
	user, err := client.Users().Read(claims.Claims.Subject)
	if err != nil {
		return "", err
	}

	// Add Clerk user ID to the context
	c.Locals(ClerkUserIDContext, user.ID)

	userID := user.PrivateMetadata.(map[string]interface{})[UserID_metadata]
	if userID == nil || userID == "" {
		// Add userID to the Clerk metadata
		userID, err = AddUserID(user)
		if err != nil {
			return "", errors.New("user not found")
		}
	}

	return userID.(string), nil
}

func AddUserID(user *clerk.User) (string, error) {
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

func GetClerkUserID(c *fiber.Ctx) (string, error) {
	clerkID, ok := c.Locals(ClerkUserIDContext).(string)
	if !ok || clerkID == "" {
		return "", errors.New("clerk ID not found")
	}

	return clerkID, nil
}

func GetClerkClient() clerk.Client {
	return client
}
