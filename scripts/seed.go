package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pdrm26/hotel-reservation/api"
	"github.com/pdrm26/hotel-reservation/db"
	"github.com/pdrm26/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx          = context.Background()
	client       *mongo.Client
	hotelStore   db.HotelStore
	roomStore    db.RoomStore
	userStore    db.UserStore
	bookingStore db.BookingStore
)

func init() {
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
	bookingStore = db.NewMongoBookingStore(client)

}

func seedUser(fname, lname, email, password string, isAdmin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  password,
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	insertedUser, errr := userStore.InsertUser(ctx, user)
	if errr != nil {
		log.Fatal(errr)
	}
	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))
	return insertedUser
}

func seedHotel(name, location string, rating int) *types.Hotel {
	var hotel = types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func seedRoom(size string, price float64, hotelID primitive.ObjectID) *types.Room {

	var room = types.Room{Size: size, Price: price, HotelID: hotelID}
	insertedRoom, err := roomStore.InsertRoom(ctx, &room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func seedBooking(roomId, userId primitive.ObjectID, numGuests int, startDate, endDate time.Time) *types.Booking {

	book := &types.Booking{
		RoomID:    roomId,
		UserID:    userId,
		NumGuests: numGuests,
		StartDate: startDate,
		EndDate:   endDate,
	}

	insertedBook, err := bookingStore.InsertBooking(ctx, book)
	if err != nil {
		log.Fatal(err)
	}

	return insertedBook

}

func main() {
	admin := seedUser("admin", "admin", "admin@admin.com", "admin123", true)
	pedram := seedUser("pedram", "baradarian", "pedram@gmail.com", "123123123", false)
	jack := seedUser("jack", "joe", "jack@gmail.com", "123123123", false)

	seedHotel("The cozy hotel", "France", 5)
	seedHotel("Bautopa", "South Africa", 4)
	hilton := seedHotel("Hilton", "USA", 5)

	room1 := seedRoom("small", 99.10, hilton.ID)
	room2 := seedRoom("normal", 120, hilton.ID)
	room3 := seedRoom("kingsize", 320, hilton.ID)

	seedBooking(room1.ID, pedram.ID, 2, time.Now(), time.Now().AddDate(0, 0, 2))
	seedBooking(room2.ID, jack.ID, 1, time.Now(), time.Now().AddDate(0, 0, 7))
	seedBooking(room3.ID, admin.ID, 3, time.Now(), time.Now().AddDate(0, 0, 3))
}
