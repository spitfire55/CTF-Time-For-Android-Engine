package main

import (
	"golang.org/x/net/context"

	"io/ioutil"
	"net/http"

	"strconv"
	"encoding/base64"
)

func updateAllTeams(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	highestNode := getHighestNode()
	for x := 0; x < highestNode; x++ {
		nodeStr:= strconv.Itoa(x)
		go updateSingleTeam(nodeStr, ctx, w, r)
	}
}

func updateSingleTeam(node string, ctx context.Context, w http.ResponseWriter, r *http.Request) {
	body := checkNewTeam(node, ctx, w)
	if body != nil {
		bodyKeyed, id := getSingleTeam(body)
		saveNewTeam(bodyKeyed)
		convertNewTeam(bodyKeyed, id)
	}
}

func checkNewTeam(node string, ctx context.Context, w http.ResponseWriter) []byte {
	client := &http.Client{}
	baseUrl := "https://ctftime.org/api/v1/teams/"
	resp, err := client.Get(baseUrl + node + "/")
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
