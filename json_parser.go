package main

import (
	"encoding/json"
	"log"
	"strconv"
)

type AllRankings struct {
	Sixteen   []Rankings `json:"2016"`
	Seventeen []Rankings `json:"2017"`
	Eleven    interface{} `json:"-"`
	Twelve    []Rankings `json:"2012"`
	Thirteen  interface{} `json:"-"`
	Fourteen  interface{} `json:"-"`
	Fifteen   []Rankings `json:"2015"`

}

type CurrentRankings struct {
	Rankings []Rankings `json:"2017"`
}

type Rankings struct {
	TeamName string  `json:"team_name"`
	Points   float64 `json:"points"`
	Id       int     `json:"team_id"`
}

// will be used to store only the rankings from AllRankings that actually exist (ignore json:"-" values)
type ValidRankings [][]Rankings

type KeyedRankingsYear map[string]KeyedRankingValue
type KeyedRankingsAll map[string]KeyedRankingsYear

type KeyedRankingValue struct {
	TeamName string `json:"team_name"`
	Points float64 `json:"points"`
}
func getAllRankings(jsonStream []byte) KeyedRankingsAll {

	var results AllRankings
	err := json.Unmarshal(jsonStream, &results)
	if err != nil {
		log.Fatal(err)
	}
	// store the valid top 10 in its own slice, which is a 2-d map of rankings (rows are years, columns are Rankings structs)
	validRankings := ValidRankings{results.Twelve, results.Fifteen, results.Sixteen, results.Seventeen}
	// store just the years as slice of strings.
	validRankingsYears := []string{"2012", "2015", "2016", "2017"}

	// initialize an empty map that will eventually contain contents to store in Firebase. Key = year, value = KeyedRankingsYears map
	var keyRankings = make(KeyedRankingsAll, len(validRankings))
	for i, yearRankings := range validRankings {
		// initialize inner map. Key = team id, value = KeyedRankingValue interface
		keyRankings[validRankingsYears[i]] = make(map[string]KeyedRankingValue, len(validRankings[i]))
		for _, ranking := range yearRankings {
			// indices of validRankingsYears align to the order in which validRankings are stored and thus, contains the correct
			// year to use as a key value for the outer maps. 0 = 2012, 1 = 2015, etc.
			keyRankings[validRankingsYears[i]][strconv.Itoa(ranking.Id)] = KeyedRankingValue{ranking.TeamName, ranking.Points}
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

	var keyCurrentRankings = make(map[string]KeyedRankingValue, len(results.Rankings))
	for _, ranking := range results.Rankings {
		keyCurrentRankings[strconv.Itoa(ranking.Id)] = KeyedRankingValue{ranking.TeamName, ranking.Points}
	}
	return keyCurrentRankings
}