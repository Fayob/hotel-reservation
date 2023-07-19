package middleware

import (
	"fmt"
	"hotel_reservation/db"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Getting token through Request Header Authorization (Bearer token)
		// AuthorizationToken, ok := ctx.GetReqHeaders()["Authorization"]
		// if !ok {
		// 		return fmt.Errorf("Unauthorized")
		// 	}
		// fmt.Println(strings.Split(AuthorizationToken, " ")[1])

		token, ok := ctx.GetReqHeaders()["X-Api-Token"]
		if !ok {
			fmt.Println("Token not present in the headers")
			return fmt.Errorf("Unauthorized")
		}
		claims, err := validateToken(token)
		if err != nil {
			return err
		}
		expires := claims["expires"].(float64)
		expiresInt := int64(expires)
		// Check token expiration
		if time.Now().Unix() > expiresInt {
			return fmt.Errorf("token expired")
		}
		userID := claims["id"].(string)
		user, err := userStore.GetUserByID(ctx.Context(), userID)
		if err != nil {
			return fmt.Errorf("unauthorized")
		}
		// Set the current authenticated user to he context
		ctx.Context().SetUserValue("user", user)
		return ctx.Next()
	}
}

func validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("Unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse JWT token:", err)
		return nil, fmt.Errorf("unauthorized")
	}
	if !token.Valid {
		fmt.Println("invalid token")
		return nil, fmt.Errorf("unauthorized")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	return claims, nil
}
