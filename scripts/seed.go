package main

import (
	"context"
	"log"

	"github.com/pdrm26/hotel-reservation/db"
	"github.com/pdrm26/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx        = context.Background()
	client     *mongo.Client
	hotelStore db.HotelStore
	roomStore  db.RoomStore
	userStore  db.UserStore
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

}

func seedHotel(name, location string, rating int) {
	var hotel = types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	var rooms = []types.Room{
		{Size: "small", Price: 99.10},
		{Size: "normal", Price: 120},
		{Size: "kingsize", Price: 320},
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func seedUser(fname, lname, email, password string, isAdmin bool) {
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
	_, errr := userStore.InsertUser(ctx, user)
	if errr != nil {
		log.Fatal(errr)
	}
}

func main() {
	seedHotel("The cozy hotel", "France", 5)
	seedHotel("Bautopa", "South Africa", 4)
	seedHotel("Hilton", "USA", 5)

	seedUser("admin", "admin", "admin@admin.com", "admin123", true)
	seedUser("pedram", "baradarian", "pedram@gmail.com", "123123123", false)
	seedUser("jack", "joe", "jack@gmail.com", "123123123", false)
}
