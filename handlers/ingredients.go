package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mjande/recipes-microservice/models"
	"github.com/mjande/recipes-microservice/utils"
)

type IngredientNamesResponse struct {
	Message string   `json:"message"`
	Data    []string `json:"data"`
}

type IngredientsResponse struct {
	Message string              `json:"message"`
	Data    []models.Ingredient `json:"data"`
}

// Handles getting a unique list of ingredients used in other recipes.
func GetIngredients(w http.ResponseWriter, r *http.Request) {
	// Call database function to query ingredients
	ingredients, err := models.ListIngredients(r.Context())
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseData := IngredientNamesResponse{
		Data: ingredients,
	}

	// Encode the ingredients in JSON and send as response
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(responseData)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

// Handles getting a list of ingredients used the given set of recipes
func GetIngredientsByMultipleRecipes(w http.ResponseWriter, r *http.Request) {
	var recipeIds []int64
	err := json.NewDecoder(r.Body).Decode(&recipeIds)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Call database function to query ingredients
	ingredients, err := models.ListIngredientsByMultipleRecipes(r.Context(), recipeIds)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseData := IngredientsResponse{
		Data: ingredients,
	}

	// Encode the ingredients in JSON and send as response
	err = json.NewEncoder(w).Encode(responseData)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}
