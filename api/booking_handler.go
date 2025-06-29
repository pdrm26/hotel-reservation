package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pdrm26/hotel-reservation/core"
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
		return core.BadRequestError()
	}

	booking, err := h.store.Booking.GetBookingByID(c.Context(), bookingID)
	if err != nil {
		return core.NotFoundError("booking")
	}

	user, err := getAuthUser(c)
	if err != nil {
		return err
	}

	if user.ID != booking.UserID {
		return core.ForbiddenError()
	}

	if err := h.store.Booking.UpdateBookingByID(c.Context(), bookingID, bson.M{"canceled": true}); err != nil {
		return err
	}
	return c.JSON(genericResp{Message: "Updated"})

}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return core.NotFoundError("bookings")
	}

	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	bookingID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return core.NotFoundError("booking")
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
		return core.ForbiddenError()
	}

	return c.JSON(bookings)
}
