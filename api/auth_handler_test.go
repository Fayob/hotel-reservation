package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hotel_reservation/db/fixtures"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestAuthenticateSuccess(t *testing.T)  {
	tdb := setup(t)
	defer tdb.teardown(t)
	insertedUser := fixtures.AddUser(tdb.Store, "james", "foo", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuthentication)

	params := AuthParams{
		Email: "james@foo.com",
		Password: "james_foo",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code of 200 but got %d", resp.StatusCode)
	}
	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Error(err)
	}
	if authResp.Token == "" {
		t.Fatalf("Expected the JWT token to be present the auth response")
	}
	// Set the encrypted password to an empty string because 
	// we didn't return in anu JSON response
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, authResp.User) {
		t.Fatalf("Expected the user to be inserted user")
	}
}

func TestAuthenticateWithWrongPasswordFailure(t *testing.T)  {
	tdb := setup(t)
	defer tdb.teardown(t)
	fixtures.AddUser(tdb.Store, "james", "foo", false)


	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuthentication)

	params := AuthParams{
		Email: "james@foo.com",
		Password: "notcorrectpassword",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code of 400 but got %d", resp.StatusCode)
	}
	var genResp genericResponse
	fmt.Println("resp.body", resp.Body)
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Error(err)
	}
	if genResp.Type != "error" {
		t.Fatalf("Expected the gen response type to be error but got %s", genResp.Type)
	}
	if genResp.Msg != "invalid credentials" {
		t.Fatalf("Expected the msg to be <invalid credentials> but got %s", genResp.Msg)
	}
}
