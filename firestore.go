// Copyright 2017-2018 Dale Lakes <spitfire@spitfy.re>. All rights reserved.
// Use of this source code is governed by the MIT license located in the LICENSE file.

package goctftime

import (
	"context"
	"errors"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

// FirebaseContext contains the necessary variables to get/set Firestore data.
type FirebaseContext struct {
	Ctx context.Context  // context used in connection to Firestore
	Fb  firestore.Client // client used in connection to Firestore
}

// NewFirebaseContext creates a new FirebaseContext object for a Firestore request.
func NewFirebaseContext(ctx context.Context, token option.ClientOption) (FirebaseContext, error) {
	teamFbClient, err := Connect(token)
	if err != nil {
		return FirebaseContext{}, err
	}
	return FirebaseContext{ctx, *teamFbClient}, nil
}

// Connect connects to Firestore and returns an authenticated client that can read to/write from the database.
func Connect(token option.ClientOption) (*firestore.Client, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "ctf-time-for-android", token)
	if err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

// GenerateToken reads an API key and returns an option to be used by Connect.
func GenerateToken() (option.ClientOption, error) {
	if apiKey, ok := os.LookupEnv("CTF_TIME_KEY"); ok {
		return option.WithCredentialsFile(apiKey), nil
	} else {
		return nil, errors.New("API Key not found")
	}
}
