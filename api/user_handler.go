package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pdrm26/hotel-reservation/types"
)

func HandleGetUsers(c *fiber.Ctx) error {
	return c.JSON(types.User{FirstName: "Pedram Baradarian", LastName: "Jamshid rostami"})
}

func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"user": "Pedram Baradarian"})
}
