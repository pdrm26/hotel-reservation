package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/pdrm26/hotel-reservation/db"
)

type AuthHandler struct {
	userStore db.UserStore
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var AuthParams AuthParams
	if err := c.BodyParser(&AuthParams); err != nil {
		return err
	}

	fmt.Println(AuthParams, "**************************")
	t := jwt.New(jwt.SigningMethodES256)
	s, err := t.SigningString()
	if err != nil {
		return err
	}
	return c.JSON(map[string]string{"token": s})
}
