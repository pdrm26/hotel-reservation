package middleware

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/pdrm26/hotel-reservation/types"
)

func AdminAuth(c *fiber.Ctx) error {

	user, ok := c.Context().UserValue("user").(*types.User)

	if !ok {
		return fmt.Errorf("not authorized")
	}

	if !user.IsAdmin {
		return c.Status(http.StatusUnauthorized).SendString("not authorized")

	}

	return c.Next()

}
