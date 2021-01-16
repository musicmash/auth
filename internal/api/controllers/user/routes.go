package user

import "github.com/go-chi/chi"

func (c *Controller) GetRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/user/info", c.Get)
	r.Post("/user/logout", c.Delete)

	return r
}
