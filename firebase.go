package main

import (
	"log"
	"os"

	"github.com/zabawaba99/firego"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
)

var fb *firego.Firebase

func connect(ctx context.Context) {

	hc, err := google.DefaultClient(ctx,
		"https://www.googleapi.com/auth/firebase.database",
		"https://www.googleapis.com/auth/userinfo.email")
	if err != nil {
		log.Fatal(err)
	}
	fb = firego.New(os.Getenv("FIREBASE_BASE"), hc)
}

func saveAllRankings(teamRankings interface{}) {
	if err := fb.Child("Rankings").Set(teamRankings); err != nil {
		log.Fatal(err)
	}
}

func saveCurrentRankings(teamRankings interface{}) {
	if err := fb.Child("Rankings").Child("2017").Set(teamRankings); err != nil {
		log.Fatal(err)
	}
}
