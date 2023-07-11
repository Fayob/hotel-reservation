package api

import (
	"context"
	"hotel_reservation/db"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (u *UserHandler) HandlerGetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	user, err := u.userStore.GetUserByID(context.Background(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(user)
}