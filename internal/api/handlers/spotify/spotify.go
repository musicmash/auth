package spotify

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/musicmash/auth/internal/backend"
	"github.com/musicmash/auth/internal/log"
)

const (
	state = "auth"
)

type Handler struct {
	backend *backend.Backend
}

func NewHandler(b *backend.Backend) *Handler {
	return &Handler{backend: b}
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

func validateStateAndCode(values url.Values) error {
	code := values.Get("code")
	if code == "" {
		return errors.New("didn't get access code")
	}

	actualState := values.Get("state")
	if actualState != state {
		return errors.New("redirect state parameter doesn't match")
	}

	return nil
}

func (h *Handler) DoAuth(w http.ResponseWriter, r *http.Request) {
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

	sid, err := h.backend.GetSession(values.Get("code"))
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	http.SetCookie(w, newSidCookie(sid))
	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}
