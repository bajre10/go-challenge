package main

import (
	"fmt"
	"log"
	"net/http"
	"shortcut-challenge/config"
	"shortcut-challenge/database"
	"shortcut-challenge/router"
	"shortcut-challenge/utils"
)

// @title Swagger Example API
// @version 1.0
// @securityDefinitions.oauth2 BearerAuth
// @description Go challenge.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
func main() {
	// Load .env.docker
	config.LoadConfig()

	// Connect to DB
	database.InitDB()

	// Setup routes
	r := router.SetupRoutes()

	err := http.ListenAndServe(fmt.Sprintf(":%s", utils.GetEnv("PORT", "8080")), r)

	if err != nil {
		log.Fatal(err)
	}
}
