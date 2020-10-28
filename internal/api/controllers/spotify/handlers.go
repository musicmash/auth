package spotify

import (
	"net/http"
	"time"

	"github.com/musicmash/auth/internal/backend"
	"github.com/musicmash/auth/internal/log"
)

type Controller struct {
	backend *backend.Backend
}

func New(backend *backend.Backend) *Controller {
	return &Controller{backend: backend}
}

func newSidCookie(sid string) *http.Cookie {
	return &http.Cookie{
		Name:     "sid",
		Value:    sid,
		Path:     "/v1",
		Expires:  time.Now().UTC().AddDate(0, 3, 0),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func (c *Controller) DoAuth(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	if err := values.Get("error"); err != "" {
		if err != "access_denied" {
			log.Errorf("got '%v' error query when try to sync artists", err)
		}

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	if err := validateStateAndCode(values); err != nil {
		log.Error(err.Error())
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if err := validateStateAndCode(values); err != nil {
		log.Error(err.Error())
		http.Error(w, "bad request", http.StatusBadRequest)
	}

	sid, err := c.backend.GetSession(r.Context(), values.Get("code"))
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	http.SetCookie(w, newSidCookie(sid))
	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}
