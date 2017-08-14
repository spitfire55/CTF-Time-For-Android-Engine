package main

import (
	"log"
	"gopkg.in/zabawaba99/firego.v1"
	"io/ioutil"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"golang.org/x/net/context"
)

func test(ctx context.Context) map[string]interface{} {
	authClient := authenticate()
	fb := firego.New("https://ctf-time-for-android.firebaseio.com/", authClient.Client(ctx))

	var v map[string]interface{}
	if err := fb.Value(&v); err != nil {
		log.Fatal(err)
	}
	return v
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