package main

import (
	"os"
	"strconv"
	"net/http"
	"encoding/base64"
	"context"
	"log"
	"google.golang.org/api/option"
	"cloud.google.com/go/firestore"
)

func connect() *firestore.Client {
	ctx := context.Background()
	token := option.WithCredentialsFile("ctf-time-engine.json")
	client, err := firestore.NewClient(ctx, os.Getenv("FIREBASE_ID"), token)
	if err != nil {
		return nil
	} else {
		return client
	}
}

func saveCurrentRankings(teamRankings interface{}, fbc *FirebaseContext) {
	_, err := fbc.fb.Collection("Rankings").Doc("2017").Set(fbc.r.Context(), teamRankings)
	if err != nil {
		http.Error(fbc.w, err.Error(), http.StatusInternalServerError)
	}
}

func setHighestNode(node int, fbc *FirebaseContext) {
	_, err := fbc.fb.Doc("TeamHighestNode").Set(fbc.r.Context(), node)
	if err != nil {
		http.Error(fbc.w, err.Error(), http.StatusInternalServerError)
	}
}

func getHighestNode(fbc *FirebaseContext) int {
	node, err:= fbc.fb.Doc("TeamHighestNode").Get(fbc.r.Context())
	if err != nil {
		http.Error(fbc.w, err.Error(), http.StatusInternalServerError)
	}
	return node.Data()["node"].(int)
}

func saveNewTeam(node int, team interface{}, fbc *FirebaseContext) {
	// nil value passed in for team if we have reached highest team ID
	if team != nil {
		_, err := fbc.fb.Collection("Teams").Doc(strconv.Itoa(node)).Set(fbc.r.Context(), team)
		if err != nil {
			log.Fatal(err)
		}
	}
}

//TODO: Overhaul this to be less ugly
func convertTeams(fbc *FirebaseContext) {
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
