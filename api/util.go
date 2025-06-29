package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pdrm26/hotel-reservation/core"
	"github.com/pdrm26/hotel-reservation/types"
)

func getAuthUser(c *fiber.Ctx) (*types.User, error) {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return nil, core.UnAuthorizedError()
	}

	return user, nil
}
