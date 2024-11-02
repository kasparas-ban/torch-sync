package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"torch/torch-sync/middleware"
	"torch/torch-sync/storage"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/gofiber/fiber/v2"
)

type EmailAddresses struct {
	EmailAddress string `json:"email_address"`
}

type NewReqData struct {
	ID             string           `json:"id"`
	Username       string           `json:"username"`
	EmailAddresses []EmailAddresses `gorm:"embedded" json:"email_addresses"`
}

type NewUserReq struct {
	Data NewReqData `gorm:"embedded" json:"data"`
	Type string
}

type ConfirmSignInReq struct {
	ClerkID string `json:"clerkId"`
	Email   string `json:"email"`
}

func UserHandler(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).SendString("unauthorized")
	}

	user, err := storage.GetUser(userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}

func ConfirmSignInHandler(c *fiber.Ctx) error {
	var data ConfirmSignInReq
	if err := c.BodyParser(&data); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request object"})
	}

	// If user does not exist - create it
	if _, err := storage.GetUserByClerkID(data.ClerkID); err != nil {
		user, err := middleware.GetClerkClient().Users().Read(data.ClerkID)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
		}
		if user.Username == nil && *user.Username == "" {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "User data is corrupted"})
		}

		slog.Info("Signed in user not found - creating new one")
		newUser := storage.NewUser{ClerkID: data.ClerkID, Email: data.Email, Username: *user.Username}
		_, err = storage.AddUser(newUser)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "User data is invalid"})
		}
	}

	return c.SendStatus(http.StatusNoContent)
}

func RegisterUserHandler(c *fiber.Ctx) error {
	var userData storage.RegisterUserReq
	if err := c.BodyParser(&userData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user object"})
	}

	if err := userData.Validate(); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user object"})
	}

	var params clerk.CreateUserParams
	params.EmailAddresses = []string{userData.Email}
	params.Password = &userData.Password
	params.Username = &userData.Username

	clerkUser, err := middleware.GetClerkClient().Users().Create(params)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user"})
	}

	newUser, err := storage.RegisterUser(userData, clerkUser.ID)
	if err != nil {
		middleware.GetClerkClient().Users().Delete(clerkUser.ID)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user"})
	}

	return c.JSON(newUser)
}

func AddNewUserHandler(c *fiber.Ctx) error {
	var data NewUserReq
	if err := c.BodyParser(&data); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user object"})
	}

	if err := validatePayload(data); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user object"})
	}

	newUser := storage.NewUser{
		ClerkID:  data.Data.ID,
		Username: data.Data.Username,
		Email:    data.Data.EmailAddresses[0].EmailAddress,
	}

	user, err := storage.AddUser(newUser)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}

func UpdateUserHandler(c *fiber.Ctx) error {
	var userData storage.UpdateUserReq
	if err := c.BodyParser(&userData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user object"})
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).SendString("unauthorized")
	}

	if err := userData.Validate(); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user object"})
	}

	updatedUser, err := storage.UpdateUser(userID, userData)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Failed to update the user"})
	}

	return c.JSON(updatedUser)
}

func DeleteUserHandler(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized access"})
	}

	err = storage.DeleteUser(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(http.StatusOK)
}

func validatePayload(data NewUserReq) error {
	if data.Data.ID == "" || len(data.Data.EmailAddresses) == 0 || data.Data.Username == "" || data.Type != "user.created" {
		return errors.New("payload invalid")
	}
	return nil
}
