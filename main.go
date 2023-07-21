package main

import (
	"context"
	"flag"
	"hotel_reservation/api"
	"hotel_reservation/db"
	"hotel_reservation/middleware"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		if apiError, ok := err.(api.Error); ok {
			return ctx.Status(apiError.Code).JSON(apiError)
		}
		apiError := api.NewError(http.StatusInternalServerError, err.Error())
		return ctx.Status(apiError.Code).JSON(apiError)
	},
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	// handlers initialization
	var (
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		userStore    = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			Room:    roomStore,
			Hotel:   hotelStore,
			User:    userStore,
			Booking: bookingStore,
		}
		userHandler    = api.NewUserHandler(userStore)
		hotelhandler   = api.NewHotelHandler(store)
		authHandler    = api.NewAuthHandler(userStore)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		app            = fiber.New(config)
		apiv1          = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
		admin          = apiv1.Group("/admin", middleware.AdminAuth)
	)

	listenAddr := flag.String("listenAddr", ":5000", "This is the address the app will listen on")
	flag.Parse()

	app.Get("/", handleFoo)

	// Authentication
	app.Post("api/auth", authHandler.HandleAuthentication)

	// Admin route
	apiv1.Post("/admin", middleware.AdminAuth)

	// User Handlers
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandlerGetUsers)
	apiv1.Get("/user/:id", userHandler.HandlerGetUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandleUpdateUser)

	// Hotel Handlers
	apiv1.Get("/hotel", hotelhandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelhandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelhandler.HandleGetRooms)

	// Rooms Handler
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)
	apiv1.Get("/room", roomHandler.HandleGetRooms)

	// Booking Handler
	admin.Get("/booking", bookingHandler.HandleGetBookings)
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)
	app.Listen(*listenAddr)
}

func handleFoo(ctx *fiber.Ctx) error {
	return ctx.JSON(map[string]string{"msg": "Welcome to this API!!!"})
}
