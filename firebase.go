package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

func connect(token *option.ClientOption) (*firestore.Client, context.Context) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("FIREBASE_ID"), *token)
	if err != nil {
		fmt.Println(err.Error())
		return nil, nil
	} else {
		return client, ctx
	}
}

func generateToken() option.ClientOption {
	return option.WithCredentialsFile(os.Getenv("CTF_TIME_KEY"))
}

func saveCurrentRankings(teamRankings interface{}, fbc *FirebaseContext) {
	currentRankings, ok := teamRankings.(KeyedRankingsYear)
	if ok {
		_, err := fbc.fb.Collection("Rankings").Doc("2017").Set(fbc.ctx, currentRankings)
		if err != nil {
			http.Error(fbc.w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func saveAllRankings(teamRankings KeyedRankingsAll, fbc *FirebaseContext) {
	rankingYears := fbc.fb.Collection("Rankings")
	for id, ranking := range teamRankings {
		_, err := rankingYears.Doc(id).Set(fbc.ctx, ranking)
		if err != nil {
			http.Error(fbc.w, err.Error(), http.StatusInternalServerError)
		}
	}
}


/*
func setHighestNode(node int, fbc *FirebaseContext) {
	_, err := fbc.fb.Doc("TeamHighestNode").Set(fbc.r.Context(), node)
	if err != nil {
		http.Error(fbc.w, err.Error(), http.StatusInternalServerError)
	}
}
*/

func getHighestNode(fbc *FirebaseContext) int64 {
	node, err := fbc.fb.Collection("Teams").Doc("HighestID").Get(fbc.ctx)
	if err != nil {
		http.Error(fbc.w, err.Error(), http.StatusInternalServerError)
		return 0
	}

	var idInt int64
	idRaw, err := node.DataAt("id")
	if err != nil {
		http.Error(fbc.w, err.Error(), http.StatusInternalServerError)
		return 0
	}
	idInt, ok := idRaw.(int64)
	if !ok {
		http.Error(fbc.w, err.Error(), http.StatusInternalServerError)
		return 0
	}
	return idInt
}


func saveNewTeam(node int, team KeyedTeam, fbc *FirebaseContext) {
	// nil value passed in for team if we have reached highest team ID
	fbc.fb.Collection("Teams").Doc(strconv.Itoa(node)).Set(fbc.ctx, team)
}
