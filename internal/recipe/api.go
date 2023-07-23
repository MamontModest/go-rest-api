package recipe

import (
	"github.com/go-ozzo/ozzo-routing/v2"
	"github.com/mamontmodest/go-rest-api/pkg/log"
	"net/http"
	"strconv"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, logger *log.Logger) {
	res := resource{service, logger}

	r.Get("/recipe/<id>", res.get)
	r.Get("/recipe", res.getAll)

	// the following endpoints require a valid BasicAuth
	r.Post("/recipe", res.create)
	r.Put("/recipe/<id>", res.update)
	r.Delete("/recipe/<id>", res.delete)
}

type resource struct {
	service Service
	logger  *log.Logger
}

func (r *resource) get(c *routing.Context) error {
	recipeId, err := strconv.Atoi(c.Param("id"))
	album, err := r.service.Get(c.Request.Context(), recipeId)
	if err != nil {
		return c.WriteWithStatus(BadRequestResponse{}, http.StatusBadRequest)
	}

	return c.Write(album)
}

func (r *resource) getAll(c *routing.Context) error {
	ctx := c.Request.Context()
	recipesList, err := r.service.GetRecipes(ctx)
	if err != nil {
		return c.WriteWithStatus(BadRequestResponse{}, http.StatusBadRequest)
	}
	return c.Write(recipesList.Recipes)
}

func (r *resource) create(c *routing.Context) error {
	var input CreateRecipeRequest
	if err := c.Read(&input); err != nil {
		return c.WriteWithStatus(BadRequestResponse{}, http.StatusBadRequest)
	}
	album, err := r.service.Create(c.Request.Context(), input)
	if err != nil {
		return c.WriteWithStatus(BadRequestResponse{}, http.StatusBadRequest)
	}

	return c.WriteWithStatus(album, http.StatusCreated)
}

func (r *resource) update(c *routing.Context) error {
	var input UpdateRecipeRequest
	if err := c.Read(&input); err != nil {
		return c.WriteWithStatus(BadRequestResponse{}, http.StatusBadRequest)
	}
	recipeId, err := strconv.Atoi(c.Param("id"))
	noContent, err := r.service.Update(c.Request.Context(), recipeId, input)
	if err != nil {
		return c.WriteWithStatus(BadRequestResponse{}, http.StatusBadRequest)
	}

	return c.WriteWithStatus(noContent, http.StatusNoContent)
}

func (r *resource) delete(c *routing.Context) error {
	recipeId, err := strconv.Atoi(c.Param("id"))
	recipe, err := r.service.Delete(c.Request.Context(), recipeId)
	if err != nil {
		return c.WriteWithStatus(BadRequestResponse{}, http.StatusBadRequest)
	}

	return c.WriteWithStatus(recipe, http.StatusNoContent)
}
