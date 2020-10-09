package backend

import (
	"context"
	"fmt"

	"github.com/musicmash/auth/internal/db"
	"github.com/musicmash/auth/internal/log"
	"github.com/musicmash/auth/internal/secure"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

type Backend struct {
	config *oauth2.Config
	mgr    *db.Mgr
}

func New(mgr *db.Mgr, redirectURL, appID, appSecret string, scopes ...string) *Backend {
	cfg := oauth2.Config{
		ClientID:     appID,
		ClientSecret: appSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  spotify.AuthURL,
			TokenURL: spotify.TokenURL,
		},
	}
	return &Backend{config: &cfg, mgr: mgr}
}

func (b *Backend) newSpotifyClient(token *oauth2.Token) spotify.Client {
	return spotify.NewClient(b.config.Client(context.Background(), token))
}

func getUserPhoto(images []spotify.Image) *string {
	if len(images) == 0 {
		return nil
	}

	return &images[0].URL
}

func (b *Backend) GetSession(code string) (string, error) {
	// retrieve access token
	token, err := b.config.Exchange(context.Background(), code)
	if err != nil {
		return "", fmt.Errorf("can't get access_token: %w", err)
	}

	// get user info
	client := b.newSpotifyClient(token)
	user, err := client.CurrentUser()
	if err != nil {
		return "", fmt.Errorf("couldn't get user info: %v", err)
	}

	log.Infof("user successfully logged in: %s", user.ID)

	// ensure that user exists in the db
	err = b.mgr.EnsureUserExists(&db.User{Name: user.ID, Photo: getUserPhoto(user.Images)})
	if err != nil {
		return "", fmt.Errorf("can't ensure that user exists: %w", err)
	}

	// generate sha256 string
	sid := secure.GenerateHash(token.AccessToken)

	// save session into the db
	err = b.mgr.CreateSession(&db.Session{UserName: user.ID, Value: sid})
	if err != nil {
		return "", fmt.Errorf("can't create new user session: %w", err)
	}

	return sid, nil
}

func (b *Backend) GetAuthURL(state string) string {
	return b.config.AuthCodeURL(state, oauth2.SetAuthURLParam("show_dialog", "true"))
}
