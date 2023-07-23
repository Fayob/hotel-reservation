package api

import (
	"hotel_reservation/db"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetRooms(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := db.Map{"hotelID": oid}
	rooms, err := h.store.Room.GetRooms(ctx.Context(), filter)
	if err != nil {
		return err
	}
	return ctx.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotel(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	hotel, err := h.store.Hotel.GetHotelByID(ctx.Context(), id)
	if err != nil {
		return ErrResourceNotFound("hotel")
	}
	return ctx.JSON(hotel)
}

type ResourceResponse struct {
	Results int   `json:"results"`
	Data    any `json:"data"`
	Page    int   `json:"page"`
}

type HotelQueryParams struct {
	db.Pagination
	Rating int
}

func (h *HotelHandler) HandleGetHotels(ctx *fiber.Ctx) error {
	var query HotelQueryParams
	if err := ctx.QueryParser(&query); err != nil {
		return ErrBadRequest()
	}

	filter := db.Map{
		"rating": query.Rating,
	}
	
	hotels, err := h.store.Hotel.GetHotels(ctx.Context(), filter, &query.Pagination)
	if err != nil {
		return err
	}
	resp := ResourceResponse{
		Results: len(hotels),
		Data: hotels,
		Page: int(query.Page),
	}
	return ctx.JSON(resp)
}
