package main

import (
	"io/ioutil"
	"log"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"gopkg.in/zabawaba99/firego.v1"
)

type Firebase struct {
	fbInstance *firego.Firebase
	context    context.Context
	jwtConfig  *jwt.Config
}
var fb Firebase

func connect(ctx context.Context) {
	fb.jwtConfig = authenticate()
	fb.context = ctx
	fb.fbInstance = firego.New("https://ctf-time-for-android.firebaseio.com/", fb.jwtConfig.Client(fb.context))
}

func saveAllRankings(teamRankings interface{}) {
	if err := fb.fbInstance.Child("Rankings").Set(teamRankings); err != nil {
		log.Fatal(err)
	}
}

func saveCurrentRankings(teamRankings interface{}) {
	if err := fb.fbInstance.Child("Rankings").Child("2017").Set(teamRankings); err != nil {
		log.Fatal(err)
	}
}

func authenticate() *jwt.Config {
	d, err := ioutil.ReadFile("/home/spitfire/External/CTF-Time-For-Android/firebase-adminsdk-token.json")
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
