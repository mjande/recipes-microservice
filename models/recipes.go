package models

type Recipe struct {
	ID           int64        `json:"id"`
	Name         string       `json:"name"`
	CookingTime  string       `json:"cookingTime"`
	Description  string       `json:"description"`
	Instructions string       `json:"instructions"`
	Ingredients  []Ingredient `json:"ingredients"`
	Tags         []string     `json:"tags"`
}
