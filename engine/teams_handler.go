package engine

import (
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
	var highestTeamId, offset int
	var fbc FirebaseContext
	var debug bool
	newTeam := true
	maxRoutines := 10
	guard := make(chan struct{}, maxRoutines)

	if debugQuery := r.URL.Query().Get("debug"); debugQuery == "true" {
		debug = true
	}

	token, err := GenerateToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !debug {
		offset = 1
		fbClient, err := Connect(token, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fbc = FirebaseContext{
			Ctx: appengine.NewContext(r), Fb: *fbClient,
		}
		if highestTeamId = GetLastTeamId(fbc); highestTeamId == 0 {
			http.Error(w, "Failed to acquire last team id value from Firestore.", http.StatusInternalServerError)
			fbc.Fb.Close()
			return
		}
		fbc.Fb.Close()
	} else {
		highestTeamId = 6
		offset = 570
	}

	// Phase One
	for i := offset; i < offset+highestTeamId; i++ {
		guard <- struct{}{}
		go func(teamId int) {
			fbClient, err := Connect(token, r)
			if err != nil {
				fmt.Println(err)
				fmt.Printf("Unable to connect to Firestore for team id %d\n", teamId)
				<-guard
				return
			}
			fbc = FirebaseContext{
				Ctx: appengine.NewContext(r), Fb: *fbClient,
			}
			defer fbc.Fb.Close()

			teamUrl := fmt.Sprintf("https://ctftime.org/team/%d", teamId)
			response, err := Fetch(teamUrl)
			if err != nil {
				fmt.Println(err.Error())
			} else if err := ParseAndStoreTeam(i, response, fbc); err != nil {
				fmt.Println(err.Error())
			}

			<-guard // must be last line of goroutine
		}(i)
	}

	for newTeam && !debug {
		fbClient, err := Connect(token, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fbc = FirebaseContext{
			Ctx: appengine.NewContext(r), Fb: *fbClient,
		}

		teamUrl := fmt.Sprintf("https://ctftime.org/team/%d", highestTeamId)
		response, err := Fetch(teamUrl)
		if err != nil {
			newTeam = false
			UpdateLastTeamId(fbc, highestTeamId)
		} else {
			err := ParseAndStoreTeam(highestTeamId, response, fbc)
			if err != nil {
				fmt.Println(err)
			}
			highestTeamId++
		}

		fbc.Fb.Close()
	}
	w.Write([]byte("Finished doing work"))
}
