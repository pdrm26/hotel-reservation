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

}

func seedHotel(name, location string) {
	var hotel = types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
	}

	var rooms = []types.Room{
		{Type: types.SingleRoomType, BasePrice: 99.10},
		{Type: types.DoubleRoomType, BasePrice: 120},
		{Type: types.DeluxRoomType, BasePrice: 320},
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

func main() {
	seedHotel("The cozy hotel", "France")
	seedHotel("Bautopa", "South Africa")
	seedHotel("Hilton", "USA")
}
