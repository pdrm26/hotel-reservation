package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/pdrm26/hotel-reservation/db/fixtures"
)

func TestAuthenticateSuccess(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	insertedUser := fixtures.AddUser(db.Store, "jack", "joe", false)

	app := fiber.New()
	authHandler := NewAuthHandler(db.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	b, _ := json.Marshal(AuthParams{Email: "jack@joe.com", Password: "jack_joe"})
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected http status to be 200 got %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}

	if authResp.Token == "" {
		t.Fatal("Expected token")
	}

	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, authResp.User) {
		t.Fatalf("Expected user to be inserted user %v", insertedUser)
	}
}

func TestAuthenticateWithWrongPassword(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	fixtures.AddUser(db.Store, "jack", "joe", false)

	app := fiber.New()
	authHandler := NewAuthHandler(db.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	b, _ := json.Marshal(AuthParams{Email: "jack@joe.com", Password: "incorrectpassword"})
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected http status to be %d got %d", http.StatusBadRequest, resp.StatusCode)
	}

	var genResp genericResp
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}

	if genResp.Message != "invalid credentials" {
		t.Fatalf("expected message like <invalid credentials> but got %s", genResp.Message)
	}
}
