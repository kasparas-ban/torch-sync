package handlers

import (
	"net/http"
	"torch/torch-sync/middleware"
	"torch/torch-sync/storage"

	"github.com/gofiber/fiber/v2"
)

func ItemsHandler(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).SendString("unauthorized")
	}

	items, err := storage.GetAllItemsByUser(userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(items)
}
