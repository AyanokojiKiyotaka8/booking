package main

import (
	"context"
	"log"
	"os"

	"github.com/AyanokojiKiyotaka8/booking.git/api"
	"github.com/AyanokojiKiyotaka8/booking.git/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Database clients
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	// db
	userStore := db.NewMongoUserStore(client, db.DBNAME)
	hotelStore := db.NewMongoHotelStore(client, db.DBNAME)
	roomStore := db.NewMongoRoomStore(client, db.DBNAME, hotelStore)
	bookingStore := db.NewMongoBookingStore(client, db.DBNAME)
	store := db.Store{
		User:    userStore,
		Hotel:   hotelStore,
		Room:    roomStore,
		Booking: bookingStore,
	}

	// Handlers
	userHandler := api.NewUserHandler(userStore)
	hotelHandler := api.NewHotelHandler(&store)
	roomHandler := api.NewRoomHandler(&store)
	authHandler := api.NewAuthHandler(userStore)
	bookingHandler := api.NewBookingHandler(&store)

	// App and API's
	app := fiber.New(fiber.Config{ErrorHandler: api.ErrorHandler})
	apiv1 := app.Group("/api/v1", api.JWTAuthentication(userStore))
	auth := app.Group("/api")
	admin := apiv1.Group("/admin", api.AdminAuth)

	// Auth
	auth.Post("/auth", authHandler.HandleAuth)

	// User API's
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)

	// Hotel API's
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// Room API's
	apiv1.Get("/room", roomHandler.HandleGetRooms)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	// Booking API with user auth required
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	// Booking API with admin auth required
	admin.Get("/booking", bookingHandler.HandleGetBookings)

	app.Listen(os.Getenv("HTTP_LISTEN_ADDRESS"))
}

func init() {
	db.LoadConfig()
}
