package main

import (
	"github.com/Qwerci/sps_backend/database"
	// "github.com/Qwerci/sps_backend/middlewares"
	// routes "github.com/Qwerci/sps_backend/routes"
	"github.com/Qwerci/sps_backend/controllers"
	"github.com/Qwerci/sps_backend/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

)

func loadDatabase() {
	database.LoadEnv()
	database.Connect()
	models.SyncDatabase()
}

func main() {
	loadDatabase()



	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(cors.Default())


	// Route
	r.GET("/ws", controllers.HandleConnections)
	r.GET("/registered-contacts", controllers.GetRegisteredContacts)
	r.POST("/import-contacts", controllers.ImportContacts)
	r.POST("/register", controllers.RegisterUser)
	r.POST("/login", controllers.LoginUser)
	r.POST("/initiate-chat", controllers.InitiateChat)


	r.Run(":8084")

}