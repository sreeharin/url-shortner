package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/sreeharin/url-shortner/internal/models"
)

// Registration handles user registration
// It accepts a JSON payload with username and password
// It hashes the password and stores the user in the database
// It returns a 201 Created status on success
func (h *Handler) Registration(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bytes, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(bytes)

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created"})
}
