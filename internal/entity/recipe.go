package entity

type Recipe struct {
	RecipeId    int           `json:"recipe_id"`
	Name        string        `json:"name"`
	Title       string        `json:"title"`
	Ingredients []*Ingredient `json:"ingredients"`
	Steps       []*Step       `json:"steps"`
}
