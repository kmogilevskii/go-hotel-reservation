package api

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kmogilevskii/hotel-reservation/db"
	"github.com/kmogilevskii/hotel-reservation/errors"
	"github.com/kmogilevskii/hotel-reservation/types"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{store: store}
}

// should be admin authorized
func (b *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	var params db.BookingQueryParams
	if err := c.QueryParser(&params); err != nil {
		return errors.ErrBadRequest()
	}
	filter := db.Map{}
	if params.UserID != "" {
		filter["userID"] = params.UserID
	}
	bookings, err := b.store.Booking.GetBookings(c.Context(), filter, &params.Pagination)
	if err != nil {
		return err
	}
	resp := db.ResourceResponse{
		Results: len(bookings),
		Data:    bookings,
		Page:    params.Page,
	}
	return c.JSON(resp)
}

// should be user authorized
func (b *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	bookingID := c.Params("id")
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return fmt.Errorf("authorization problems")
	}
	booking, err := b.store.Booking.GetBookingByID(c.Context(), bookingID)
	if err != nil {
		return err
	}
	if booking.UserID != user.ID && !user.IsAdmin {
		return c.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg:  "unauthorized",
		})
	}
	return c.JSON(booking)
}

func (b *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := b.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return fmt.Errorf("authorization problems")
	}
	if booking.UserID != user.ID && !user.IsAdmin {
		return c.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg:  "unauthorized",
		})
	}
	if err := b.store.Booking.UpdateBooking(c.Context(), id); err != nil {
		return err
	}
	return c.JSON(genericResp{
		Type: "success",
		Msg:  "booking canceled",
	})
}
