package main

import (
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func init() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/current-rankings", checkCurrentRankingsHandler)
	//http.HandleFunc("/all-teams", allTeamsHandler)
	http.HandleFunc("/check-new-team", checkNewTeamHandler)
	http.HandleFunc("/convert-team", convertTeamHandler)
}

func setup(url string, ctx context.Context, w http.ResponseWriter) []byte {
	client := &http.Client {
		Transport: &urlfetch.Transport{
			Context: ctx,
			AllowInvalidServerCertificate: appengine.IsDevAppServer(),
		},
	}
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
	ctx := appengine.NewContext(r)
	body := setup("https://ctftime.org/api/v1/top/", ctx, w)
	ranking := getAllRankings(body)
	connect(ctx)
	saveAllRankings(ranking, ctx)
}

func checkCurrentRankingsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	body := setup("https://ctftime.org/api/v1/top/2017/", ctx, w)
	ranking := getCurrentRankings(body)
	connect(ctx)
	saveCurrentRankings(ranking, ctx)
}

/*
func allTeamsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	body := setup("https://ctftime.org/api/v1/teams/", ctx, w)
	ranking := getAllTeams(body)
	connect(ctx)
	saveAllTeams(ranking, ctx)
}
*/

func checkNewTeamHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	connect(ctx)
	recursiveTeamCheck(ctx, w)
}

func convertTeamHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	connect(ctx)
	convertTeams(ctx)
}

func main() {
	appengine.Main()
}
