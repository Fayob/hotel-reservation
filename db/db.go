package db

const (
	DBURI      = "mongodb://localhost:27017"
	DBNAME     = "hotel_reservation"
	TESTDBNAME = "hotel_reservation_test"
)

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}
