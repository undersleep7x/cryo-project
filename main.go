package main

import(
	"log"

	"github.com/gin-gonic/gin"
	"github.com/undersleep7x/cryptowallet-v0.1/routes"
)

func main() {
	router := gin.Default()

	routes.SetupRoutes(router) //register the routes associated with this application

	log.Println("Server is now running on port 8080")
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start server on port 8080: %v", err)
	}
}