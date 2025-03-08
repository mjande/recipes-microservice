package models

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mjande/recipes-microservice/database"
	"github.com/mjande/recipes-microservice/utils"
)

type Ingredient struct {
	ID       int64   `json:"id"`
	UserId   string  `json:"userId"`
	Name     string  `json:"name"`
	RecipeID int64   `json:"recipeId"`
	Quantity float32 `json:"quantity"`
	Unit     string  `json:"unit"`
}

// Queries the database for all unique ingredients used in any recipe.
func ListIngredients(ctx context.Context) ([]string, error) {
	userId := utils.ExtractUserIDFromContext(ctx)
	query := `SELECT name FROM ingredients WHERE user_id = $1`

	// Execute query
	rows, err := database.DB.Query(ctx, query, userId)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	// Map database response onto ingredients set
	ingredientsSet := map[string]bool{}
	for rows.Next() {
		var ingredient string

		err = rows.Scan(&ingredient)
		if err != nil {
			return []string{}, err
		}

		ingredientsSet[ingredient] = true
	}

	if err = rows.Err(); err != nil {
		return []string{}, err
	}

	// Extract all unique ingredients in map into slice
	var ingredients []string
	for ingredient := range ingredientsSet {
		ingredients = append(ingredients, ingredient)
	}

	return ingredients, nil
}

// Queries all ingredients for a given recipe.
func ListIngredientsByRecipe(ctx context.Context, recipeId int64) ([]Ingredient, error) {
	query := `SELECT id, name, recipe_id, quantity, unit FROM ingredients WHERE recipe_id = $1`

	rows, err := database.DB.Query(ctx, query, recipeId)
	if err != nil && err == pgx.ErrNoRows {
		return []Ingredient{}, nil
	} else if err != nil {
		return []Ingredient{}, err
	}

	// Map database response onto recipes slice
	var ingredients []Ingredient
	for rows.Next() {
		var ingredient Ingredient

		err = rows.Scan(&ingredient.ID, &ingredient.Name, &ingredient.RecipeID, &ingredient.Quantity, &ingredient.Unit)
		if err != nil {
			return []Ingredient{}, err
		}

		ingredients = append(ingredients, ingredient)
	}

	if err = rows.Err(); err != nil {
		return []Ingredient{}, err
	}

	return ingredients, nil
}

// Queries all ingredients for any recipe in given ID list
func ListIngredientsByMultipleRecipes(ctx context.Context, recipeIds []int64) ([]Ingredient, error) {
	var ingredients []Ingredient
	for _, recipeId := range recipeIds {
		recipeIngredients, err := ListIngredientsByRecipe(ctx, recipeId)
		if err != nil {
			return ingredients, err
		}

		ingredients = append(ingredients, recipeIngredients...)
	}

	return ingredients, nil
}

// Queries a single ingredient by its name and recipe id.
func FindIngredient(ctx context.Context, name string, recipeId int64) (Ingredient, error) {
	query := `SELECT id, name, recipe_id, quantity, unit FROM ingredients WHERE name = $1 AND recipe_id = $2`

	// Query the database
	result := database.DB.QueryRow(ctx, query, name, recipeId)

	// Scan database result into recipe object
	var ingredient Ingredient
	err := result.Scan(&ingredient.ID, &ingredient.Name, &ingredient.RecipeID, &ingredient.Quantity, &ingredient.Unit)
	if err != nil {
		return Ingredient{}, err
	}

	return ingredient, nil
}

// Creates a new ingredient in the database
func CreateIngredient(ctx context.Context, ingredient Ingredient) (int64, error) {
	userId := utils.ExtractUserIDFromContext(ctx)
	query := `INSERT INTO ingredients (name, user_id, recipe_id, quantity, unit) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	row := database.DB.QueryRow(ctx, query, ingredient.Name, userId, ingredient.RecipeID, ingredient.Quantity, ingredient.Unit)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}
