package middleware

import (
	"fmt"
	"hotel_reservation/types"

	"github.com/gofiber/fiber/v2"
)

func AdminAuth(ctx *fiber.Ctx) error {
	user, ok := ctx.Context().UserValue("user").(*types.User)
	if !ok {
		return fmt.Errorf("not authorized")
	}
	if !user.IsAdmin {
		return fmt.Errorf("not authorized")
	}
	return ctx.Next()
}