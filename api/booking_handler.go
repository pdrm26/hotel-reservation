package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/pdrm26/hotel-reservation/db"
	"github.com/pdrm26/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}

	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	bookingID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}

	bookings, err := h.store.Booking.GetBookingByID(c.Context(), bookingID)
	if err != nil {
		return err
	}

	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(genericResp{Message: "Unauthorized"})
	}

	if bookings.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResp{Message: "Unauthorized"})
	}

	return c.JSON(bookings)
}
