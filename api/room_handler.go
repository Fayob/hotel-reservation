package api

import (
	"fmt"
	"hotel_reservation/db"
	"hotel_reservation/types"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
	FromDate time.Time `json:"fromDate"`
	TillDate time.Time `json:"tillDate"`
	NumPersons int `json:"numPerson"`
}

func (p BookRoomParams) validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.TillDate) {
		return fmt.Errorf("cannot book a room in the past")
	}
	return nil
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (r *RoomHandler) HandleGetRooms(ctx *fiber.Ctx) error {
	rooms, err := r.store.Room.GetRooms(ctx.Context(), bson.M{})
	if err != nil {
		return nil
	}
	return ctx.JSON(rooms)
}

func (r *RoomHandler) HandleBookRoom(ctx *fiber.Ctx) error {
	var param BookRoomParams
	if err := ctx.BodyParser(&param); err != nil {
		return err
	}
	if err := param.validate(); err != nil {
		return nil
	}
	roomID, err := primitive.ObjectIDFromHex(ctx.Params("id"))
	if err != nil {
		return err
	}
	user, ok := ctx.Context().Value("user").(*types.User)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(genericResponse{
			Type: "error",
			Msg: "internal sever error",
		})
	}

	ok, err = r.isRoomAvailableForBooking(ctx, roomID, param)
	if err != nil {
		return err
	}
	if !ok {
		return ctx.Status(http.StatusBadRequest).JSON(genericResponse{
			Type: "error",
			Msg: fmt.Sprintf("room %s already booked", ctx.Params("id")),
		})
	}

	booking := types.Booking{
		UserID: user.ID,
		RoomID: roomID,
		NumPerson: param.NumPersons,
		FromDate: param.FromDate,
		TillDate: param.TillDate,
	}

	inserted, err := r.store.Booking.InsertBooking(ctx.Context(), &booking)
	if err !=nil {
		return err
	}

	return ctx.JSON(inserted)
}

func (r RoomHandler) isRoomAvailableForBooking(ctx *fiber.Ctx, roomID primitive.ObjectID, param BookRoomParams) (bool, error) {
	where := bson.M{
		"roomID": roomID,
		"fromDate": bson.M{
			"$gte": param.FromDate,
		},
		"tillDate": bson.M{
			"$gte": param.TillDate,
		},
	}

	bookings, err := r.store.Booking.GetBookings(ctx.Context(), where)
	if err != nil {
		return false, err
	}
	ok := len(bookings) == 0
	return ok, nil
}