package main

import (
	routes "github.com/pedromessetti/auth_project/routes"
	"os"
	"github.com/gin-gonic/gin"
)

// Sets up a server using the Gin framework, defines some routes, and runs the server on the specified port.
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{"sucess": "Access granted for api-1"})
	})

	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"sucess": "Access granted for api-2"})
	})

	router.Run(":" + port)
}