package main

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"

	"io/ioutil"
	"net/http"

	"strconv"
	"encoding/base64"
)

func updateAllTeams(ctx context.Context, w http.ResponseWriter) {

}

func recursiveTeamCheck(ctx context.Context, w http.ResponseWriter) {
	body := checkNewTeam(ctx, w)
	if body != nil {
		bodyKeyed, id := getSingleTeam(body)
		saveNewTeam(bodyKeyed)
		convertNewTeam(bodyKeyed, id)
		recursiveTeamCheck(ctx, w)
	}
}

func checkNewTeam(ctx context.Context, w http.ResponseWriter) []byte {
	highestNode := getHighestNode()
	client := &http.Client{}
	highestNodeString:= strconv.Itoa(highestNode)
	baseUrl := "https://ctftime.org/api/v1/teams/"
	resp, err := client.Get(baseUrl + highestNodeString + "/")
	// if this team ID does not exist
	if err != nil || resp.StatusCode != 200 {
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return body
}

func convertNewTeam(team KeyedTeam, id int) {

	fb.Child("TeamsByName/" +
		base64.URLEncoding.EncodeToString([]byte(team.Name))).Set(id)
	if team.Country != "" {
		fb.Child("TeamsByCountry/" + team.Country + "/" +
			strconv.Itoa(id)).Set(true)
	}
	if team.Aliases != nil {
		for _, name := range team.Aliases {
			fb.Child("TeamsByName/" +
				base64.URLEncoding.EncodeToString([]byte(name))).Set(id)
		}
	}
	if team.Ratings != nil {
		for year, rating := range team.Ratings {
			simpleRating := SimpleRating{
				rating.RatingPoints,
				id,
			}
			fb.Child("TeamsByPlace/" + year + "/" +
				strconv.Itoa(rating.RatingPlace)).Set(simpleRating)
		}
	}
}
