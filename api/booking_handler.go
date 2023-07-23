package api

import (
	"hotel_reservation/db"
	"net/http"

	"github.com/gofiber/fiber/v2"
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
	bookings, err := b.store.Booking.GetBookings(ctx.Context(), db.Map{})
	if err != nil {
		return err
	}
	return ctx.JSON(bookings)
}

func (b *BookingHandler) HandleCancelBooking(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	booking, err := b.store.Booking.GetBookingByID(ctx.Context(), id)
	if err != nil {
		return err
	}
	user, err := getAuthUser(ctx)
	if err != nil {
		return err
	}
	if booking.UserID != user.ID && !user.IsAdmin {
		return ctx.Status(http.StatusUnauthorized).JSON(genericResponse{
			Type: "error",
			Msg: "not authorized",
		})
	}
	if err := b.store.Booking.UpdateBooking(ctx.Context(), ctx.Params("id"), db.Map{"canceled": true}); err != nil {
		return err
	}
	return ctx.JSON(genericResponse{Type: "msg", Msg: "Book Canceled"})
}

// This need to be both admin and user authorized!
func (b *BookingHandler) HandleGetBooking(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	booking, err := b.store.Booking.GetBookingByID(ctx.Context(), id)
	if err != nil {
		return err
	}
	user, err := getAuthUser(ctx)
	if err != nil {
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