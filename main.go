package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/mjande/recipes-microservice/database"
)

func main() {
	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatal(err)
	}

	err = database.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.DB.Close()

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{os.Getenv("CLIENT_URL")},
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE"},
	}))

	// Register routes here

	log.Printf("Recipes service listening on port %s", os.Getenv("PORT"))
	err = http.ListenAndServe(":"+os.Getenv("PORT"), router)
	if err != nil {
		log.Fatal(err)
	}
}
