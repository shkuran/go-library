package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shkuran/go-library/auth"
	"github.com/shkuran/go-library/controllers"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/reservations", controllers.GetReservations)
	server.POST("/reservations", auth.Authenticate, controllers.AddReservation)
	server.POST("/reservations/:id", auth.Authenticate, controllers.ReturnBook)

	server.GET("/books", controllers.GetBooks)
	server.POST("/books", controllers.AddBook)

	server.POST("/signup", controllers.CreateUser)
	server.POST("/login", controllers.Login)
	server.GET("/users", controllers.GetUsers)
}
