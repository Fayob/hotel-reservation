package main

import (
	"context"
	"hotel_reservation/db"
	"hotel_reservation/types"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	hotelStore *db.MongoHotelStore
	roomStore  *db.MongoRoomStore
	userStore *db.MongoUserStore
)

func seedUser(fname, lname, email string) {
	user, err :=  types.NewUserFromParams(types.CreateUserParams{
		Email: email,
		FirstName: fname,
		LastName: lname,
		Password: "supersecurepassword",
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = userStore.CreateUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
}
func seedHotel(name, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}
	rooms := []types.Room{
		{
			Size:  "small",
			Price: 99.9,
		},
		{
			Size:  "normal",
			Price: 199.9,
		},
		{
			Size:  "kingsize",
			Price: 322.9,
		},
	}
	insertedHotel, err := hotelStore.InsertHotel(context.Background(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(context.Background(), &room)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	seedHotel("Bellucia", "France", 3)
	seedHotel("The cozy hotel", "Netherland", 4)
	seedHotel("All king", "Lagos", 1)
	seedUser("James", "foo", "foo@foo.com")
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(context.Background()); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
}
