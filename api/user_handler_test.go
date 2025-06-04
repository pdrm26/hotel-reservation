package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/pdrm26/hotel-reservation/db"
	"github.com/pdrm26/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	db.UserStore
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	return &testdb{
		UserStore: db.NewMongoUserStore(client, db.DBNAMETEST),
	}
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}

}

func TestHandlePostUser_Success(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(db.UserStore)

	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{FirstName: "Negar", LastName: "Yekta", Email: "negar@gmail.com", Password: "123456Negar"}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, _ := app.Test(req)
	fmt.Println(resp.Status)
}
