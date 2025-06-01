package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/pdrm26/hotel-reservation/db"
	"github.com/pdrm26/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"message": "Not Found"})
		}
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

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	var userID = c.Params("id")
	if err := h.userStore.DeleteUser(c.Context(), userID); err != nil {
		return err
	}
	return c.JSON(map[string]string{"deleted": userID})
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		params types.UpdateUserParams
		userID = c.Params("id")
	)

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	if err := h.userStore.UpdateUser(c.Context(), bson.M{"_id": oid}, params); err != nil {
		return err
	}

	return c.JSON(map[string]string{"updated": userID})
}
