package api

import (
	"errors"
	"hotel_reservation/db"
	"hotel_reservation/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (u *UserHandler) HandleUpdateUser(ctx *fiber.Ctx) error {
	var (
		params types.UpdateUserParams
		userID = ctx.Params("id")
	)
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	if err := ctx.BodyParser(&params); err != nil {
		return err
	}
	filter := bson.M{"_id": oid}
	if err := u.userStore.UpdateUser(ctx.Context(), filter, params); err != nil {
		return err
	}
	return ctx.JSON(map[string]string{"response": "User with userID " + userID + " was Updated successfully"})
}

func (u *UserHandler) HandleDeleteUser(ctx *fiber.Ctx) error {
	userID := ctx.Params("id")
	err := u.userStore.DeleteUser(ctx.Context(), userID)
	if err != nil {
		return err
	}
	return ctx.JSON(map[string]string{"response": "User with userID " + userID + " was deleted successfully"})
}

func (u *UserHandler) HandlePostUser(ctx *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := ctx.BodyParser(&params); err != nil {
		return err
	}
	if err := params.Validate(); len(err) > 0 {
		return ctx.JSON(err)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}
	insertedUser, err := u.userStore.CreateUser(ctx.Context(), user)
	if err != nil {
		return err
	}
	return ctx.JSON(insertedUser)
}

func (u *UserHandler) HandlerGetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	user, err := u.userStore.GetUserByID(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ctx.JSON(map[string]string{"error": "Not Found"})
		}
		return err
	}

	return ctx.JSON(user)
}

func (u *UserHandler) HandlerGetUsers(ctx *fiber.Ctx) error {
	user, err := u.userStore.GetUsers(ctx.Context())
	if err != nil {
		return ErrResourceNotFound("user")
	}

	return ctx.JSON(user)
}
