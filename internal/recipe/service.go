package recipe

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mamontmodest/go-rest-api/internal/entity"
	"github.com/mamontmodest/go-rest-api/pkg/log"
)

// Service encapsulates useCase logic for recipes.
type Service interface {
	Get(ctx context.Context, id int) (Recipe, error)
	GetRecipes(ctx context.Context) (RecipesList, error)
	Create(ctx context.Context, input CreateRecipeRequest) (Recipe, error)
	Update(ctx context.Context, id int, input UpdateRecipeRequest) (NoContentRequestResponse, error)
	Delete(ctx context.Context, id int) (NoContentRequestResponse, error)
	GetRecipesFilter(ctx context.Context, filter filter) (RecipesList, error)
	Login(ctx context.Context, login, password string) error
}

// Recipe represents the data about a recipe.
type Recipe struct {
	entity.Recipe
}

// RecipesList represents the data about recipes.
type RecipesList struct {
	entity.RecipesList
}

// CreateRecipeRequest represents an recipe creation request.
type CreateRecipeRequest struct {
	Name        string               `json:"name"`
	Title       string               `json:"title"`
	Ingredients []*entity.Ingredient `json:"ingredients"`
	Steps       []*entity.Step       `json:"steps"`
}

// Validate validates the CreateRecipeRequest fields.
func (m CreateRecipeRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Title, validation.Required, validation.Length(0, 255)),
		validation.Field(&m.Ingredients, validation.Required, validation.Length(0, 31)),
		validation.Field(&m.Steps, validation.Required, validation.Length(0, 31)),
	)
}

type UpdateRecipeRequest struct {
	Name        string               `json:"name"`
	Title       string               `json:"title"`
	Ingredients []*entity.Ingredient `json:"ingredients"`
	Steps       []*entity.Step       `json:"steps"`
}

type BadRequestResponse struct {
}

type CreateRequestResponse struct {
}

type NoContentRequestResponse struct {
}

// Validate validates the UpdateRecipeRequest fields.
func (m UpdateRecipeRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Title, validation.Required, validation.Length(0, 255)),
		validation.Field(&m.Ingredients, validation.Required, validation.Length(0, 31)),
		validation.Field(&m.Steps, validation.Required, validation.Length(0, 31)),
	)
}

type service struct {
	repo   Repository
	logger *log.Logger
}

// NewService creates a new recipe service.
func NewService(repo Repository, logger *log.Logger) Service {
	return service{repo, logger}
}

// Get returns the recipe with the specified the recipe ID.
func (s service) Get(ctx context.Context, id int) (Recipe, error) {
	recipe, err := s.repo.Get(ctx, id)
	if err != nil {
		return Recipe{}, err
	}
	return Recipe{recipe}, nil
}

// GetRecipesFilter returns the list of recipes with filter .
func (s service) GetRecipesFilter(ctx context.Context, filter filter) (RecipesList, error) {
	f := entity.Filter{
		TimeBetween: filter.TimeBetween,
		SortTime:    filter.SortTime,
		Ingredients: filter.Ingredients,
	}
	recipesList, err := s.repo.GetRecipesFilter(ctx, f)
	if err != nil {
		return RecipesList{recipesList}, err
	}
	return RecipesList{recipesList}, nil
}

func (s service) GetRecipes(ctx context.Context) (RecipesList, error) {
	recipesList, err := s.repo.GetRecipes(ctx)
	if err != nil {
		return RecipesList{recipesList}, err
	}
	return RecipesList{recipesList}, nil
}

// Create creates a new recipe.
func (s service) Create(ctx context.Context, req CreateRecipeRequest) (Recipe, error) {
	if err := req.Validate(); err != nil {
		return Recipe{}, err
	}
	recipe := entity.Recipe{
		Name:        req.Name,
		Title:       req.Title,
		Steps:       req.Steps,
		Ingredients: req.Ingredients,
	}
	rec, err := s.repo.Create(ctx, recipe)
	if err != nil {
		return Recipe{}, err
	}
	return Recipe{rec}, err
}

// Update updates the recipe with the specified ID.
func (s service) Update(ctx context.Context, id int, req UpdateRecipeRequest) (NoContentRequestResponse, error) {
	if err := req.Validate(); err != nil {
		return NoContentRequestResponse{}, err
	}
	recipe, err := s.repo.Get(ctx, id)
	if err != nil {
		return NoContentRequestResponse{}, err
	}
	err = s.repo.Delete(ctx, id)
	if err != nil {
		return NoContentRequestResponse{}, err
	}
	recipe.RecipeId = id
	rc := entity.Recipe{
		RecipeId:    id,
		Name:        req.Name,
		Title:       req.Title,
		Ingredients: req.Ingredients,
		Steps:       req.Steps,
	}
	if err := s.repo.CreateWithId(ctx, rc); err != nil {
		return NoContentRequestResponse{}, err
	}
	return NoContentRequestResponse{}, nil
}

// Delete deletes the recipe with the specified ID.
func (s service) Delete(ctx context.Context, id int) (NoContentRequestResponse, error) {
	_, err := s.Get(ctx, id)
	if err != nil {
		return NoContentRequestResponse{}, err
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return NoContentRequestResponse{}, err
	}
	return NoContentRequestResponse{}, nil
}

func (s service) Login(ctx context.Context, login, password string) error {
	person := entity.Identity{
		Login:    login,
		Password: password,
	}
	err := s.repo.Login(ctx, person)
	if err != nil {
		return err
	}
	return nil
}

type Ingredient struct {
	entity.Ingredient
}

type filter struct {
	entity.Filter
}
