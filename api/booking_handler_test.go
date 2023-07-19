package api

import (
	"encoding/json"
	"fmt"
	"hotel_reservation/db/fixtures"
	"hotel_reservation/middleware"
	"hotel_reservation/types"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestGetBooking(t *testing.T)  {
	db := setup(t)
	defer db.teardown(t)
	var (
		nonAuthUser = fixtures.AddUser(db.Store, "jimmy", "watercooler", false)
		user = fixtures.AddUser(db.Store, "james", "foo", false)
		hotel = fixtures.AddHotel(db.Store, "bar hotel", "Lagos", 4, nil)
		room = fixtures.AddRoom(db.Store, "small", true, 4.4, hotel.ID)
		from = time.Now()
		till = from.AddDate(0,0,4)
		booking = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		app = fiber.New()
		route = app.Group("/", middleware.JWTAuthentication(db.User))
		bookingHandler = NewBookingHandler(db.Store)
	)
	route.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateToken(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	var bookingResp *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}
	if booking.ID != bookingResp.ID {
		t.Fatalf("expected %s got %s", booking.ID, bookingResp.ID)
	}
	if booking.UserID != bookingResp.UserID {
		t.Fatalf("expected %s got %s", booking.UserID, bookingResp.UserID)
	}

	// test non-authorized cannot access the booking
	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateToken(nonAuthUser))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected a non 200 status code but got %d", resp.StatusCode)
	}
}

func TestAdminGetBookings(t *testing.T)  {
	db := setup(t)
	defer db.teardown(t)
	var (
		adminUser = fixtures.AddUser(db.Store, "john", "peter", true)
		user = fixtures.AddUser(db.Store, "james", "foo", false)
		hotel = fixtures.AddHotel(db.Store, "bar hotel", "Lagos", 4, nil)
		room = fixtures.AddRoom(db.Store, "small", true, 4.4, hotel.ID)
		from = time.Now()
		till = from.AddDate(0,0,4)
		booking = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		app = fiber.New()
		admin = app.Group("/", middleware.JWTAuthentication(db.User), middleware.AdminAuth)
		bookingHandler = NewBookingHandler(db.Store)
	)
	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateToken(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response: got %d", resp.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking but got %d", len(bookings))
	}
	if booking.ID != bookings[0].ID {
		t.Fatalf("expected %s but got %s", booking.ID, bookings[0].ID)
	}
	if booking.UserID != bookings[0].UserID {
		t.Fatalf("expected %s but got %s", booking.ID, bookings[0].ID)
	}
	
	// test non-admin cannot access the bookings
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateToken(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected a non 200 status code but got %d", resp.StatusCode)
	}
	fmt.Println(bookings)
}