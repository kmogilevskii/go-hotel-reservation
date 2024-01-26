package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kmogilevskii/hotel-reservation/db/fixtures"
)

func TestAuthenticateWithWrongPasswordFailure(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	fixtures.AddUser(tdb.store.User, "foo", "bar", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.store.User)
	app.Post("/", authHandler.HandleAuthenticate)
	params := AuthParams{
		Email:    "foo@bar.com",
		Password: "wrongpassword",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status code 400, got %d", resp.StatusCode)
	}

	var genResp genericResp
	json.NewDecoder(resp.Body).Decode(&genResp)
	if genResp.Type != "error" {
		t.Fatalf("expected type to be error, got %s", genResp.Type)
	}
	if genResp.Msg != "invalid credentials" {
		t.Fatalf("expected msg to be invalid credentials, got %s", genResp.Msg)
	}

}

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	insertedUser := fixtures.AddUser(tdb.store.User, "foo", "bar", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.store.User)
	app.Post("/", authHandler.HandleAuthenticate)
	params := AuthParams{
		Email:    "foo@bar.com",
		Password: "foo_bar",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}
	var authResp AuthResponse
	json.NewDecoder(resp.Body).Decode(&authResp)

	if authResp.Token == "" {
		t.Fatalf("expected token to be present")
	}

	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(authResp.User, insertedUser) {
		t.Fatalf("expected user to be %#v, got %#v", insertedUser, authResp.User)
	}
}
