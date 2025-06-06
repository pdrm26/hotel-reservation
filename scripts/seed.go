package main

import (
	"context"
	"fmt"
	"log"

	"github.com/pdrm26/hotel-reservation/db"
	"github.com/pdrm26/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var hotel = types.Hotel{
	Name:     "hotel 1",
	Location: "NYC",
	Rooms:    []primitive.ObjectID{},
}

var rooms = []types.Room{
	{Type: types.SingleRoomType, BasePrice: 99.10},
	{Type: types.DoubleRoomType, BasePrice: 120},
	{Type: types.DeluxRoomType, BasePrice: 320},
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(context.Background()); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)

	insertedHotel, err := hotelStore.InsertHotel(context.Background(), &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		insertedRoom, err := roomStore.InsertRoom(context.Background(), &room)
		if err != nil {
			log.Fatal(err)
		}
		roomStore.HotelStore.UpdateHotel(context.Background(), bson.M{"_id": insertedHotel.ID}, bson.M{"_id": insertedRoom.ID})
		fmt.Println(insertedRoom)
	}
}
