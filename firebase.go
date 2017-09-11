package main

import (
	"os"
	"strconv"

	"github.com/zabawaba99/firego"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine/log"
	"encoding/base64"
)

var fb *firego.Firebase

func connect(ctx context.Context) {

	hc, err := google.DefaultClient(ctx,
		"https://www.googleapis.com/auth/firebase.database",
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

func saveAllTeams(teams interface{}, ctx context.Context) {
	if err:= fb.Child("Teams").Set(teams); err != nil {
		log.Errorf(ctx, err.Error())
	}
}

func setHighestNode(ctx context.Context, node int) {
	if err := fb.Child("TeamHighestNode").Set(node); err != nil {
		log.Errorf(ctx, err.Error())
	}
}

func getHighestNode(ctx context.Context) int {
	var highestNode int
	fb.Child("TeamHighestNode").Value(&highestNode)
	return highestNode
}

func saveNewTeam(team interface{}, ctx context.Context) {
	highestNode := getHighestNode(ctx)
	// nil value passed in for team if we have reached highest team ID
	if team != nil {
		err := fb.Child("Teams/" + strconv.Itoa(highestNode)).Set(team)
		if err != nil {
			log.Errorf(ctx, err.Error())
		}
		err = fb.Child("TeamHighestNode").Set( highestNode+ 1)
		if err != nil {
			log.Errorf(ctx, err.Error())
		}
	}
}

func convertTeams(ctx context.Context) {
	var teams []KeyedTeam
	if err := fb.Child("Teams").Value(&teams); err != nil {
		log.Errorf(ctx, err.Error())
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
