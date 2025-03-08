package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mjande/recipes-microservice/models"
	"github.com/mjande/recipes-microservice/utils"
)

type RecipeResponse struct {
	Message string          `json:"message"`
	Data    []models.Recipe `json:"data"`
}

// Handles getting a list of recipes.
func GetRecipes(w http.ResponseWriter, r *http.Request) {
	// Call database function to query recipes
	recipes, err := models.ListRecipes(r.Context())
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseData := RecipeResponse{
		Data: recipes,
	}

	// Encode the recipes in JSON and send as response
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(responseData)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

// Handles getting a single recipe.
func GetRecipe(w http.ResponseWriter, r *http.Request) {
	// Extract id from request
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Call database function to query recipes
	recipe, err := models.FindRecipe(r.Context(), id)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	responseData := RecipeResponse{
		Data: []models.Recipe{recipe},
	}

	// Encode the recipes in JSON and send as response
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(responseData)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

// Handles creating a recipe with ingredients
func PostRecipe(w http.ResponseWriter, r *http.Request) {
	// Decode JSON data from request
	var recipe models.Recipe
	err := json.NewDecoder(r.Body).Decode(&recipe)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Use database function to create recipe
	id, err := models.CreateRecipe(r.Context(), recipe)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get recipe from database
	recipe, err = models.FindRecipe(r.Context(), id)
	if err != nil {
		log.Println(err, "Could not find recipe ", id)
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseData := RecipeResponse{
		Message: "Recipe successfull created!",
		Data:    []models.Recipe{recipe},
	}

	// Encode recipe as JSON and send response
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(responseData)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func PatchRecipe(w http.ResponseWriter, r *http.Request) {
	// Extract id from request
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check that recipe exists
	_, err = models.FindRecipe(r.Context(), id)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	// Decode JSON data from request
	var recipe models.Recipe
	err = json.NewDecoder(r.Body).Decode(&recipe)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Use database function to create recipe
	id, err = models.UpdateRecipe(r.Context(), id, recipe)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get recipe from database
	recipe, err = models.FindRecipe(r.Context(), id)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseData := RecipeResponse{
		Message: "Recipe successfully updated!",
		Data:    []models.Recipe{recipe},
	}

	// Encode recipe as JSON and send response
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(responseData)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	// Extract id from request
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = models.DeleteRecipe(r.Context(), id)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseData := RecipeResponse{
		Message: "Recipe successfully deleted",
	}

	// Encode recipe as JSON and send response
	w.WriteHeader(http.StatusNoContent)
	err = json.NewEncoder(w).Encode(responseData)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}
