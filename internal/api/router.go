package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/musicmash/auth/internal/api/controllers/auth"
	"github.com/musicmash/auth/internal/api/controllers/spotify"
	"github.com/musicmash/auth/internal/api/controllers/user"
	"github.com/musicmash/auth/internal/db"
	"github.com/musicmash/auth/internal/services/backend"
)

func GetRouter(mgr *db.Conn, backend *backend.Backend) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		// user logger inside /v1 route
		// to avoid logging of healthz requests
		r.Use(middleware.Logger)

		r.Mount("/auth", auth.New(mgr).GetRouter())
		r.Mount("/callbacks/spotify", spotify.New(backend).GetRouter())
		r.Mount("/", user.New(mgr).GetRouter())
	})

	return r
}
