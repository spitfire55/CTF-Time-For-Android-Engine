package main

import (
	"os"

	"github.com/zabawaba99/firego"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine/log"
)

var fb *firego.Firebase

func connect(ctx context.Context) {

	hc, err := google.DefaultClient(ctx,
		"https://www.googleapi.com/auth/firebase.database",
		"https://www.googleapis.com/auth/userinfo.email")
	if err != nil {
		log.Errorf(ctx, err.Error())
	}
	fb = firego.New(os.Getenv("FIREBASE_BASE"), hc)
}

func saveAllRankings(teamRankings interface{}, ctx context.Context) {
	if err := fb.Child("Rankings").Set(teamRankings); err != nil {
		log.Errorf(ctx, err.Error())
	}
}

func saveCurrentRankings(teamRankings interface{}, ctx context.Context) {
	if err := fb.Child("Rankings/2017").Set(teamRankings); err != nil {
		log.Errorf(ctx, err.Error())
	}
}
