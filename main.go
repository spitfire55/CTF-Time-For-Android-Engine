package main

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"fmt"
	"strconv"
	"google.golang.org/appengine"
	"sync"
)

type FirebaseContext struct {
	w   http.ResponseWriter
	r   http.Request
	c   http.Client      // client used to GET from ctftime.org
	ctx context.Context  // context used in connection to Firestore
	fb  firestore.Client // client used in connection to Firestore
}

func fetch(url string, fbc FirebaseContext) *http.Response {
	resp, err := fbc.c.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if resp.StatusCode != 200 {
		return nil
	}
	return resp
}

func updateCurrentRankingsHandler(w http.ResponseWriter, r *http.Request) {
	token := generateToken()
	year := r.URL.Query().Get("year")
	FbClient := connect(token, r)
	var highestNode int
	if FbClient != nil {
		fbc := FirebaseContext{
			w, *r, http.Client{}, appengine.NewContext(r), *FbClient,
		}
		highestNode = getLastPageNumber(fbc, year)
		fbc.fb.Close()
	} else {
		highestNode = 0
	}
	success := true
	for i := 1; i < highestNode; i += 10 {
		var wg sync.WaitGroup
		wg.Add(10)
		for j := i; j < i + 10 && j < highestNode; j++ {
			go func(j int) {
				defer wg.Done()
				FbClient := connect(token, r)
				if FbClient != nil {
					fbc := FirebaseContext{
						w, *r, http.Client{}, appengine.NewContext(r), *FbClient,
					}
					rankingsUrl := fmt.Sprintf("https://ctftime.org/stats/%s?page=%d", year, j)
					response := fetch(rankingsUrl, fbc)
					if response != nil {
						err := parseAndStoreRankings(response, j, year, fbc)
						if err != nil {
							fmt.Println(err)
						}
					}
					fbc.fb.Close()
					fbc.ctx.Done()
					fmt.Println(strconv.Itoa(j) + " is complete")
				} else {
					http.Error(w,
						"Failed to connect to Firestore",
						http.StatusInternalServerError)
				}
			}(j)
		}
		wg.Wait()
	}
	for success {
		FbClient := connect(token, r)
		if FbClient != nil {
			fbc := FirebaseContext{
				w, *r, http.Client{}, appengine.NewContext(r), *FbClient,
			}
			rankingsUrl := fmt.Sprintf("https://ctftime.org/stats/%s?page=%d", year, highestNode)
			response := fetch(rankingsUrl, fbc)
			if response != nil {
				err := parseAndStoreRankings(response, highestNode, year, fbc)
				if err != nil {
					fmt.Println(err)
				}
				highestNode++
			} else {
				success = false
				updateLastPageNumber(fbc, year, highestNode)
			}
			fbc.fb.Close()
		} else {
			http.Error(w,
			"Failed to connect to Firestore",
			http.StatusInternalServerError)
		}
	}
	w.Write([]byte("Finished loading contents"))
}

func defaultHandler(_ http.ResponseWriter, _ *http.Request) {}

func main() {
	http.HandleFunc("/rankings", updateCurrentRankingsHandler)
	http.HandleFunc("/", defaultHandler)
	http.ListenAndServe("localhost:8080", nil)
}
