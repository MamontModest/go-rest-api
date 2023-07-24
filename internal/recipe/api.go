package recipe

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-ozzo/ozzo-routing/v2"
	"github.com/mamontmodest/go-rest-api/pkg/log"
	"net/http"
	"strconv"
	"strings"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, logger *log.Logger) {
	res := resource{service, logger}

	r.Get("/recipe/<id>", res.get)
	r.Get("/recipe", res.getAll)

	// the following endpoints require a valid BasicAuth
	r.Put("/recipe/<id>", res.update)
	r.Post("/recipe", res.create)
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
	f := new(filter)
	err := json.Unmarshal([]byte(c.Request.FormValue("filter")), &f)
	//проверяем валидность фильтра при его наличие
	if err != nil && c.Request.FormValue("filter") != "" {
		return c.WriteWithStatus(BadRequestResponse{}, http.StatusBadRequest)
	}
	//значит нет фильтров
	if f.Ingredients == nil && f.TimeBetween == "" && f.SortTime == "" {
		recipesList, err := r.service.GetRecipes(ctx)
		if err != nil {
			return c.WriteWithStatus(BadRequestResponse{}, http.StatusBadRequest)
		}
		return c.Write(recipesList.Recipes)
	}
	err = f.Validate()
	if err != nil {
		return c.WriteWithStatus(BadRequestResponse{}, http.StatusBadRequest)
	}
	recipesList, err := r.service.GetRecipesFilter(ctx, *f)
	if err != nil {
		return c.WriteWithStatus(BadRequestResponse{}, http.StatusBadRequest)
	}
	return c.Write(recipesList.Recipes)

}

func (r *resource) create(c *routing.Context) error {
	login, password, ok := c.Request.BasicAuth()
	if !ok {
		return c.WriteWithStatus(BadRequestResponse{}, http.StatusUnauthorized)
	}
	err := r.service.Login(c.Request.Context(), login, password)
	if err != nil {
		return c.WriteWithStatus(BadRequestResponse{}, http.StatusBadRequest)
	}
	var input CreateRecipeRequest
	if err := c.Read(&input); err != nil {
		return c.WriteWithStatus(BadRequestResponse{}, http.StatusBadRequest)
	}
	recipe, err := r.service.Create(c.Request.Context(), input)
	if err != nil {
		return c.WriteWithStatus(BadRequestResponse{}, http.StatusBadRequest)
	}

	return c.WriteWithStatus(recipe, http.StatusCreated)
}

func (r *resource) update(c *routing.Context) error {
	var input UpdateRecipeRequest
	login, password, ok := c.Request.BasicAuth()
	if !ok {
		return c.WriteWithStatus(BadRequestResponse{}, http.StatusUnauthorized)
	}
	err := r.service.Login(c.Request.Context(), login, password)
	if err != nil {
		return c.WriteWithStatus(BadRequestResponse{}, http.StatusBadRequest)
	}
	if err := c.Read(&input); err != nil {
		return c.WriteWithStatus(BadRequestResponse{}, http.StatusBadRequest)
	}
	recipeId, err := strconv.Atoi(c.Param("id"))
	_, err = r.service.Update(c.Request.Context(), recipeId, input)
	if err != nil {
		return c.WriteWithStatus(BadRequestResponse{}, http.StatusBadRequest)
	}
	return nil
}

func (r *resource) delete(c *routing.Context) error {
	login, password, ok := c.Request.BasicAuth()
	if !ok {
		return c.WriteWithStatus(BadRequestResponse{}, http.StatusUnauthorized)
	}
	err := r.service.Login(c.Request.Context(), login, password)
	if err != nil {
		return c.WriteWithStatus(BadRequestResponse{}, http.StatusBadRequest)
	}
	recipeId, err := strconv.Atoi(c.Param("id"))
	recipe, err := r.service.Delete(c.Request.Context(), recipeId)
	if err != nil {
		return c.WriteWithStatus(BadRequestResponse{}, http.StatusBadRequest)
	}

	return c.WriteWithStatus(recipe, http.StatusNoContent)
}

func (m filter) Validate() error {
	if m.SortTime != "asc" && m.SortTime != "desc" && m.SortTime != "" {
		return errors.New(fmt.Sprintf("Field SortTime not asc or desc : SortTime is %s", m.SortTime))
	}
	if m.TimeBetween != "" {
		mass := strings.Split(m.TimeBetween, ":")
		_, err := strconv.Atoi(mass[0])
		if err != nil {
			errors.New(fmt.Sprintf("Field TimeBeetween not correct format : TimeBeetwen is %s", m.TimeBetween))
		}
		_, err = strconv.Atoi(mass[1])
		if err != nil {
			errors.New(fmt.Sprintf("Field TimeBeetween not correct format : TimeBeetwen is %s", m.TimeBetween))
		}
	}
	return nil
}
