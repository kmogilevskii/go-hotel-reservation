package api

import (
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

func TestHandleGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	var (
		admin          = fixtures.AddUser(tdb.store.User, "admin", "admin", true)
		user           = fixtures.AddUser(tdb.store.User, "foo", "bar", false)
		hotel          = fixtures.AddHotel(tdb.store.Hotel, "Hilton", "New York", 5)
		room           = fixtures.AddRoom(tdb.store.Room, hotel.ID, "Single", 99.99)
		booking        = fixtures.AddBooking(tdb.store.Booking, user.ID, room.ID, time.Now().AddDate(0, 0, 2), time.Now().AddDate(0, 0, 8), 2)
		app            = fiber.New(fiber.Config{ErrorHandler: errors.ErrorHandler})
		bookingHandler = NewBookingHandler(tdb.store)
	)
	jwt := app.Group("/", JWTAuthentication(tdb.store.User), AdminAuth)
	jwt.Get("/", bookingHandler.HandleGetBookings)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Token", CreateTokenFromUser(admin))
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}

	var bookings []*types.Booking
	json.NewDecoder(resp.Body).Decode(&bookings)
	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking, got %d", len(bookings))
	}

	if bookings[0].ID != booking.ID {
		t.Fatalf("got wrong booking, expected to be %v, got %v", booking, bookings[0])
	}

	if bookings[0].UserID != booking.UserID {
		t.Fatalf("got wrong booking, expected to be %v, got %v", booking, bookings[0])
	}

	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Token", CreateTokenFromUser(user))
	resp, _ = app.Test(req)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status code 401, got %d", resp.StatusCode)
	}

}

func TestGetBooking(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	var (
		admin          = fixtures.AddUser(tdb.store.User, "admin", "admin", true)
		user           = fixtures.AddUser(tdb.store.User, "foo", "bar", false)
		non_auth_user  = fixtures.AddUser(tdb.store.User, "bar", "foo", false)
		hotel          = fixtures.AddHotel(tdb.store.Hotel, "Hilton", "New York", 5)
		room           = fixtures.AddRoom(tdb.store.Room, hotel.ID, "Single", 99.99)
		booking        = fixtures.AddBooking(tdb.store.Booking, user.ID, room.ID, time.Now().AddDate(0, 0, 2), time.Now().AddDate(0, 0, 8), 2)
		app            = fiber.New(fiber.Config{ErrorHandler: errors.ErrorHandler})
		bookingHandler = NewBookingHandler(tdb.store)
	)
	jwt := app.Group("/", JWTAuthentication(tdb.store.User))
	jwt.Get("/:id", bookingHandler.HandleGetBooking)

	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Token", CreateTokenFromUser(user))
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}

	var respBooking *types.Booking
	json.NewDecoder(resp.Body).Decode(&respBooking)

	if respBooking.ID != booking.ID {
		t.Fatalf("got wrong booking, expected to be %v, got %v", booking, respBooking)
	}

	if respBooking.UserID != booking.UserID {
		t.Fatalf("got wrong booking, expected to be %v, got %v", booking, respBooking)
	}

	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Token", CreateTokenFromUser(admin))
	resp, _ = app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}

	json.NewDecoder(resp.Body).Decode(&respBooking)

	if respBooking.ID != booking.ID {
		t.Fatalf("got wrong booking, expected to be %v, got %v", booking, respBooking)
	}

	if respBooking.UserID != booking.UserID {
		t.Fatalf("got wrong booking, expected to be %v, got %v", booking, respBooking)
	}

	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Token", CreateTokenFromUser(non_auth_user))
	resp, _ = app.Test(req)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status code 401, got %d", resp.StatusCode)
	}
}

func TestCancelBooking(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	var (
		//admin          = fixtures.AddUser(tdb.store.User, "admin", "admin", true)
		user = fixtures.AddUser(tdb.store.User, "foo", "bar", false)
		//non_auth_user  = fixtures.AddUser(tdb.store.User, "bar", "foo", false)
		hotel          = fixtures.AddHotel(tdb.store.Hotel, "Hilton", "New York", 5)
		room           = fixtures.AddRoom(tdb.store.Room, hotel.ID, "Single", 99.99)
		booking        = fixtures.AddBooking(tdb.store.Booking, user.ID, room.ID, time.Now().AddDate(0, 0, 2), time.Now().AddDate(0, 0, 8), 2)
		app            = fiber.New(fiber.Config{ErrorHandler: errors.ErrorHandler})
		bookingHandler = NewBookingHandler(tdb.store)
	)
	jwt := app.Group("/", JWTAuthentication(tdb.store.User))
	jwt.Delete("/:id", bookingHandler.HandleCancelBooking)

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Token", CreateTokenFromUser(user))
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}

	var genResp *genericResp
	json.NewDecoder(resp.Body).Decode(&genResp)
	if genResp.Type != "success" || genResp.Msg != "booking canceled" {
		t.Fatalf("expected success response, but got %v", genResp)
	}

}
