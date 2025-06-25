package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pdrm26/hotel-reservation/db"
	"github.com/pdrm26/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
	NumGuests int       `json:"numGuests"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

func isRoomAvailable(ctx context.Context, bookingStore db.BookingStore, roomID primitive.ObjectID, startDate, endDate time.Time) (bool, error) {
	filter := bson.M{
		"roomID":    roomID,
		"startDate": bson.M{"$lt": endDate},
		"endDate":   bson.M{"$gt": startDate},
	}

	bookings, err := bookingStore.GetBookings(ctx, filter)
	if err != nil {
		return false, err
	}

	return len(bookings) == 0, nil
}

func (p BookRoomParams) validate() error {
	now := time.Now()

	// Validate that the booking is not in the past
	if now.After(p.StartDate) || now.After(p.EndDate) {
		return fmt.Errorf("cannot book a room in the past")
	}

	// Validate that start date is before end date
	if p.StartDate.After(p.EndDate) {
		return fmt.Errorf("cannot set enddate before startdate")
	}

	return nil
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	if err := params.validate(); err != nil {
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResp{
			Message: "internal server error",
		})
	}

	available, err := isRoomAvailable(c.Context(), h.store.Booking, roomID, params.StartDate, params.EndDate)
	if err != nil {
		return fmt.Errorf("could not check room availability: %w", err)
	}
	if !available {
		return c.Status(http.StatusBadRequest).JSON(genericResp{
			Message: "This room is already booked during the selected time.",
		})
	}

	book := types.Booking{
		UserID:    user.ID,
		RoomID:    roomID,
		NumGuests: params.NumGuests,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
	}

	inserted, err := h.store.Booking.InsertBooking(c.Context(), &book)
	if err != nil {
		return err
	}

	return c.JSON(inserted)
}
