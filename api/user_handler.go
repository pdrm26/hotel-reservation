package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pdrm26/hotel-reservation/db"
	"github.com/pdrm26/hotel-reservation/types"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var id = c.Params("id")

	user, err := h.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var userParams types.CreateUserParams
	if err := c.BodyParser(&userParams); err != nil {
		return err
	}

	user, err := types.NewUserFromParams(userParams)
	if err != nil {
		return c.JSON(err)
	}

	insertedUser, er := h.userStore.InsertUser(c.Context(), user)
	if er != nil {
		return er
	}

	return c.JSON(insertedUser)
}
