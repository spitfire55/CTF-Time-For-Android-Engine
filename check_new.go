package main

import (
	"strconv"
	"net/http"
)

func updateAllTeams(fbc *FirebaseContext, w http.ResponseWriter, r *http.Request) {
	highestNode, ok := getHighestNode(fbc).(int)
	if ok {
		for x := 0; x < highestNode; x++ {
			go updateSingleTeam(x, &w, r)
		}
	}
}

func updateSingleTeam(node int, w *http.ResponseWriter, r *http.Request) {
	FbClient, ctx := connect()
	if FbClient != nil && ctx != nil {
		fbc := &FirebaseContext{
			*w, *r, http.Client{}, ctx, *FbClient,
		}
		body := checkNewTeam(node, fbc)
		if body != nil {
			bodyKeyed := getSingleTeam(body)
			saveNewTeam(node, bodyKeyed, fbc)
		}
	} else {
		http.Error(*w,
			"Failed to connect to Firestore",
			http.StatusInternalServerError)
	}
}

func checkNewTeam(node int, fbc *FirebaseContext) []byte {
	baseUrl := "https://ctftime.org/api/v1/teams/"
	return fetch(baseUrl+strconv.Itoa(node)+"/", fbc)
}
