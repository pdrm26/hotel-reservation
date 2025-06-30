package api

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pdrm26/hotel-reservation/core"
	"github.com/pdrm26/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{store: store}
}

func (h *HotelHandler) HanldeGetHotels(c *fiber.Ctx) error {
	pageParam := c.Query("page", "1")
	limitParam := c.Query("limit", "10")

	fmt.Println(pageParam, limitParam)

	page, err := strconv.Atoi(pageParam)
	fmt.Println(err)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitParam)
	fmt.Println(err)
	if err != nil || limit < 1 {
		limit = 10
	}

	fmt.Println(page, limit)

	opts := options.FindOptions{}
	opts.SetSkip(int64(page-1) * int64(limit))
	opts.SetLimit(int64(limit))
	hotels, err := h.store.Hotel.GetHotels(c.Context(), nil, &opts)
	if err != nil {
		return core.NotFoundError("hotels")
	}

	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	hotelId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return core.InvalidIDError()
	}

	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), hotelId)
	if err != nil {
		return core.NotFoundError("hotel")
	}

	return c.JSON(hotel)
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	hotelIDStr := c.Params("id")
	hotelID, err := primitive.ObjectIDFromHex(hotelIDStr)
	if err != nil {
		return core.InvalidIDError()
	}
	rooms, err := h.store.Room.GetRooms(c.Context(), bson.M{"hotelID": hotelID})
	if err != nil {
		return core.NotFoundError("rooms")
	}
	return c.JSON(rooms)
}
