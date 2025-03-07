package models

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/mjande/recipes-microservice/database"
	"github.com/mjande/recipes-microservice/utils"
)

type Recipe struct {
	ID           int64        `json:"id"`
	Name         string       `json:"name"`
	CookingTime  string       `json:"cookingTime"`
	Description  string       `json:"description"`
	Instructions string       `json:"instructions"`
	Ingredients  []Ingredient `json:"ingredients"`
	Tags         []string     `json:"tags"`
	UserID       int64        `json:"userId"`
}

// Queries the database for all recipes (while only loading basic
// data for index page)
func ListRecipes(ctx context.Context) ([]Recipe, error) {
	userId := utils.ExtractUserIDFromContext(ctx)

	query := `SELECT id, name, cooking_time, description FROM recipes WHERE user_id = $1`

	// Execute query
	rows, err := database.DB.Query(ctx, query, userId)
	if err != nil {
		return []Recipe{}, err
	}
	defer rows.Close()

	// Map database response onto recipes slice
	recipes := []Recipe{}
	for rows.Next() {
		var recipe Recipe

		err = rows.Scan(&recipe.ID, &recipe.Name, &recipe.CookingTime, &recipe.Description)
		if err != nil {
			return []Recipe{}, err
		}

		// Get all tags for this recipe
		tags, err := FindTagsByRecipe(ctx, recipe.ID)
		if err != nil {
			return []Recipe{}, err
		}

		// Add the tag name to this recipe
		for _, tag := range tags {
			recipe.Tags = append(recipe.Tags, tag.Name)
		}

		recipes = append(recipes, recipe)
	}

	if err = rows.Err(); err != nil {
		return []Recipe{}, err
	}

	return recipes, nil
}

// Queries the database for an ingredient that has the given id.
func FindRecipe(ctx context.Context, id int64) (Recipe, error) {
	query := `SELECT id, name, cooking_time, description, instructions FROM recipes WHERE id = $1`

	// Query the database
	result := database.DB.QueryRow(ctx, query, id)

	// Scan database result into recipe object
	var recipe Recipe
	err := result.Scan(&recipe.ID, &recipe.Name, &recipe.CookingTime, &recipe.Description, &recipe.Instructions)
	if err != nil {
		return Recipe{}, err
	}

	ingredientsQuery := `SELECT id, name, quantity, unit FROM ingredients WHERE recipe_id = $1`

	// Get all ingredients used in this recipe
	rows, err := database.DB.Query(ctx, ingredientsQuery, recipe.ID)
	if err != nil {
		return Recipe{}, err
	}
	defer rows.Close()

	// Map database response onto ingredients slice
	var ingredients []Ingredient
	for rows.Next() {
		var ingredient Ingredient

		err = rows.Scan(&ingredient.ID, &ingredient.Name, &ingredient.Quantity, &ingredient.Unit)
		if err != nil {
			return Recipe{}, err
		}

		ingredients = append(ingredients, ingredient)
	}

	if err = rows.Err(); err != nil {
		return Recipe{}, err
	}

	tags, err := FindTagsByRecipe(ctx, recipe.ID)
	if err != nil {
		return Recipe{}, err
	}

	var tagStrs []string
	for _, tag := range tags {
		tagStrs = append(tagStrs, tag.Name)
	}

	// Add ingredients and tags to recipe object
	recipe.Ingredients = ingredients
	recipe.Tags = tagStrs

	return recipe, nil
}

func CreateRecipe(ctx context.Context, recipe Recipe) (int64, error) {
	userId := utils.ExtractUserIDFromContext(ctx)
	query := `INSERT INTO recipes (name, user_id, cooking_time, description, instructions) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	// Send query
	row := database.DB.QueryRow(ctx, query, recipe.Name, userId, recipe.CookingTime, recipe.Description, recipe.Instructions)

	// Get id of created recipe
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}

	// Create ingredients
	for i := 0; i < len(recipe.Ingredients); i++ {
		ingredient := recipe.Ingredients[i]
		ingredient.RecipeID = id

		_, err = CreateIngredient(ctx, ingredient)
		if err != nil {
			return -1, err
		}
	}

	// Create tags
	for _, tag := range recipe.Tags {
		_, err := CreateTag(ctx, id, tag)
		if err != nil {
			return -1, err
		}
	}

	return id, nil
}

func UpdateRecipe(ctx context.Context, id int64, recipe Recipe) (int64, error) {
	query := `UPDATE recipes SET name = $1, cooking_time = $2, description = $3, instructions = $4 WHERE id = $5 RETURNING id`

	// Send query
	_, err := database.DB.Exec(ctx, query, recipe.Name, recipe.CookingTime, recipe.Description, recipe.Instructions, id)
	if err != nil {
		return -1, nil
	}

	err = updateRecipeIngredients(ctx, id, recipe)
	if err != nil {
		return -1, nil
	}

	err = updateRecipeTags(ctx, id, recipe)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func DeleteRecipe(ctx context.Context, id int64) error {
	query := `DELETE FROM recipes WHERE id = $1`

	_, err := database.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

// Helper Functions
func updateRecipeIngredients(ctx context.Context, recipeId int64, recipe Recipe) error {
	ingredientsToDeleteSlice, err := ListIngredientsByRecipe(ctx, recipeId)
	if err != nil {
		return err
	}

	ingredientsToDelete := map[int64]bool{}
	for _, recipe := range ingredientsToDeleteSlice {
		ingredientsToDelete[recipe.ID] = true
	}

	for i := 0; i < len(recipe.Ingredients); i++ {
		ingredient := recipe.Ingredients[i]
		ingredient.RecipeID = recipeId

		// Check if ingredient was previously in recipe
		prevIngredient, err := FindIngredient(ctx, ingredient.Name, recipeId)
		if err != nil && err == pgx.ErrNoRows {
			// A previous version of this recipe's ingredient does not exist, so
			// create it
			_, err = CreateIngredient(ctx, ingredient)
			if err != nil {
				return err
			}
		} else if err != nil {
			fmt.Println("HERE")
			return err
		} else {
			// Update previous version of ingredient
			updateQuery := `UPDATE ingredients SET name = $1, quantity = $2, unit = $3 WHERE id = $4`

			_, err = database.DB.Exec(ctx, updateQuery, ingredient.Name, ingredient.Quantity, ingredient.Unit, prevIngredient.ID)
			if err != nil {
				return err
			}

			ingredientsToDelete[prevIngredient.ID] = false
		}
	}

	// Remove ingredients that are no longer used
	deleteQuery := `DELETE FROM ingredients WHERE id = $1`

	for id, delete := range ingredientsToDelete {
		if delete {
			_, err = database.DB.Exec(ctx, deleteQuery, id)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Takes a recipe ID and an updated recipe object. Updates the recipe in the
// datase to reflect the new list of tags.
func updateRecipeTags(ctx context.Context, recipeId int64, recipe Recipe) error {
	tagsToDeleteSlice, err := FindTagsByRecipe(ctx, recipeId)
	if err != nil {
		return err
	}

	tagsToDelete := map[int64]bool{}
	for _, tag := range tagsToDeleteSlice {
		tagsToDelete[tag.ID] = true
	}

	for _, tag := range recipe.Tags {
		// Check if recipe previously included tag
		prevTag, err := FindTag(ctx, recipeId, tag)
		if err != nil && err == pgx.ErrNoRows {
			// A previous version of this tag does not exist, so create it
			_, err = CreateTag(ctx, recipeId, tag)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			tagsToDelete[prevTag.ID] = false
		}
	}

	deleteQuery := `DELETE FROM recipe_tags WHERE id = $1`

	for id, delete := range tagsToDelete {
		if delete {
			_, err = database.DB.Exec(ctx, deleteQuery, id)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
