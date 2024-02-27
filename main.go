package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/shkuran/go-library/db"
	"github.com/shkuran/go-library/routes"
)

func main() {

	// varDb, err := db.InitMySQLDB()
	varDb, err := db.InitPostgresDB()

	if err != nil {
		log.Fatal(err)
		return
	}

	db.CreateTables(varDb)

	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8080")
}
