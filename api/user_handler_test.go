package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hotel_reservation/db"
	"hotel_reservation/types"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	return &testdb{
		UserStore: db.NewMongoUserStore(client),
	}
}

func TestPostUser(t *testing.T)  {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email: "some@foo.com",
		FirstName: "James",
		LastName: "Foo",
		Password: "1234567",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	if len(user.ID) == 0 {
		t.Errorf("expecting a userID to be set")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected firstname %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected lastName %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email %s but got %s", params.Email, user.Email)
	}
	fmt.Println(resp.Status)
}

// func TestGetUser(t *testing.T)  {
// 	tdb := setup(t)
// 	defer tdb.teardown(t)

// 	app := fiber.New()
// 	userHandler := NewUserHandler(tdb.UserStore)
// 	app.Get("/", userHandler.HandlerGetUsers)

// 	http := httptest.NewRequest("GET", "/", nil)
// 	res, err := app.Test(http)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	var users []types.User
// 	json.Marshal(res.Body)
// }