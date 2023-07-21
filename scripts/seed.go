package main

import (
	"context"
	"fmt"
	"hotel_reservation/api"
	"hotel_reservation/db"
	"hotel_reservation/db/fixtures"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	var (
		mongoEndpoint = os.Getenv("MONGO_DB_URL")
		mongoDBName = os.Getenv("MONGO_DB_NAME")
	)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(mongoDBName).Drop(context.Background()); err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)
	userStore := db.NewMongoUserStore(client)
	bookStore := db.NewMongoBookingStore(client)
	store := &db.Store{
		Hotel: hotelStore,
		Room: roomStore,
		Booking: bookStore,
		User: userStore,
	}
	user := fixtures.AddUser(store, "James", "foo", false)
	fmt.Println("James ->", api.CreateToken(user))
	admin := fixtures.AddUser(store, "Admin", "admin", true)
	fmt.Println("Admin ->", api.CreateToken(admin))
	hotel := fixtures.AddHotel(store, "hotel name", "Abuja", 4, nil)
	room := fixtures.AddRoom(store, "room name", true, 199.0, hotel.ID)
	booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 4))
	fmt.Println("Booking ->", booking)

	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("hotel name %d", i)
		loc := fmt.Sprintf("hotel location %d", i)
		fixtures.AddHotel(store, name, loc, rand.Intn(5)+1, nil)
	}
}
