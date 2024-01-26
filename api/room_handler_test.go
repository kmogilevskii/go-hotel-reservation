package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kmogilevskii/hotel-reservation/db/fixtures"
	"github.com/kmogilevskii/hotel-reservation/errors"
	"github.com/kmogilevskii/hotel-reservation/types"
)

func TestHandleBookRoom(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	var (
		user        = fixtures.AddUser(tdb.store.User, "foo", "bar", false)
		hotel       = fixtures.AddHotel(tdb.store.Hotel, "Hilton", "New York", 5)
		room        = fixtures.AddRoom(tdb.store.Room, hotel.ID, "Single", 99.99)
		booking     = fixtures.AddBooking(tdb.store.Booking, user.ID, room.ID, time.Now().AddDate(0, 0, 2), time.Now().AddDate(0, 0, 8), 2)
		app         = fiber.New(fiber.Config{ErrorHandler: errors.ErrorHandler})
		roomHandler = NewRoomHandler(tdb.store)
	)
	book := app.Group("/", JWTAuthentication(tdb.store.User))
	book.Post("/room/:id/book", roomHandler.HandleBookRoom)

	params := types.BookParams{
		FromDate:   time.Now().AddDate(0, 0, 9),
		TillDate:   time.Now().AddDate(0, 0, 10),
		NumPersons: 2,
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", fmt.Sprintf("/room/%s/book", room.ID.Hex()), bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Token", CreateTokenFromUser(user))
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status code 201, got %d", resp.StatusCode)
	}
	var respBooking *types.Booking
	json.NewDecoder(resp.Body).Decode(&respBooking)

	if !params.FromDate.Equal(respBooking.FromDate) || !params.TillDate.Equal(respBooking.TillDate) {
		t.Fatalf("booking dates do not match; expected: %s-%s, got %s-%s", params.FromDate, params.TillDate, respBooking.FromDate, respBooking.TillDate)
	}

	if booking.TillDate.After(respBooking.FromDate) && booking.TillDate.Before(respBooking.TillDate) {
		t.Fatalf("bookings overlap")
	}

	// testing overlapping dates
	params = types.BookParams{
		FromDate:   time.Now().AddDate(0, 0, 1),
		TillDate:   time.Now().AddDate(0, 0, 4),
		NumPersons: 2,
	}
	b, _ = json.Marshal(params)

	req = httptest.NewRequest("POST", fmt.Sprintf("/room/%s/book", room.ID.Hex()), bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Token", CreateTokenFromUser(user))
	resp, _ = app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status code 400, got %d", resp.StatusCode)
	}
}
