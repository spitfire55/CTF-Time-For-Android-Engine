package main

import (
	"encoding/json"
	"log"
	"strconv"
)

/*
 * RANKINGS
 */
type AllRankings struct {
	Sixteen   []Rankings  `json:"2016"`
	Seventeen []Rankings  `json:"2017"`
	Eleven    []Rankings  `json:"2011"`
	Twelve    []Rankings  `json:"2012"`
	Thirteen  []Rankings  `json:"2013"`
	Fourteen  []Rankings  `json:"2014"`
	Fifteen   []Rankings  `json:"2015"`
}

type CurrentRankings struct {
	Rankings []Rankings `json:"2017"`
}

type Rankings struct {
	TeamName string  `json:"team_name,omitempty"`
	Points   float64 `json:"points,omitempty"`
	Id       int     `json:"team_id,omitempty"`
}

type ValidRankings [][]Rankings

type KeyedRankingsYear map[string]Rankings
type KeyedRankingsAll map[string]KeyedRankingsYear

/*
 * TEAMS
 */

type Team struct {
	Country  string       `json:"country"`
	Academic bool         `json:"academic"`
	Id       int          `json:"id"`
	Name     string       `json:"name"`
	Aliases  []string     `json:"aliases"`
	Ratings  []RatingYear `json:"rating"`
}

type KeyedTeam struct {
	Country  string            `json:"country,omitempty"`
	Academic bool              `json:"academic"`
	Id       int               `json:"-"`
	Name     string            `json:"name"`
	Aliases  []string          `json:"aliases,omitempty"`
	Ratings  map[string]Rating `json:"rating,omitempty"`
}

type RatingYear struct {
	Seventeen Rating `json:"2017"`
	Sixteen   Rating `json:"2016"`
	Fifteen   Rating `json:"2015"`
	Fourteen  Rating `json:"2014"`
	Thirteen  Rating `json:"2013"`
	Twelve    Rating `json:"2012"`
	Eleven	  Rating `json:"2011"`
}

type Rating struct {
	OrganizerPoints float64 `json:"organizer_points"`
	RatingPoints    float64 `json:"rating_points"`
	RatingPlace     int     `json:"rating_place"`
}

// key value will be place
/*
type SimpleRating struct {
	Points float64
	Id     int
}
*/

//type Teams []Team
//type KeyedTeams map[string]KeyedTeam

func getAllRankings(jsonStream []byte) KeyedRankingsAll {

	var results AllRankings
	err := json.Unmarshal(jsonStream, &results)
	if err != nil {
		log.Fatal(err)
	}
	// store the valid top 10 in its own slice, which is a 2-d map of rankings (rows are years, columns are Rankings structs)
	validRankings := ValidRankings{results.Eleven, results.Twelve, results.Thirteen, results.Fourteen, results.Fifteen, results.Sixteen, results.Seventeen}
	// store just the years as slice of strings.
	validRankingsYears := []string{"2011", "2012", "2013", "2014", "2015", "2016", "2017"}

	// initialize an empty map that will eventually contain contents to store in Firebase. Key = year, value = KeyedRankingsYears map
	var keyRankings = make(KeyedRankingsAll, len(validRankings))
	for i, yearRankings := range validRankings {
		// initialize inner map. Key = team id, value = KeyedRankingValue interface
		keyRankings[validRankingsYears[i]] = make(map[string]Rankings, len(validRankings[i]))
		for j, ranking := range yearRankings {
			// indices of validRankingsYears align to the order in which validRankings are stored and thus, contains the correct
			// year to use as a key value for the outer maps. 0 = 2012, 1 = 2015, etc.
			keyRankings[validRankingsYears[i]][strconv.Itoa(j)] = ranking
		}
	}
	return keyRankings
}

func getCurrentRankings(jsonStream []byte) KeyedRankingsYear {
	var results CurrentRankings
	err := json.Unmarshal(jsonStream, &results)
	if err != nil {
		log.Fatal(err)
	}

	var keyCurrentRankings = make(map[string]Rankings, len(results.Rankings))
	for i, ranking := range results.Rankings {
		keyCurrentRankings[strconv.Itoa(i)] = ranking
	}
	return keyCurrentRankings
}

/*
func getAllTeams(jsonStream []byte) KeyedTeams {
	var teams Teams
	err := json.Unmarshal(jsonStream, &teams)
	if err != nil {
		log.Fatal(err)
	}

	keyTeams := make(map[string]KeyedTeam)
	for _, team := range teams {
		key := strconv.Itoa(team.Id)
		var value KeyedTeam
		value = KeyedTeam{
			team.Country,
			team.Academic,
			team.Id,
			team.Name,
			team.Aliases,
			nil,
		}
		keyTeams[key] = value
	}
	return keyTeams
}
*/

func getSingleTeam(jsonStream []byte) KeyedTeam {
	var team Team
	err := json.Unmarshal(jsonStream, &team)
	if err != nil {
		log.Fatal(err)
	}

	finalRatings := make(map[string]Rating)
	for _, team := range team.Ratings {
		if team.Eleven.RatingPlace != 0 {
			finalRatings["2011"] = team.Eleven
		}
		if team.Twelve.RatingPlace != 0 {
			finalRatings["2012"] = team.Twelve
		}
		if team.Thirteen.RatingPlace != 0 {
			finalRatings["2013"] = team.Thirteen
		}
		if team.Fourteen.RatingPlace != 0 {
			finalRatings["2014"] = team.Fourteen
		}
		if team.Fifteen.RatingPlace != 0 {
			finalRatings["2015"] = team.Fifteen
		}
		if team.Sixteen.RatingPlace != 0 {
			finalRatings["2016"] = team.Sixteen
		}
		if team.Seventeen.RatingPlace != 0 {
			finalRatings["2017"] = team.Seventeen
		}
	}
	var value KeyedTeam
	value = KeyedTeam{
		team.Country,
		team.Academic,
		team.Id,
		team.Name,
		team.Aliases,
		finalRatings,
	}
	return value
}
