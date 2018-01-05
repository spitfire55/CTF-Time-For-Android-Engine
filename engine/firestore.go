package engine

import (
	"context"
	"errors"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

type FirebaseContext struct {
	Ctx context.Context  // context used in connection to Firestore
	Fb  firestore.Client // client used in connection to Firestore
}

func Connect(token option.ClientOption, r *http.Request) (*firestore.Client, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "ctf-time-for-android", token)
	if err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

func GenerateToken() (option.ClientOption, error) {
	if apiKey, ok := os.LookupEnv("CTF_TIME_KEY"); ok {
		return option.WithCredentialsFile(apiKey), nil
	} else {
		return nil, errors.New("API Key not found")
	}
}
