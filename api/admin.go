package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kmogilevskii/hotel-reservation/errors"
	"github.com/kmogilevskii/hotel-reservation/types"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return errors.ErrUnauthorized()
	}
	if !user.IsAdmin {
		return errors.ErrUnauthorized()
	}
	return c.Next()
}
