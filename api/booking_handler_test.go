package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pdrm26/hotel-reservation/api/middleware"
	"github.com/pdrm26/hotel-reservation/db/fixtures"
	"github.com/pdrm26/hotel-reservation/types"
)

func TestNonOwnerCannotGetBooking(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	anotherUser := fixtures.AddUser(db.Store, "jack", "joe", false)
	user := fixtures.AddUser(db.Store, "jack", "joe", false)
	hotel := fixtures.AddHotel(db.Store, "Hilton", "UK", 5, nil)
	room := fixtures.AddRoom(db.Store, "small", true, 99.9, hotel.ID)
	booking := fixtures.AddBooking(db.Store, room.ID, user.ID, 3, time.Now(), time.Now().AddDate(0, 0, 3))

	app := fiber.New()
	route := app.Group("/", middleware.JWTAuthentication(db.User))
	bookingHandler := NewBookingHandler(db.Store)
	route.Get("/:id", bookingHandler.HandleGetBooking)

	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.Id.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(anotherUser))

	res, err := app.Test(req)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode == http.StatusOK {
		log.Fatalf("expected %d but got %d", http.StatusUnauthorized, res.StatusCode)
	}
}
func TestGetBooking(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	user := fixtures.AddUser(db.Store, "jack", "joe", false)
	hotel := fixtures.AddHotel(db.Store, "Hilton", "UK", 5, nil)
	room := fixtures.AddRoom(db.Store, "small", true, 99.9, hotel.ID)
	booking := fixtures.AddBooking(db.Store, room.ID, user.ID, 3, time.Now(), time.Now().AddDate(0, 0, 3))

	app := fiber.New()
	route := app.Group("/", middleware.JWTAuthentication(db.User))
	bookingHandler := NewBookingHandler(db.Store)
	route.Get("/:id", bookingHandler.HandleGetBooking)

	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.Id.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))

	res, err := app.Test(req)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalf("expected %d but got %d", http.StatusOK, res.StatusCode)
	}

	var book *types.Booking
	if err := json.NewDecoder(res.Body).Decode(&book); err != nil {
		log.Fatal(err)
	}

	if book.Id != booking.Id {
		t.Fatalf("expected %s but got %s", booking.Id, book.Id)
	}

}

func TestAdminCanRetrieveAllBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	adminUser := fixtures.AddUser(db.Store, "admin", "admin", true)
	hotel := fixtures.AddHotel(db.Store, "Hilton", "USA", 5, nil)
	room := fixtures.AddRoom(db.Store, "small", true, 99.10, hotel.ID)
	booking := fixtures.AddBooking(db.Store, room.ID, adminUser.ID, 2, time.Now(), time.Now().AddDate(0, 0, 2))

	app := fiber.New()
	admin := app.Group("/", middleware.JWTAuthentication(db.User), middleware.AdminAuth)
	bookingHandler := NewBookingHandler(db.Store)
	admin.Get("/", bookingHandler.HandleGetBookings)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))

	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected %d but got %d", http.StatusOK, res.StatusCode)
	}

	var bookings []*types.Booking
	if err := json.NewDecoder(res.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}

	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking but got %d", len(bookings))
	}

	if bookings[0].Id != booking.Id {
		t.Fatalf("expected %s but got %s", booking.Id, bookings[0].Id)
	}

}

func TestNormalUserCannotRetrieveAllBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	user := fixtures.AddUser(db.Store, "jack", "joe", false)
	hotel := fixtures.AddHotel(db.Store, "Hilton", "UK", 5, nil)
	room := fixtures.AddRoom(db.Store, "small", true, 99.9, hotel.ID)
	fixtures.AddBooking(db.Store, room.ID, user.ID, 3, time.Now(), time.Now().AddDate(0, 0, 3))

	app := fiber.New()
	admin := app.Group("/", middleware.JWTAuthentication(db.User), middleware.AdminAuth)
	bookingHandler := NewBookingHandler(db.Store)
	admin.Get("/", bookingHandler.HandleGetBookings)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))

	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected %d but got %d", http.StatusUnauthorized, res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	if string(body) != "not authorized" {
		t.Fatalf("expected not authorized but got %s", string(body))
	}

}
