package auth

import (
	"net/http"

	"github.com/musicmash/auth/internal/db"
	"github.com/musicmash/auth/internal/log"
)


type Controller struct {
	mgr *db.Conn
}

func New(mgr *db.Conn) *Controller {
	return &Controller{mgr: mgr}
}

func (c *Controller) DoAuth(w http.ResponseWriter, r *http.Request) {
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

	session, err := c.mgr.GetSession(r.Context(), cookie.Value)
	if err != nil {
		log.Info("someone provided empty sid cookie")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("x-user-name", session.UserName)
	w.WriteHeader(http.StatusOK)
}
