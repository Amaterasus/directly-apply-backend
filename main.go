package main

import (
	"fmt"
	"os"
	"log"

	"github.com/joho/godotenv"	

	"github.com/Amaterasus/direct-apply-backend/api/models"
	"github.com/Amaterasus/direct-apply-backend/api/controllers"
)

func main() {
	fmt.Println("Welcome to direct apply")

	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	initialMigrations()
	controllers.HandleRequests(port)
}

func initialMigrations() {
	models.InitialUserMigration()
}