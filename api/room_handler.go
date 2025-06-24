package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pdrm26/hotel-reservation/db"
	"github.com/pdrm26/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
	NumGuests int       `json:"numGuests"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

func (p BookRoomParams) validate() error {
	now := time.Now()
	if now.After(p.StartDate) || now.After(p.EndDate) {
		return fmt.Errorf("cannot book a room in the past")
	}
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

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.validate(); err != nil {
		return err
	}

	roomId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResp{
			Message: "internal server error",
		})
	}

	book := types.Booking{
		UserID:    user.ID,
		RoomID:    roomId,
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
