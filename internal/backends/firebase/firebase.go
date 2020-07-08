package firebase

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/musicmash/auth/internal/backends"
	"google.golang.org/api/option"
)

type Backend struct {
	client *auth.Client
}

func New(serviceAccountFilePath string) (backends.Backend, error) {
	opt := option.WithCredentialsFile(serviceAccountFilePath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting Auth client: %w", err)
	}

	backend := Backend{client: client}
	return &backend, nil
}

func (b *Backend) GetUserID(idToken string) (string, error) {
	token, err := b.client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return "", fmt.Errorf("error verifying ID token: %w", err)
	}

	return token.UID, nil
}
