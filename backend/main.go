package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rak-nark/sparkpass/config"
	"github.com/rak-nark/sparkpass/routes"
	"golang.org/x/crypto/bcrypt"
)

// @title SparkPass API
// @version 1.0
// @description API for subscription management
func main() {
	password := "123456"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Println("Hash correcto:", string(hash))
	// Initialize Echo
	e := echo.New()
	
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Database
	config.ConnectDB()

	// Routes
	routes.SetupRoutes(e)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
}