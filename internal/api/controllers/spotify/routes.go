package spotify

import "github.com/go-chi/chi"

func (c *Controller) GetRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/auth", c.DoAuth)

	return r
}
