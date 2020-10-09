package auth

import (
	"net/http"

	"github.com/musicmash/auth/internal/db"
	"github.com/musicmash/auth/internal/log"
)

type Handler struct {
	mgr *db.Mgr
}

func NewHandler(mgr *db.Mgr) *Handler {
	return &Handler{mgr: mgr}
}

func (h *Handler) DoAuth(w http.ResponseWriter, r *http.Request) {
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

	session, err := h.mgr.GetSession(cookie.Value)
	if err != nil {
		log.Info("someone provided empty sid cookie")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("x-user-name", session.UserName)
	w.WriteHeader(http.StatusOK)
}
