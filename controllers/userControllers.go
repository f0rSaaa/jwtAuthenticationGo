package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/f0rSaaaa/JWTAuthenticationGO/initializers"
	"github.com/f0rSaaaa/JWTAuthenticationGO/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	//get the email and pass of req body
	var body struct {
		Email    string
		Password string
	}

	if c.BindJSON(&body) != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read body",
		})
		return
	}

	//hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to hash password",
		})
	}

	//create the user
	user := models.User{Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to create user",
		})
	}
	//respond

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "Successful signup",
	})
}

func Login(c *gin.Context) {
	//get the email and pass of the body
	var body struct {
		Email    string
		Password string
	}

	if c.BindJSON(&body) != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read body",
		})
		return
	}

	//look up the requested user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid email or password",
		})
		return
	}

	//compare sent in pass with saved user pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid email or password",
		})
		return
	}

	//generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to generate token",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	//send it back
	c.IndentedJSON(http.StatusOK, gin.H{})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": user,
	})
}
