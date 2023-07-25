package main

import (
	"fmt"

	"github.com/f0rSaaaa/JWTAuthenticationGO/controllers"
	"github.com/f0rSaaaa/JWTAuthenticationGO/initializers"
	"github.com/f0rSaaaa/JWTAuthenticationGO/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	fmt.Println("Hello")

	r := gin.Default()
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
