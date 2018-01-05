package engine

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
)

type FirebaseContext struct {
	W   http.ResponseWriter
	R   http.Request
	C   http.Client      // client used to GET from ctftime.org
	Ctx context.Context  // context used in connection to Firestore
	Fb  firestore.Client // client used in connection to Firestore
}

func Connect(token option.ClientOption, r *http.Request) *firestore.Client {
	ctx := appengine.NewContext(r)
	client, err := firestore.NewClient(ctx, "ctf-time-for-android", token)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	} else {
		return client
	}
}

func GenerateToken() option.ClientOption {
	return option.WithCredentialsFile(os.Getenv("CTF_TIME_KEY"))
}
