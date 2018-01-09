package engine

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/appengine"
)

// UpdateRankingsHandler handles any requests to <engine_hostname_or_ip>/rankings. In order to work correctly, the request must
// include the query string 'year'. The year query value will be used to request the respective year's rankings from
// ctftime.org. In order to use 'debug mode', which limits the maximum number of ranking pages requested to 2, set the query
// key 'debug' to true in the request. This handler operates in two phases.
//
// First Phase
//
// The first phase triggers multiple goroutines to parse and store ranking pages concurrently. By default, the maximum number
// of goroutines running at once is 10. To change the maximum number of goroutines running at once, modify the maxRoutines
// variable. This concurrent phase only requests pages that we have scraped before.
//
// Second Phase
//
// The second phase operates on the main thread and checks to see if a new rankings page exists. If a new page exists, it is
// parsed and stored in Firestore. Once we identify that phase two have reached the final page, the final page value is
// updated and stored in Firestore.
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
