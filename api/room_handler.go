package api

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kmogilevskii/hotel-reservation/db"
	"github.com/kmogilevskii/hotel-reservation/errors"
	"github.com/kmogilevskii/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{store: store}
}

func (r *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	var params db.RoomQueryParams
	if err := c.QueryParser(&params); err != nil {
		return errors.ErrBadRequest()
	}
	filter := db.Map{}
	if params.Size != "" {
		filter["size"] = params.Size
	}

	rooms, err := r.store.Room.GetRooms(c.Context(), filter, &params.Pagination)
	if err != nil {
		return err
	}
	resp := db.ResourceResponse{
		Results: len(rooms),
		Data:    rooms,
		Page:    params.Page,
	}
	return c.Status(http.StatusOK).JSON(resp)
}

func (r *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params types.BookParams
	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err != nil {
		return err
	}
	if err := params.Validate(); err != nil {
		return errors.ErrBadRequest()
	}
	ok, err := r.isRoomAvailableForBooking(c.Context(), c.Params("id"), params)
	if err != nil {
		return err
	}
	if !ok {
		return errors.ErrAlreadyBooked()
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return err
	}

	booking := types.Booking{
		ID:         primitive.NewObjectID(),
		RoomID:     roomID,
		UserID:     user.ID,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
		NumPersons: params.NumPersons,
		Canceled:   false,
	}

	insertedBooking, err := r.store.Booking.Insert(c.Context(), &booking)
	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(insertedBooking)
}

func (r *RoomHandler) isRoomAvailableForBooking(ctx context.Context, roomID string, params types.BookParams) (bool, error) {
	filter := db.Map{
		"roomID": roomID,
	}

	bookings, err := r.store.Booking.GetBookings(ctx, filter, &db.Pagination{Page: 1, Limit: 100})
	if err != nil {
		return false, err
	}
	ok := true
	for _, booking := range bookings {
		if params.FromDate.Before(booking.FromDate) && params.TillDate.After(booking.FromDate) {
			ok = false
			break
		}
		if params.FromDate.Before(booking.TillDate) && params.TillDate.After(booking.TillDate) {
			ok = false
			break
		}
		if params.FromDate.After(booking.FromDate) && params.TillDate.Before(booking.TillDate) {
			ok = false
			break
		}
		if params.FromDate.Equal(booking.FromDate) && params.TillDate.Equal(booking.TillDate) {
			ok = false
			break
		}
	}

	return ok, nil
}
