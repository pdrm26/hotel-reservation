package main

import (
	"context"
	"flag"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/pdrm26/hotel-reservation/api"
	"github.com/pdrm26/hotel-reservation/api/middleware"
	"github.com/pdrm26/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	// uri := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	var (
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		userStore    = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			User:    userStore,
			Hotel:   hotelStore,
			Room:    roomStore,
			Booking: bookingStore,
		}
		hotelHandler = api.NewHotelHandler(store)
		userHandler  = api.NewUserHandler(userStore)
		authHandler  = api.NewAuthHandler(userStore)
		roomHandler  = api.NewRoomHandler(store)
		app          = fiber.New(config)
		auth         = app.Group("/api")
		apiv1        = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
	)

	// auth
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// User handler
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)

	// Hotel handler
	apiv1.Get("/hotel", hotelHandler.HanldeGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// Room handler
	apiv1.Get("/room", roomHandler.HandleGetRooms)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	app.Listen(*listenAddr)
}
