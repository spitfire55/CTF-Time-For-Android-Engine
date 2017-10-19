package main

import (
	"os"
	"strconv"

	"github.com/zabawaba99/firego"
	"golang.org/x/oauth2/google"
	"encoding/base64"
	"context"
	"log"
	"io/ioutil"
)

func connect(ctx context.Context) *firego.Firebase {

	token, err := ioutil.ReadFile("ctf-time-token.json")
	if err != nil {
		log.Fatal(err)
	}
	conf, err := google.JWTConfigFromJSON(token,
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/firebase.database")
	if err != nil {
		log.Fatal(err)
	}
	return firego.New(os.Getenv("FIREBASE_BASE"), conf.Client(ctx))
}

func saveAllRankings(teamRankings interface{}) {
	if err := fb.Child("Rankings").Set(teamRankings); err != nil {
		log.Fatal(err)
	}
}

func saveCurrentRankings(teamRankings interface{}) {
	if err := fb.Child("Rankings/2017").Set(teamRankings); err != nil {
		log.Fatal(err)
	}
}

func saveAllTeams(teams interface{}) {
	if err:= fb.Child("Teams").Set(teams); err != nil {
		log.Fatal(err)
	}
}

func setHighestNode(node int) {
	if err := fb.Child("TeamHighestNode").Set(node); err != nil {
		log.Fatal(err)
	}
}

func getHighestNode() int {
	var highestNode int
	fb.Child("TeamHighestNode").Value(&highestNode)
	return highestNode
}

func saveNewTeam(team interface{}) {
	highestNode := getHighestNode()
	// nil value passed in for team if we have reached highest team ID
	if team != nil {
		err := fb.Child("Teams/" + strconv.Itoa(highestNode)).Set(team)
		if err != nil {
			log.Fatal(err)
		}
		err = fb.Child("TeamHighestNode").Set( highestNode+ 1)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func convertTeams() {
	var teams []KeyedTeam
	if err := fb.Child("Teams").Value(&teams); err != nil {
		log.Fatal(err)
	}
	for i, v := range teams {
		if v.Name != "" {
			intId := i
			fb.Child("TeamsByName/" + base64.URLEncoding.EncodeToString([]byte(v.Name))).Set(intId)
			if v.Country != "" {
				fb.Child("TeamsByCountry/" + v.Country + "/" +
					strconv.Itoa(intId)).Set(true)
			}
			if v.Aliases != nil {
				for _, name := range v.Aliases {
					fb.Child("TeamsByName/" + base64.URLEncoding.EncodeToString([]byte(name))).Set(intId)
				}
			}
			if v.Ratings != nil {
				for year, rating := range v.Ratings {
					simpleRating := SimpleRating{
						rating.RatingPoints,
						intId,
					}
					fb.Child("TeamsByPlace/" + year + "/" +
						strconv.Itoa(rating.RatingPlace)).Set(simpleRating)
				}
			}
		}
	}
}
