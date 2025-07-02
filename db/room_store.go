package db

import (
	"context"
	"os"

	"github.com/pdrm26/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error)
	GetRooms(context.Context, bson.M) ([]*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection
	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	dbname := os.Getenv("MONGO_DB_NAME")
	return &MongoRoomStore{
		client:     client,
		coll:       client.Database(dbname).Collection("rooms"),
		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) GetRooms(ctx context.Context, filter bson.M) ([]*types.Room, error) {
	cursor, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rooms []*types.Room
	if err := cursor.All(ctx, &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	resp, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}

	if oid, ok := resp.InsertedID.(primitive.ObjectID); ok {
		room.ID = oid
	}

	filter := bson.M{"_id": room.HotelID}
	update := bson.M{"$push": bson.M{"rooms": room.ID}}
	if err := s.HotelStore.UpdateHotel(ctx, filter, update); err != nil {
		return nil, err
	}

	return room, nil
}
