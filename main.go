package main

import (
	"context"
	"io/ioutil"
	"net/http"

	"cloud.google.com/go/firestore"
)

type FirebaseContext struct {
	w  http.ResponseWriter
	r  http.Request
	c  http.Client      // client used to GET from ctftime.org
	ctx context.Context // context used in connection to Firestore
	fb firestore.Client // client used in connection to Firestore
}

func fetch(url string, fbc *FirebaseContext) []byte {
	resp, err := fbc.c.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(fbc.w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	defer resp.Body.Close()
	return body
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	FbClient, ctx := connect()
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

func checkCurrentRankingsHandler(w http.ResponseWriter, r *http.Request) {
	FbClient, ctx := connect()
	if FbClient != nil && ctx != nil {
		fbc := &FirebaseContext{
			w, *r, http.Client{}, ctx, *FbClient,
		}
		body := fetch("https://ctftime.org/api/v1/top/2017/", fbc)
		ranking := getCurrentRankings(body)
		saveCurrentRankings(ranking, fbc)
	} else {
		http.Error(w,
			"Failed to connect to Firestore",
			http.StatusInternalServerError)
	}
}

func updateAllTeamsHandler(w http.ResponseWriter, r *http.Request) {
	FbClient, ctx := connect()
	if FbClient != nil && ctx != nil {
		fbc := &FirebaseContext{
			w, *r, http.Client{}, ctx, *FbClient,
		}
		updateAllTeams(fbc)
	} else {
		http.Error(w,
			"Failed to connect to Firestore",
			http.StatusInternalServerError)
	}
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/current-rankings", checkCurrentRankingsHandler)
	//http.HandleFunc("/all-teams", allTeamsHandler)
	http.HandleFunc("/update-new-team", updateAllTeamsHandler)
	http.ListenAndServe("localhost:8080", nil)
}
