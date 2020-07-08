package main

import (
	"fmt"
	"log"

	"github.com/musicmash/auth/internal/backends/firebase"
)

const (
	serviceAccountFilePath = "serviceAccountKey.json"
	idToken                = "xxx"
)

func main() {
	backend, err := firebase.New(serviceAccountFilePath)
	if err != nil {
		log.Fatal(err)
	}

	uid, err := backend.GetUserID(idToken)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("user_id: %v\n", uid)
}
