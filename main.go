package main

import (
	"io/ioutil"
	"net/http"

	"cloud.google.com/go/firestore"
)

type FirebaseContext struct {
	w  http.ResponseWriter
	r  http.Request
	c  http.Client // client used to GET from ctftime.org
	fb firestore.Client //client used to POST to Firestore
}

func fetch(url string, fbc *FirebaseContext) []byte {
	resp, err := fbc.c.Get(url)
	if err != nil {
		http.Error(fbc.w, err.Error(), http.StatusInternalServerError)
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
	fbc := FirebaseContext{
		w, *r, http.Client{}, *connect(),

	}
	body := fetch("https://ctftime.org/api/v1/top/", &fbc)
	ranking := getAllRankings(body)
	saveAllRankings(ranking, fbc.fb)
}

func checkCurrentRankingsHandler(w http.ResponseWriter, r *http.Request) {
	fbc := FirebaseContext {
		w, r, &http.Client{}, connect(),
	}
	body := fetch("https://ctftime.org/api/v1/top/2017/", fbc)
	ranking := getCurrentRankings(body)
	saveCurrentRankings(ranking, fbc.fb)
}

func updateAllTeamsHandler(w http.ResponseWriter, r *http.Request) {
	fbc := FirebaseContext {
		&w, r, &http.Client{}, connect(),
	}
	updateAllTeams(fbc)
}

func convertTeamHandler(w http.ResponseWriter, r *http.Request) {
	fbc := FirebaseContext{
		&w, r, &http.Client{}, connect(),
	}
	convertTeams(fbc.fb)
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/current-rankings", checkCurrentRankingsHandler)
	//http.HandleFunc("/all-teams", allTeamsHandler)
	http.HandleFunc("/update-new-team", updateAllTeamsHandler)
	http.HandleFunc("/convert-team", convertTeamHandler)
	http.ListenAndServe(":80", nil)
}
