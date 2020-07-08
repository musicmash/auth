package main

import (
	"log"

	"golang.org/x/net/context"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

const (
	serviceAccountFilePath = "serviceAccountKey.json"
	idToken                = "xxx"
)

func main() {
	opt := option.WithCredentialsFile(serviceAccountFilePath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	token, err := client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		log.Fatalf("error verifying ID token: %v\n", err)
	}

	log.Printf("Verified ID token: %v\n", token)
	log.Printf("user_id: %v", token.UID)
}
