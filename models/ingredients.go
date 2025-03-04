package models

type Ingredient struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	RecipeID int64   `json:"recipeId"`
	Quantity float32 `json:"quantity"`
	Unit     string  `json:"unit"`
}
