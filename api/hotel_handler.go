package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kmogilevskii/hotel-reservation/db"
	"github.com/kmogilevskii/hotel-reservation/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(s *db.Store) *HotelHandler {
	return &HotelHandler{store: s}
}

func (hh *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	//?page=1&limit=10
	var params db.HotelQueryParams
	if err := c.QueryParser(&params); err != nil {
		return errors.ErrBadRequest()
	}
	filter := db.Map{
		//"rating": params.Rating,
	}
	hotels, err := hh.store.Hotel.GetHotels(c.Context(), filter, &params.Pagination)
	if err != nil {
		return errors.ErrBadRequest()
	}
	resp := db.ResourceResponse{
		Results: len(hotels),
		Data:    hotels,
		Page:    params.Page,
	}
	return c.JSON(resp)
}

func (hh *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	hotelID := c.Params("id")
	hotels, err := hh.store.Hotel.GetHotelByID(c.Context(), hotelID)
	if err != nil {
		return errors.ErrInvalidID()
	}
	return c.JSON(hotels)
}

func (hh *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	var params db.RoomQueryParams
	if err := c.QueryParser(&params); err != nil {
		return errors.ErrBadRequest()
	}
	hotelID := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		return errors.ErrInvalidID()
	}
	filter := db.Map{
		"hotelID": oid,
	}
	if params.Size != "" {
		filter["size"] = params.Size
	}

	rooms, err := hh.store.Room.GetRooms(c.Context(), filter, &params.Pagination)
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
