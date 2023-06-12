package go_web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/your-username/go_web" // Replace with the correct import path
)

func TestRegisterHandler(t *testing.T) {
	// Define the User struct
	type User struct {
		ID       int    `db:"id" json:"id"`
		Nama     string `db:"nama" json:"nama"`
		Email    string `db:"email" json:"email"`
		Password string `db:"password" json:"password"`
	}

	// Create the payload
	payload := User{
		Nama:     "John Doe",
		Email:    "johndoe@example.com",
		Password: "password123",
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	// Create an HTTP request
	req, err := http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatal(err)
	}

	// Create an HTTP response recorder to capture the response
	rr := httptest.NewRecorder()

	// Run the RegisterHandler handler by making the HTTP request
	handler := http.HandlerFunc(go_web.RegisterHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %v, but got %v", http.StatusOK, status)
	}

	// Check the response body
	expectedResponse := `{"id":1,"nama":"John Doe","email":"johndoe@example.com","password":"<hashed-password>"}`
	if rr.Body.String() != expectedResponse {
		t.Errorf("Expected response body %v, but got %v", expectedResponse, rr.Body.String())
	}
}

// Implement testing functions for other handlers such as loginHandler and logoutHandler.
// ...
