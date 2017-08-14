package main

import (
	"encoding/json"
	"log"
)

type Rankings struct {
	teamName string `json:"team_name"`
	points float64 `json:"points"`
	id int `json:"team_id"`
}

func decode(jsonStream *[]byte) []Rankings {
	var results map[string][]interface{}
	err := json.Unmarshal(*jsonStream, &results); if err != nil {
		log.Fatal(err)
	}
	var rankings []Rankings
	for _, v := range results {
		if len(v) > 0 {
			for _, rankV := range v {
				var rankStruct Rankings
				rank := rankV.(map[string]interface{})
				rankStruct.teamName = rank["team_name"].(string)
				rankStruct.points = rank["points"].(float64)
				rankStruct.id = int(rank["team_id"].(float64))
				rankings = append(rankings, rankStruct)
			}
		}
	}
	return rankings
}
