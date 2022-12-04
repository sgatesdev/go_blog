package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"samgates.io/server/controllers"
	mongodb "samgates.io/server/database"
	token "samgates.io/server/utils"
)

func main() {
	router := gin.Default()

	public := router.Group("/api")

	public.POST("/signup", controllers.Register)
	public.POST("/login", controllers.Login)

	protected := router.Group("/action")
	protected.Use(token.JwtMiddleware())
	protected.GET("/user", controllers.CurrentUser)
	//protected.GET("/post", ...)

	// get values from .env or abort startup
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	mongodb.Connect()
	token.SetEnvs()
	router.Run("localhost:8080")

	// close mongo connection when router is stopped
	defer mongodb.Disconnect()
}
