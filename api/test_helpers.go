package api

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/pdrm26/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	client *mongo.Client
	*db.Store
}

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func setup(t *testing.T) *testdb {
	dburi := os.Getenv("MONGO_DB_URL_TEST")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	store := &db.Store{
		User:    db.NewMongoUserStore(client),
		Hotel:   hotelStore,
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Booking: db.NewMongoBookingStore(client),
	}
	return &testdb{
		client: client,
		Store:  store,
	}
}

func (tdb *testdb) teardown(t *testing.T) {
	mongoDBName := os.Getenv("MONGO_DB_NAME")
	if err := tdb.client.Database(mongoDBName).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}
