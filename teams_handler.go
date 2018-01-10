// Copyright 2017-2018 Dale Lakes <spitfire@spitfy.re>. All rights reserved.
// Use of this source code is governed by the MIT license located in the LICENSE file.

package goctftime

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/appengine"
)

func UpdateTeamsHandler(w http.ResponseWriter, r *http.Request) {
	var highestTeamId int
	var debug bool
	newTeam := true
	ctx := appengine.WithContext(context.Background(), r)
	maxRoutines := 10
	guard := make(chan bool, maxRoutines)

	if debugQuery := r.URL.Query().Get("debug"); debugQuery == "true" {
		debug = true
	}

	token, err := GenerateToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !debug {
		fbc, err := NewFirebaseContext(ctx, token)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if highestTeamId = GetLastTeamId(fbc); highestTeamId == 0 {
			http.Error(w, "Failed to acquire last team id value from Firestore.", http.StatusInternalServerError)
			fbc.Fb.Close()
			return
		}

		fbc.Fb.Close()
	} else {
		highestTeamId = 11
	}

	// Phase One
	for i := 1; i < highestTeamId; i++ {
		guard <- true
		go func(teamId int) {
			defer func() { <-guard }()

			fbc, err := NewFirebaseContext(ctx, token)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			defer fbc.Fb.Close()

			teamUrl := fmt.Sprintf("https://ctftime.org/team/%d", teamId)
			if response, err := Fetch(teamUrl); err != nil {
				return
			} else if err = ParseAndStoreTeam(teamId, response, fbc); err != nil {
				fmt.Println(err.Error())
			}
		}(i)
	}

	for i := 0; i < maxRoutines; i++ {
		guard <- true
	}

	// Phase Two
	fbc, err := NewFirebaseContext(ctx, token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer fbc.Fb.Close()

	for newTeam && !debug {
		teamUrl := fmt.Sprintf("https://ctftime.org/team/%d", highestTeamId)
		response, err := Fetch(teamUrl)

		if err != nil {
			newTeam = false
			UpdateLastTeamId(fbc, highestTeamId)
			goto finish
		}
		if err := ParseAndStoreTeam(highestTeamId, response, fbc); err != nil {
			fmt.Println(err)
		}

		highestTeamId++
	}

finish:
	w.Write([]byte("Finished doing work"))
}
