package engine

import (
	"context"
	"errors"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

type FirebaseContext struct {
	Ctx context.Context  // context used in connection to Firestore
	Fb  firestore.Client // client used in connection to Firestore
}

func NewFirebaseContext(ctx context.Context, token option.ClientOption) (FirebaseContext, error) {
	teamFbClient, err := Connect(token)
	if err != nil {
		return FirebaseContext{}, err
	}
	return FirebaseContext{ctx, *teamFbClient}, nil
}

func Connect(token option.ClientOption) (*firestore.Client, error) {
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
