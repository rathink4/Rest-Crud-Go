package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rathink4/rest-crud-go/initializers"
	"github.com/rathink4/rest-crud-go/models"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	// Get the user details
	var userDetail struct {
		Email    string
		Password string
	}

	if c.Bind(&userDetail) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read the body",
		})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(userDetail.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create the user in the database
	user := models.User{Email: userDetail.Email, Password: string(hash)}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create the user",
		})
		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{})

}

func LoginIn(c *gin.Context) {
	// Get the user details
	var userDetail struct {
		Email    string
		Password string
	}

	if c.Bind(&userDetail) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read the body",
		})
		return
	}

	// Get the user detail from the database
	var user models.User
	initializers.DB.First(&user, "email = ?", userDetail.Email)

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid email",
		})
		return
	}

	// Check if the password matches
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userDetail.Password))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid password",
		})
		return
	}

	// Send a token if it matches
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to generate the token",
		})
	}

	// Return logged in status with the token set
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{})

}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}
