package main

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"fmt"
)

type FirebaseContext struct {
	w   http.ResponseWriter
	r   http.Request
	c   http.Client      // client used to GET from ctftime.org
	ctx context.Context  // context used in connection to Firestore
	fb  firestore.Client // client used in connection to Firestore
}

func fetch(url string, fbc FirebaseContext) *http.Response {
	resp, err := fbc.c.Get(url)
	if err != nil {
		http.Error(fbc.w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	return resp
}

/*func checkAllRankingsHandler(w http.ResponseWriter, r *http.Request) {
	token := generateToken()
	FbClient, ctx := connect(token)
	if FbClient != nil && ctx != nil {
		// create pointer to FirebaseContext
		fbc := &FirebaseContext{
			w, *r, http.Client{}, ctx, *FbClient,
		}
		body := fetch("https://ctftime.org/api/v1/top/", fbc)
		ranking := getAllRankings(body)
		saveAllRankings(ranking, fbc)
	} else {
		http.Error(w,
			"Failed to connect to Firestore",
			http.StatusInternalServerError)
	}

}
*/

func checkCurrentRankingsHandler(w http.ResponseWriter, r *http.Request) {
	token := generateToken()
	FbClient, ctx := connect(token)
	if FbClient != nil && ctx != nil {
		fbc := FirebaseContext{
			w, *r, http.Client{}, ctx, *FbClient,
		}
		//highestNode := int(getHighestNode(fbc))
		for i := 1; i < 2; i++ {
			url := fmt.Sprintf("https://ctftime.org/stats/2017?page=%d", i)
			response := fetch(url, fbc)
			if response != nil {
				getCurrentRankings(response, fbc)
				//saveCurrentRankings(ranking, fbc)
			}
		}
	} else {
		http.Error(w,
			"Failed to connect to Firestore",
			http.StatusInternalServerError)
	}
}

/*
func updateAllTeamsHandler(w http.ResponseWriter, r *http.Request) {
	token := generateToken()
	FbClient, ctx := connect(token)
	if FbClient != nil && ctx != nil {
		fbc := &FirebaseContext{
			w, *r, http.Client{}, ctx, *FbClient,
		}
		updateAllTeams(fbc, token)
	} else {
		http.Error(w,
			"Failed to connect to Firestore",
			http.StatusInternalServerError)
	}
}
*/

func defaultHandler(w http.ResponseWriter, r *http.Request) {}

func main() {
	//http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/favicon.ico", defaultHandler)
	//http.HandleFunc("/rankings", checkAllRankingsHandler)
	http.HandleFunc("/", checkCurrentRankingsHandler)
	//http.HandleFunc("/all-teams", allTeamsHandler)
	//http.HandleFunc("/update-teams", updateAllTeamsHandler)
	http.ListenAndServe("localhost:8080", nil)
}
