package models

import (
	"context"

	"github.com/jackc/pgx"
	"github.com/mjande/recipes-microservice/database"
)

type Tag struct {
	ID       int64
	RecipeId int64
	Name     string
}

// Returns all tags for a given recipe.
func FindTagsByRecipe(ctx context.Context, recipeId int64) ([]Tag, error) {
	query := `SELECT id, recipe_id, name FROM recipe_tags WHERE recipe_id = $1`

	rows, err := database.DB.Query(ctx, query, recipeId)
	if err != nil && err == pgx.ErrNoRows {
		return []Tag{}, nil
	} else if err != nil {
		return []Tag{}, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag

		err = rows.Scan(&tag.ID, &tag.RecipeId, &tag.Name)
		if err != nil {
			return []Tag{}, err
		}

		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return []Tag{}, err
	}

	return tags, nil
}

// Returns a tag by recipe ID and tag name.
func FindTag(ctx context.Context, recipeId int64, name string) (Tag, error) {
	query := `SELECT id, recipe_id, name FROM recipe_tags WHERE recipe_id = $1 AND name = $2`

	result := database.DB.QueryRow(ctx, query, recipeId, name)

	var tag Tag
	err := result.Scan(&tag.ID, &tag.RecipeId, &tag.Name)
	if err != nil {
		return Tag{}, err
	}

	return tag, nil
}

// Create new recipe tag.
func CreateTag(ctx context.Context, recipeId int64, tag string) (int64, error) {
	query := `INSERT INTO recipe_tags (recipe_id, name) VALUES ($1, $2) RETURNING id`

	row := database.DB.QueryRow(ctx, query, recipeId, tag)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}
