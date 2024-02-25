package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shkuran/go-library/db"
	"github.com/shkuran/go-library/reservation"
	"github.com/shkuran/go-library/routes"
)

func main() {

	db.InitDB()
	db.CreateTables()

	server := gin.Default()

	reservation_db := reservation.NewMySQL(db.DB)
	reservation_repo := reservation.NewRepository(reservation_db)
	reservation_controller := reservation.NewReservationController(reservation_repo)

	routes.RegisterRoutes(server, reservation_controller)

	server.Run(":8080")
}
