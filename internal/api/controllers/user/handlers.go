package user

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/musicmash/auth/internal/db"
	"github.com/musicmash/auth/internal/log"
)

type Controller struct {
	mgr *db.Conn
}

func New(mgr *db.Conn) *Controller {
	return &Controller{mgr: mgr}
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sid")
	if err != nil {
		log.Info("someone forget to provided sid cookie")
		// return 400 instead of 401 to avoid network error on the browser side
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(cookie.Value) == 0 {
		log.Info("someone provided empty sid cookie")
		// return 400 instead of 401 to avoid network error on the browser side
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = c.mgr.GetSession(r.Context(), cookie.Value)
	if errors.Is(err, sql.ErrNoRows) {
		// return 400 instead of 401 to avoid network error on the browser side
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sid")
	if err != nil {
		log.Info("someone forget to provided sid cookie")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if len(cookie.Value) == 0 {
		log.Info("someone provided empty sid cookie")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// TODO (m.kalinin): extract that code to services/backend
	session, err := c.mgr.GetSession(r.Context(), cookie.Value)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = c.mgr.DeleteSession(r.Context(), session.Value)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, getDeleteSidCookie())
	w.WriteHeader(http.StatusNoContent)
}

func getDeleteSidCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "sid",
		Value:    "",
		Path:     "/v1",
		Expires:  time.Unix(0, 0),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}
