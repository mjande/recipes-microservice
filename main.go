package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mjande/recipes-microservice/database"
	"github.com/mjande/recipes-microservice/handlers"
)

func main() {
	// Connect to database
	err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.DB.Close()

	// Create new router
	router := chi.NewRouter()

	// Middleware
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{os.Getenv("CLIENT_URL")},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	router.Use(middleware.Logger)

	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("SECRET_KEY")), nil)

	router.Use(jwtauth.Verifier(tokenAuth))

	// Routes
	router.Route("/ingredients", func(r chi.Router) {
		r.Get("/", handlers.GetIngredients)
	})

	router.Route("/recipes", func(r chi.Router) {
		r.Get("/", handlers.GetRecipes)
		r.Get("/{id}", handlers.GetRecipe)
		r.Post("/", handlers.PostRecipe)
		r.Patch("/{id}", handlers.PatchRecipe)
		r.Delete("/{id}", handlers.DeleteRecipe)
	})

	// Start server
	log.Printf("Recipes service listening on port %s", os.Getenv("PORT"))
	err = http.ListenAndServe(":"+os.Getenv("PORT"), router)
	if err != nil {
		log.Fatal(err)
	}
}
