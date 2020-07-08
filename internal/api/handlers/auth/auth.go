package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/musicmash/auth/internal/backends"
)

type Handler struct {
	backend backends.Backend
}

func NewHandler(backend backends.Backend) *Handler {
	return &Handler{backend: backend}
}

func (h *Handler) DoAuth(w http.ResponseWriter, r *http.Request) {
	idToken := r.Header.Get("x-musicmash-access-token")
	if strings.Trim(idToken, " ") == "" {
		fmt.Println("someone provided empty token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	uid, err := h.backend.GetUserID(idToken)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("x-user-name", uid)
	w.WriteHeader(http.StatusOK)
}
