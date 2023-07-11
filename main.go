package main

import (
	"context"
	"flag"
	"hotel_reservation/api"
	"hotel_reservation/db"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbUri = "mongodb://localhost:27017"
	dbName = "hotel_reservation"
	userCollection = "users"
)

var config = fiber.Config{
	ErrorHandler: func (ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

func main()  {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUri))
	if err != nil {
		log.Fatal(err)
	}
	
	// handlers initialization
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))
	
	listenAddr := flag.String("listenAddr", ":5000", "This is the address the app will listen on")
	flag.Parse()
	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")

	app.Get("/", handleFoo)
	apiv1.Get("/:id", userHandler.HandlerGetUser)

	app.Listen(*listenAddr)
}

func handleFoo(ctx *fiber.Ctx) error {
	return ctx.JSON(map[string]string{"msg": "Welcome to this API!!!"})
}


