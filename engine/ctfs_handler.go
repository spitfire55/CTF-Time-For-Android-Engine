package engine

import (
	"context"
	"fmt"
	"google.golang.org/appengine"
	"net/http"
)

// UpdateCtfsHandler handles any requests to <engine_hostname_or_ip>/ctfs. In order to use 'debug mode',
// which limits the maximum number of ctf pages requested to 10, set the query key 'debug' to true in the request. This
// handler operates in two phases.
//
// First Phase
//
// The first phase triggers multiple goroutines to parse and store ctf pages concurrently. By default, the maximum number
// of goroutines running at once is 10. To change the maximum number of goroutines running at once, modify the maxRoutines
// variable. This concurrent phase only requests pages that we have scraped before.
//
// Second Phase
//
// The second phase operates on a single thread and checks to see if a new ctf page exists. If a new page exists, it is
// parsed and stored in Firestore. Once we identify that phase two has reached the final page, the final page value is
// updated and stored in Firestore.
func UpdateCtfsHandler(w http.ResponseWriter, r *http.Request) {
	var highestCtfId int
	var debug bool
	ctx := appengine.WithContext(context.Background(), r)
	newCtf := true
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

		highestCtfId = GetLastCtfId(fbc)
		if highestCtfId == 0 {
			http.Error(w, "Failed to acquire last rankings page value from Firestore.", http.StatusInternalServerError)
			fbc.Fb.Close()
			return
		}

		fbc.Fb.Close()
	} else {
		highestCtfId = 11
	}

	// Phase One
	for i := 1; i < highestCtfId; i++ {
		guard <- true
		go func(ctfId int) {
			defer func() { <-guard }()

			fbc, err := NewFirebaseContext(ctx, token)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			defer fbc.Fb.Close()

			ctfUrl := fmt.Sprintf("https://ctftime.org/ctf/%d", ctfId)
			if response, err := Fetch(ctfUrl); err != nil {
				fmt.Println(err.Error())
			} else if err := ParseAndStoreCtf(ctfId, response, fbc); err != nil {
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

	for newCtf && !debug {
		teamUrl := fmt.Sprintf("https://ctftime.org/ctf/%d", highestCtfId)
		response, err := Fetch(teamUrl)

		if err != nil {
			newCtf = false
			UpdateLastCtfId(fbc, highestCtfId)
			goto finish
		}

		if err := ParseAndStoreCtf(highestCtfId, response, fbc); err != nil {
			fmt.Println(err)
		}

		highestCtfId++
	}

finish:
	w.Write([]byte("Finished doing work"))
}
