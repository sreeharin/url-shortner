package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sreeharin/url-shortner/internal/models"
)

var tmpUser = models.User{
	Username: "test",
	Password: "test",
}

var inputJson, _ = json.Marshal(tmpUser)

// TestRegistration checks the registration process
// It verifies that a new user can be registered successfully
// and that the password is hashed
func TestRegistration(t *testing.T) {
	router, db := SetupTestEnvironment(t)
	handler := Handler{DB: db}

	router.POST("/register", handler.Registration)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(inputJson))

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code: %d, got: %d", http.StatusCreated, w.Code)
	}

	var user models.User
	db.First(&user)

	if user.Username != tmpUser.Username {
		t.Errorf("Expected username: %s, got: %s", tmpUser.Username, user.Username)
	}

	if user.Password == tmpUser.Password {
		t.Errorf("Password should be hashed, got: %s", user.Password)
	}

}

// TestLogin checks the login process
// It verifies that a user can log in successfully
// and that a JWT token is returned
func TestLogin(t *testing.T) {
	router, db := SetupTestEnvironment(t)
	handler := Handler{DB: db}

	router.POST("/register", handler.Registration)
	router.POST("/login", handler.Login)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(inputJson))

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code: %d, got: %d", http.StatusCreated, w.Code)
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(inputJson))

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code: %d, got: %d", http.StatusOK, w.Code)
	}

	var res map[string]any
	json.NewDecoder(w.Body).Decode(&res)

	val, exists := res["token"]
	if exists {
		if len(val.(string)) == 0 {
			t.Errorf("No token returned")
		}

	} else {
		t.Errorf("Token not found in response")
	}
}
