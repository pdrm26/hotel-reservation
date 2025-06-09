package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/pdrm26/hotel-reservation/db"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleAuthenticate(c *fiber.Ctx) error {
	t := jwt.New(jwt.SigningMethodES256)
	s, err := t.SigningString()
	if err != nil {
		return err
	}
	return c.JSON(map[string]string{"token": s})
}
