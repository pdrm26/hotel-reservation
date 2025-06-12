package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/pdrm26/hotel-reservation/db"
	"github.com/pdrm26/hotel-reservation/types"
)

func seedUser(db db.UserStore) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: "Pedram",
		LastName:  "Baradarian",
		Email:     "pedram@gmail.com",
		Password:  "12345678",
	})

	if err != nil {
		log.Fatal(err)
	}
	_, errr := db.InsertUser(context.TODO(), user)
	if errr != nil {
		log.Fatal(errr)
	}

	return user
}
func TestAuthenticateSuccess(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	insertedUser := seedUser(db.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(db.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	b, _ := json.Marshal(AuthParams{Email: "pedram@gmail.com", Password: "12345678"})
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
	seedUser(db.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(db.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	b, _ := json.Marshal(AuthParams{Email: "pedram@gmail.com", Password: "incorrectpassword"})
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
