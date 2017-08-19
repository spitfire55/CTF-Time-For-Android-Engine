package main

import (
	"io/ioutil"
	"log"

	"github.com/zabawaba99/firego"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"

	"os"
)

var fb *firego.Firebase

func connect(ctx context.Context) {
	auth := authenticate()
	fb = firego.New(os.Getenv("FIREBASE_BASE"), auth.Client(ctx))

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

func authenticate() *jwt.Config {
	d, err := ioutil.ReadFile("./creds.json")
	if err != nil {
		log.Fatal(err)
	}
	conf, err := google.JWTConfigFromJSON(d, "https://www.googleapis.com/auth/firebase",
		"https://www.googleapis.com/auth/userinfo.email")
	if err != nil {
		log.Fatal(err)
	}
	return conf

}
