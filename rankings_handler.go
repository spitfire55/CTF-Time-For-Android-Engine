// Copyright 2017-2018 Dale Lakes <spitfire@spitfy.re>. All rights reserved.
// Use of this source code is governed by the MIT license located in the LICENSE file.

package goctftime

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/appengine"
)

func UpdateRankingsHandler(w http.ResponseWriter, r *http.Request) {
	var highestRankingsPage int
	var debug bool
	newRankingsPage := true
	ctx := appengine.WithContext(context.Background(), r)
	maxRoutines := 10
	guard := make(chan bool, maxRoutines)

	if debugQuery := r.URL.Query().Get("debug"); debugQuery == "true" {
		debug = true
	}

	year := r.URL.Query().Get("year")
	if year == "" {
		http.Error(w, "Query key 'year' is not specified", http.StatusInternalServerError)
		return
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

		highestRankingsPage = GetLastRankingsPageNumber(fbc, year)
		if highestRankingsPage == 0 {
			http.Error(w, "Failed to acquire last rankings page value from Firestore.", http.StatusInternalServerError)
			fbc.Fb.Close()
			return
		}

		fbc.Fb.Close()
	} else {
		highestRankingsPage = 2
	}

	// Phase One
	for i := 1; i < highestRankingsPage; i++ {
		guard <- true
		go func(teamId int) {
			defer func() { <-guard }()

			fbc, err := NewFirebaseContext(ctx, token)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			defer fbc.Fb.Close()

			rankingsUrl := fmt.Sprintf("https://ctftime.org/stats/%s?page=%d", year, teamId)
			if response, err := Fetch(rankingsUrl); err != nil {
				fmt.Println(err.Error())
			} else if err := ParseAndStoreRankings(response, teamId, year, fbc); err != nil {
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

	for newRankingsPage && !debug {
		rankingsUrl := fmt.Sprintf("https://ctftime.org/stats/%s?page=%d", year, highestRankingsPage)
		response, err := Fetch(rankingsUrl)

		if err != nil {
			newRankingsPage = false
			UpdateLastRankingsPageNumber(fbc, year, highestRankingsPage)
			goto finish
		}

		if err = ParseAndStoreRankings(response, highestRankingsPage, year, fbc); err != nil {
			fmt.Println(err)
		}

		highestRankingsPage++
	}

finish:
	w.Write([]byte("Finished loading contents"))
}
