package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rathink4/rest-crud-go/controllers"
	"github.com/rathink4/rest-crud-go/initializers"
	"github.com/rathink4/rest-crud-go/middlewares"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.LoginIn)
	r.GET("/validate", middlewares.AuthRequired, controllers.Validate)
	r.Run()
}
