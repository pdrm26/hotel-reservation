package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pdrm26/hotel-reservation/api"
	"github.com/pdrm26/hotel-reservation/db"
	"github.com/pdrm26/hotel-reservation/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURL))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)
	store := &db.Store{
		User:    db.NewMongoUserStore(client),
		Hotel:   hotelStore,
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Booking: db.NewMongoBookingStore(client),
	}
	admin := fixtures.AddUser(store, "admin", "admin", true)
	pedram := fixtures.AddUser(store, "pedram", "baradarian", false)
	jack := fixtures.AddUser(store, "jack", "joe", false)
	fmt.Println("admin ->", api.CreateTokenFromUser(admin))
	fmt.Println("pedram ->", api.CreateTokenFromUser(pedram))
	fmt.Println("jack ->", api.CreateTokenFromUser(jack))

	fixtures.AddHotel(store, "The cozy hotel", "France", 5, nil)
	fixtures.AddHotel(store, "Bautopa", "South Africa", 4, nil)
	hilton := fixtures.AddHotel(store, "Hilton", "USA", 5, nil)

	room1 := fixtures.AddRoom(store, "small", true, 99.10, hilton.ID)
	room2 := fixtures.AddRoom(store, "normal", false, 120, hilton.ID)
	room3 := fixtures.AddRoom(store, "kingsize", true, 320, hilton.ID)

	fixtures.AddBooking(store, room1.ID, pedram.ID, 2, time.Now(), time.Now().AddDate(0, 0, 2))
	fixtures.AddBooking(store, room2.ID, jack.ID, 1, time.Now(), time.Now().AddDate(0, 0, 7))
	fixtures.AddBooking(store, room3.ID, admin.ID, 3, time.Now(), time.Now().AddDate(0, 0, 3))

	// For testing pagination and filtering
	// for i := 0; i < 100; i++ {
	// 	fixtures.AddHotel(store, fmt.Sprintf("The random hotel %d", i), fmt.Sprintf("Country %d", i), rand.Intn(5), nil)
	// }
}
