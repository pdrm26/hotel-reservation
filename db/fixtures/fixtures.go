package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pdrm26/hotel-reservation/db"
	"github.com/pdrm26/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(store *db.Store, fname, lname string, isAdmin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     fmt.Sprintf("%s@%s.com", fname, lname),
		Password:  fmt.Sprintf("%s_%s", fname, lname),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	insertedUser, errr := store.User.InsertUser(context.TODO(), user)
	if errr != nil {
		log.Fatal(errr)
	}
	return insertedUser
}

func AddHotel(store *db.Store, name, location string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	var roomIDS = rooms
	if rooms == nil {
		roomIDS = []primitive.ObjectID{}
	}
	var hotel = types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    roomIDS,
		Rating:   rating,
	}

	insertedHotel, err := store.Hotel.InsertHotel(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func AddRoom(store *db.Store, size string, seaSide bool, price float64, hotelID primitive.ObjectID) *types.Room {
	var room = types.Room{
		Size:    size,
		Seaside: seaSide,
		Price:   price,
		HotelID: hotelID,
	}
	insertedRoom, err := store.Room.InsertRoom(context.TODO(), &room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func AddBooking(store *db.Store, roomId, userId primitive.ObjectID, numGuests int, startDate, endDate time.Time) *types.Booking {

	book := &types.Booking{
		RoomID:    roomId,
		UserID:    userId,
		NumGuests: numGuests,
		StartDate: startDate,
		EndDate:   endDate,
	}

	insertedBook, err := store.Booking.InsertBooking(context.TODO(), book)
	if err != nil {
		log.Fatal(err)
	}

	return insertedBook

}
