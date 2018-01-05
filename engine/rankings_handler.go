package engine

import (
	"fmt"
	"net/http"
	"strconv"

	"cloud.google.com/go/firestore"
	"google.golang.org/appengine"
)

func UpdateCurrentRankingsHandler(w http.ResponseWriter, r *http.Request) {
	var fbClient *firestore.Client
	token := GenerateToken()
	year := r.URL.Query().Get("year")
	fbClient = Connect(token, r)
	var highestNode int
	var fbc FirebaseContext
	if fbClient != nil {
		fbc = FirebaseContext{
			W: w, R: *r, C: http.Client{}, Ctx: appengine.NewContext(r), Fb: *fbClient,
		}
		highestNode = GetLastRankingsPageNumber(fbc, year)
		fbc.Fb.Close()
	} else {
		highestNode = 0
	}
	newRankingsPage := true
	maxRoutines := 10
	guard := make(chan struct{}, maxRoutines)
	for i := 1; i < highestNode; i++ {
		guard <- struct{}{}
		fbClient = Connect(token, r)
		if fbClient != nil {
			fbc = FirebaseContext{
				w, *r, http.Client{}, appengine.NewContext(r), *fbClient,
			}
			go func(i int) {
				FbClient := Connect(token, r)
				if FbClient != nil {
					fbc = FirebaseContext{
						w, *r, http.Client{}, appengine.NewContext(r), *FbClient,
					}
					rankingsUrl := fmt.Sprintf("https://ctftime.org/stats/%s?page=%d", year, i)
					response, err := Fetch(rankingsUrl, fbc)
					if err == nil {
						err := ParseAndStoreRankings(response, i, year, fbc)
						if err != nil {
							fmt.Println(err)
						} else {
							fmt.Println(strconv.Itoa(i) + " is complete")
						}
					}
					fbc.Fb.Close()
					<-guard
				} else {
					fmt.Println("Failed to connect to Firestore")
				}
			}(i)
		}
	}
	for newRankingsPage {
		fbClient = Connect(token, r)
		if fbClient != nil {
			fbc = FirebaseContext{
				W: w, R: *r, C: http.Client{}, Ctx: appengine.NewContext(r), Fb: *fbClient,
			}
			rankingsUrl := fmt.Sprintf("https://ctftime.org/stats/%s?page=%d", year, highestNode)
			response, err := Fetch(rankingsUrl, fbc)
			if err != nil {
				newRankingsPage = false
				UpdateLastRankingsPageNumber(fbc, year, highestNode)
				goto finish
			}
			err = ParseAndStoreRankings(response, highestNode, year, fbc)
			if err != nil {
				fmt.Println(err)
			}
			highestNode++
			fbc.Fb.Close()
		} else {
			http.Error(w,
				"Failed to connect to Firestore",
				http.StatusInternalServerError)
			goto finish
		}
	}
	finish:
		w.Write([]byte("Finished loading contents"))
}
