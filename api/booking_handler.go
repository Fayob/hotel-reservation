package api

import (
	"hotel_reservation/db"
	"hotel_reservation/types"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

// This need to be admin authorized!
func (b *BookingHandler) HandleGetBookings(ctx *fiber.Ctx) error {
	bookings, err := b.store.Booking.GetBookings(ctx.Context(), bson.M{})
	if err != nil {
		return err
	}
	return ctx.JSON(bookings)
}

// This need to be user authorized!
func (b *BookingHandler) HandleGetBooking(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	booking, err := b.store.Booking.GetBookingByID(ctx.Context(), id)
	if err != nil {
		return err
	}
	user, ok := ctx.Context().UserValue("user").(*types.User)
	if !ok {
		return err
	}
	if booking.UserID != user.ID && !user.IsAdmin {
		return ctx.Status(http.StatusUnauthorized).JSON(genericResponse{
			Type: "error",
			Msg: "not authorized",
		})
	}
	return ctx.JSON(booking)
}