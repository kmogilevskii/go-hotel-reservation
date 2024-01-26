package api

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/kmogilevskii/hotel-reservation/db"
	"github.com/kmogilevskii/hotel-reservation/errors"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			return errors.ErrUnauthorized()
		}

		claims, err := validateToken(token[0])
		if err != nil {
			return err
		}
		expiresFloat := claims["exp"].(float64)
		expires := int64(expiresFloat)
		if time.Now().Unix() > expires {
			return errors.ErrUnauthorized()
		}
		userID := claims["id"].(string)
		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			return errors.ErrUnauthorized()
		}
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("failed to parse JWT token", err)
		return nil, errors.ErrUnauthorized()
	}

	if !token.Valid {
		fmt.Println("invalid token")
		return nil, errors.ErrUnauthorized()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("invalid claims")
		return nil, errors.ErrUnauthorized()
	}
	return claims, nil
}
