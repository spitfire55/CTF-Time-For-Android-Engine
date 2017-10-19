package main

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/zabawaba99/firego"
)

type FirebaseContext struct {
	w  http.ResponseWriter
	r  *http.Request
	c  *http.Client
	fb *firego.Firebase
}

func fetch(url string, fbc FirebaseContext) []byte {
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		http.Error(fbc.w, err.Error(), http.StatusInternalServerError)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(fbc.w, err.Error(), http.StatusInternalServerError)
	}
	defer resp.Body.Close()
	return body
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fbc := FirebaseContext{
		w, r, &http.Client{}, connect(r.Context()),
	}
	body := fetch("https://ctftime.org/api/v1/top/", fbc)
	ranking := getAllRankings(body)
	saveAllRankings(ranking)
}

func checkCurrentRankingsHandler(w http.ResponseWriter, r *http.Request) {
	fbc := FirebaseContext {
		w, r, &http.Client{}, connect(r.Context()),
	}
	body := fetch("https://ctftime.org/api/v1/top/2017/", fbc)
	ranking := getCurrentRankings(body)
	saveCurrentRankings(ranking)
}

func updateAllTeamsHandler(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	connect(ctx)
	updateAllTeams(ctx, w, r)
}

func convertTeamHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	connect(ctx)
	convertTeams()
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/current-rankings", checkCurrentRankingsHandler)
	//http.HandleFunc("/all-teams", allTeamsHandler)
	http.HandleFunc("/update-new-team", updateAllTeamsHandler)
	http.HandleFunc("/convert-team", convertTeamHandler)
	http.ListenAndServe(":80", nil)
}
