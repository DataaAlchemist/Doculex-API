package main 

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/Data-Alchemist/doculex-api/config"
	"github.com/Data-Alchemist/doculex-api/database"
	"github.com/Data-Alchemist/doculex-api/routes"
	"github.com/Data-Alchemist/doculex-api/middleware"
)

func main() {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatal(err)
	}
	time.Local = location

	app := fiber.New() //initialize the server

	app.Use(middleware.CORSMiddleware()) //use cors middleware from other directory

	database.ConnectDB() //connect to database
	defer database.DisconnectDB() //disconnect from database

	app.Use(logger.New()) //add logger to track http request
	
	//add middleware
	jwtMiddleware := middleware.JWTMiddleware()

	routes.SetupEndpoint(app, jwtMiddleware)

	//add setup handler for false routes
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "404 Not Found",
			"status": fiber.StatusNotFound,
		})
	})

	port := config.ConfigPort()

	err = app.Listen(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}