package main

import (
	"encoding/json"
	"log"
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

type Rankings struct {
	TeamName string  `json:"team_name"`
	Points   float64 `json:"points"`
	Id       int     `json:"team_id"`
}

func getTeamRankings(jsonStream []byte) AllRankings {

	var results AllRankings
	err := json.Unmarshal(jsonStream, &results)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", results)
	return results

}
