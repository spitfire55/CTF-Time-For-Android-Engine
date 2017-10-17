package main

import (
	"io/ioutil"
	"net/http"
	"context"
)

func init() {

}

func setup(url string, w http.ResponseWriter) []byte {
	client := &http.Client {}
	resp, err := client.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer resp.Body.Close()
	return body
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	body := setup("https://ctftime.org/api/v1/top/", w)
	ranking := getAllRankings(body)
	connect(ctx)
	saveAllRankings(ranking)
}

func checkCurrentRankingsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	body := setup("https://ctftime.org/api/v1/top/2017/", w)
	ranking := getCurrentRankings(body)
	connect(ctx)
	saveCurrentRankings(ranking)
}

func updateAllTeamsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	connect(ctx)
	updateAllTeams(ctx, w)
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
