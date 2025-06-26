package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/pdrm26/hotel-reservation/db"
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

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	bookingID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}

	booking, err := h.store.Booking.GetBookingByID(c.Context(), bookingID)
	if err != nil {
		return err
	}

	user, err := getAuthUser(c)
	if err != nil {
		return err
	}

	if user.ID != booking.UserID {
		return c.Status(http.StatusForbidden).JSON(genericResp{Message: "Forbidden"})
	}

	if err := h.store.Booking.UpdateBookingByID(c.Context(), bookingID, bson.M{"canceled": true}); err != nil {
		return err
	}
	return c.JSON(genericResp{Message: "Updated"})

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

	user, err := getAuthUser(c)
	if err != nil {
		return err
	}

	if bookings.UserID != user.ID {
		return c.Status(http.StatusForbidden).JSON(genericResp{Message: "Forbidden"})
	}

	return c.JSON(bookings)
}
