package main

import (
	"context"
	"fmt"
	"hotel_reservation/api"
	"hotel_reservation/db"
	"hotel_reservation/db/fixtures"
	"log"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(context.Background()); err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)
	userStore := db.NewMongoUserStore(client)
	bookStore := db.NewMongoBookingStore(client)
	store := db.Store{
		Hotel: hotelStore,
		Room: roomStore,
		Booking: bookStore,
		User: userStore,
	}
	user := fixtures.AddUser(&store, "James", "foo", false)
	fmt.Println("James ->", api.CreateToken(user))
	admin := fixtures.AddUser(&store, "Admin", "admin", true)
	fmt.Println("Admin ->", api.CreateToken(admin))
	hotel := fixtures.AddHotel(&store, "hotel name", "Abuja", 4, nil)
	room := fixtures.AddRoom(&store, "room name", true, 199.0, hotel.ID)
	booking := fixtures.AddBooking(&store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 4))
	fmt.Println("Booking ->", booking)
}
