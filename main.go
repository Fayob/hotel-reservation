package main

import (
	"context"
	"flag"
	"hotel_reservation/api"
	"hotel_reservation/db"
	"hotel_reservation/middleware"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
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
		userStore = db.NewMongoUserStore(client)
		store = &db.Store{
			Room: roomStore,
			Hotel: hotelStore,
			User: userStore,
		}
		userHandler  = api.NewUserHandler(userStore)
		hotelhandler = api.NewHotelHandler(store)
		authHandler = api.NewAuthHandler(userStore)
		app = fiber.New(config)
		apiv1 = app.Group("/api/v1", middleware.JWTAuthentication)
	)

	listenAddr := flag.String("listenAddr", ":5000", "This is the address the app will listen on")
	flag.Parse()

	app.Get("/", handleFoo)
	
	// Authentication
	app.Post("api/auth", authHandler.HandleAuthentication)

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

	app.Listen(*listenAddr)
}

func handleFoo(ctx *fiber.Ctx) error {
	return ctx.JSON(map[string]string{"msg": "Welcome to this API!!!"})
}
