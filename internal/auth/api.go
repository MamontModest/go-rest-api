package auth

import (
	"errors"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"log"
	"net/http"
)

func RegisterHandlers(r *routing.RouteGroup, service Service) {
	res := resource{service}
	r.Post("/auth", res.create)
}
func (r *resource) create(c *routing.Context) error {
	var input Identity
	if err := c.Read(&input); err != nil {
		return c.WriteWithStatus(BadRequestResponse{}, http.StatusBadRequest)
	}
	err := r.service.CreateUser(c.Request.Context(), input.Login, input.Password)
	if err != nil {
		log.Println(err)
		return c.WriteWithStatus(BadRequestResponse{}, http.StatusBadRequest)
	}

	return nil
}
func (r *resource) Login(c *routing.Context) error {
	login, password, ok := c.Request.BasicAuth()
	if !ok {
		errors.New("unauthorized")
	}
	err := r.service.Login(c.Request.Context(), login, password)
	if err != nil {
		return err
	}
	return nil
}

type resource struct {
	service Service
}
type BadRequestResponse struct {
}
