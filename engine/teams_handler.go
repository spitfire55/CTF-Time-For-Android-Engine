package engine

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/appengine"
)

// UpdateTeamsHandler handles any requests to <engine_hostname_or_ip>/teams. In order to use 'debug mode', which limits the
// maximum number of teams requested to 5, set the query key 'debug' to true in the request.
// The handler operates in two phases.
//
// First Phase
//
// The first phase triggers multiple goroutines to parse and store teams concurrently. By default, the maximum number of
// goroutines is 5. To change the maximum number of goroutines running at once, modify the maxRoutines variable.
// This concurrent phase only requests pages that we have scraped before.
//
// Second Phase
//
// The second phase operates on the main thread and checks to see if a new team exists. If a new team exists, it is parsed
// and stored in Firestore. Once phase two has reached the final page, the final page value is updated and stored in
// Firestore.
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
				fmt.Println(err.Error())
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
