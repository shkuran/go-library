package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/shkuran/go-library/db"
	"github.com/shkuran/go-library/reservation"
	"github.com/shkuran/go-library/routes"
	"github.com/shkuran/go-library/user"
)

func main() {

	mysql, err := db.InitDB()
	if err != nil {
		log.Fatalln("Failed to initialize database")
		return
	}

	db.CreateTables(mysql)

	server := gin.Default()

	reservation_repo := reservation.NewMySQLRepository(mysql)
	reservation_controller := reservation.NewReservationController(reservation_repo)

	user_repo := user.NewMySQLRepository(mysql)
	user_controller := user.NewUserController(user_repo)

	routes.RegisterRoutes(server, reservation_controller, user_controller)

	server.Run(":8080")
}
