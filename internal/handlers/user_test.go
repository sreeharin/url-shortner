package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sreeharin/url-shortner/internal/models"
)

// TestRegistration checks the registration process
// It verifies that a new user can be registered successfully
// and that the password is hashed
func TestRegistration(t *testing.T) {
	router, db := SetupTestEnvironment(t)
	handler := Handler{DB: db}

	router.POST("/register", handler.Registration)

	tmpUser := models.User{
		Username: "test",
		Password: "test",
	}
	inputJson, _ := json.Marshal(tmpUser)

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
