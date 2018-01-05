package engine

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	"google.golang.org/appengine"
)

// UpdateCurrentRankingsHandler handles any requests to <engine_hostname_or_ip>/rankings. In order to work correctly, the
// request must include the query string 'year'. The year query value will be used to request the respective year's rankings
// from ctftime.org. In order to use 'debug mode', which limits the maximum number of ranking pages requested to 2, set the
// query key 'debug' to true in the request. This handler operates in two phases.
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
	var fbClient *firestore.Client
	var highestRankingsPage int
	var fbc FirebaseContext
	var debug bool
	newRankingsPage := true
	maxRoutines := 10
	guard := make(chan struct{}, maxRoutines)

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
		fbClient, err = Connect(token, r)
		if err != nil {
			http.Error(w, "Unable to connect to Firestore to acquire final page number", http.StatusInternalServerError)
			return
		}
		fbc = FirebaseContext{
			Ctx: appengine.NewContext(r), Fb: *fbClient,
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
		guard <- struct{}{}
		go func(i int) {
			FbClient, err := Connect(token, r)
			if err != nil {
				fmt.Printf("Unable to connect to Firestore for rankings page %d", i)
				<-guard
				return
			}
			fbc = FirebaseContext{
				Ctx: appengine.NewContext(r), Fb: *FbClient,
			}
			defer fbc.Fb.Close()

			rankingsUrl := fmt.Sprintf("https://ctftime.org/stats/%s?page=%d", year, i)
			response, err := Fetch(rankingsUrl)
			if err != nil {
				fmt.Println(err.Error())
				<-guard
				return
			}

			if err := ParseAndStoreRankings(response, i, year, fbc); err != nil {
				fmt.Println(err.Error())
				<-guard
				return
			}
			<-guard
		}(i)
	}

	// Phase Two
	for newRankingsPage {

		fbClient, err = Connect(token, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fbc = FirebaseContext{
			Ctx: appengine.NewContext(r), Fb: *fbClient,
		}

		rankingsUrl := fmt.Sprintf("https://ctftime.org/stats/%s?page=%d", year, highestRankingsPage)
		response, err := Fetch(rankingsUrl)

		// reached last page
		if err != nil {
			newRankingsPage = false
			UpdateLastRankingsPageNumber(fbc, year, highestRankingsPage)
			fbc.Fb.Close()
			w.Write([]byte("Finished loading contents"))
			return
		}

		// new page found
		if err = ParseAndStoreRankings(response, highestRankingsPage, year, fbc); err != nil {
			fmt.Println(err)
		}
		highestRankingsPage++
		fbc.Fb.Close()
	}
}
