package api

import (
	"fmt"
	"hotel_reservation/db"
	"hotel_reservation/types"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type genericResponse struct {
	Type string `json:"type"`
	Msg string `json:"msg"`
}

func invalidCredentials(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusBadRequest).JSON(genericResponse{
		Type: "error",
		Msg: "invalid credentials",
	})
}

func (u *AuthHandler) HandleAuthentication(ctx *fiber.Ctx) error {
	var params AuthParams
	if err := ctx.BodyParser(&params); err != nil {
		return err
	}

	user, err := u.userStore.GetUserByEmail(ctx.Context(), params.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("mongo: no documents in result")
			return invalidCredentials(ctx)
		}
		return err
	}
	if !types.IsPasswordValid(user.EncryptedPassword, params.Password) {
		fmt.Println("Wrong password", err)
		return invalidCredentials(ctx)
	}
	resp := AuthResponse{
		User:  user,
		Token: CreateToken(user),
	}
	return ctx.JSON(resp)
}

func CreateToken(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 4).Unix()
	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("Failed to sign token with secret", err)
	}
	return tokenStr
}
