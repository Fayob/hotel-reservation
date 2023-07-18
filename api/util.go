package api

import (
	"fmt"
	"hotel_reservation/types"

	"github.com/gofiber/fiber/v2"
)

func getAuthUser(ctx *fiber.Ctx) (*types.User, error) {
	user, ok := ctx.Context().UserValue("user").(*types.User)
	if !ok {
		return nil, fmt.Errorf("Unauthorized")
	}
	return user, nil
}