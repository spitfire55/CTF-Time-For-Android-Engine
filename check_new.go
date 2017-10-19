package main

import (
	"strconv"
	"encoding/base64"
)

func updateAllTeams(fbc FirebaseContext) {
	//highestNode := getHighestNode(fbc.fb)
	for x := 0; x < 10; x++ {
		nodeStr:= strconv.Itoa(x)
		go updateSingleTeam(nodeStr, fbc)
	}
}

func updateSingleTeam(node string, fbc FirebaseContext) {
	body := checkNewTeam(node, fbc)
	if body != nil {
		bodyKeyed, id := getSingleTeam(body)
		saveNewTeam(bodyKeyed, fbc.fb)
		convertNewTeam(bodyKeyed, id, fbc)
	}
}

func checkNewTeam(node string, fbc FirebaseContext) []byte {
	baseUrl := "https://ctftime.org/api/v1/teams/"
	return fetch(baseUrl + node + "/", fbc)
}

func convertNewTeam(team KeyedTeam, id int, fbc FirebaseContext) {

	fbc.fb.Child("TeamsByName/" +
		base64.URLEncoding.EncodeToString([]byte(team.Name))).Set(id)
	if team.Country != "" {
		fbc.fb.Child("TeamsByCountry/" + team.Country + "/" +
			strconv.Itoa(id)).Set(true)
	}
	if team.Aliases != nil {
		for _, name := range team.Aliases {
			fbc.fb.Child("TeamsByName/" +
				base64.URLEncoding.EncodeToString([]byte(name))).Set(id)
		}
	}
	if team.Ratings != nil {
		for year, rating := range team.Ratings {
			simpleRating := SimpleRating{
				rating.RatingPoints,
				id,
			}
			fbc.fb.Child("TeamsByPlace/" + year + "/" +
				strconv.Itoa(rating.RatingPlace)).Set(simpleRating)
		}
	}
}
