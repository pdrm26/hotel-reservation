package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pdrm26/hotel-reservation/core"
	"github.com/pdrm26/hotel-reservation/db"
	"github.com/pdrm26/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelsResp struct {
	Data    []*types.Hotel `json:"data"`
	Page    int64          `json:"page"`
	Results int            `json:"results"`
}

type HotelQueryParams struct {
	db.PaginateFilter
	Rating int
}

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{store: store}
}

func (h *HotelHandler) HanldeGetHotels(c *fiber.Ctx) error {
	var params HotelQueryParams
	if err := c.QueryParser(&params); err != nil {
		return core.BadRequestError()
	}
	filter := bson.M{
		"rating": params.Rating,
	}
	hotels, err := h.store.Hotel.GetHotels(c.Context(), filter, &params.PaginateFilter)
	if err != nil {
		return core.NotFoundError("hotels")
	}

	return c.JSON(HotelsResp{Data: hotels, Page: params.Page, Results: len(hotels)})
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
