package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shkuran/go-library/db"
	"github.com/shkuran/go-library/routes"
)

func main() {

	db.InitDB()

	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8080")
}