package spotify

import (
	"fmt"
	"net/http"

	"github.com/musicmash/auth/internal/log"
	"github.com/zmb3/spotify"
)

type Handler struct {
	state string
	auth  *spotify.Authenticator
}

func NewHandler(state string, auth *spotify.Authenticator) *Handler {
	return &Handler{state: state, auth: auth}
}

func (h *Handler) DoAuth(w http.ResponseWriter, r *http.Request) {
	token, err := h.auth.Token(h.state, r)
	if err != nil {
		http.Error(w, "couldn't get token", http.StatusUnauthorized)
		log.Errorf("couldn't get token: %v", err)
		return
	}

	client := h.auth.NewClient(token)
	user, err := client.CurrentUser()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		log.Errorf("couldn't get user info: %v", err)
		return
	}

	_, _ = fmt.Fprintf(w, "logged as %s", user.ID)
}
