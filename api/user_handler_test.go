package api

import (
	"context"
	"log"
	"testing"

	"github.com/pdrm26/hotel-reservation/db"
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

	// userParams := types.CreateUserParams{FirstName: "Pedram", LastName: "BR", Email: "p@gmail.com", Password: "123"}

	// user, err := types.NewUserFromParams(userParams)
	// if err != nil {
	// 	t.Fail(err)
	// }

	// if err := testStore.InsertUser(context.Context(), user); err != nil {
	// 	t.Fail(err)
	// }

}
