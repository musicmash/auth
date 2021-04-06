package auth

import "github.com/go-chi/chi"

func (c *Controller) GetRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/nginx", c.DoAuth)
	r.Get("/traefik", c.DoAuth)

	return r
}
