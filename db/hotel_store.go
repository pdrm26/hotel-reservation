package db

import (
	"context"

	"github.com/pdrm26/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HotelStore interface {
	GetHotels(context.Context, bson.M, *PaginateFilter) ([]*types.Hotel, error)
	GetHotelByID(context.Context, primitive.ObjectID) (*types.Hotel, error)
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	UpdateHotel(context.Context, bson.M, bson.M) error
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(DBNAME).Collection("hotels"),
	}
}

func (s *MongoHotelStore) GetHotelByID(ctx context.Context, hotelId primitive.ObjectID) (*types.Hotel, error) {
	var hotel types.Hotel

	if err := s.coll.FindOne(ctx, bson.M{"_id": hotelId}).Decode(&hotel); err != nil {
		return nil, err
	}

	return &hotel, nil
}

func (s *MongoHotelStore) GetHotels(ctx context.Context, filter bson.M, p *PaginateFilter) ([]*types.Hotel, error) {
	opts := options.FindOptions{}
	opts.SetSkip((p.Page - 1) * p.Limit)
	opts.SetLimit(p.Limit)
	cursor, err := s.coll.Find(ctx, filter, &opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var hotels []*types.Hotel
	if err = cursor.All(ctx, &hotels); err != nil {
		return nil, err
	}

	return hotels, nil
}

func (s *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		hotel.ID = oid
	}

	return hotel, nil
}

func (s *MongoHotelStore) UpdateHotel(ctx context.Context, filter, update bson.M) error {
	_, err := s.coll.UpdateOne(ctx, filter, update)
	return err
}
