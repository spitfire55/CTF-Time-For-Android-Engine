package main

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"fmt"
	"strconv"
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
		http.Error(fbc.w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	if resp.StatusCode != 200 {
		return nil
	}
	fmt.Println(url + resp.Status)
	return resp
}

func updateCurrentRankingsHandler(w http.ResponseWriter, r *http.Request) {
	token := generateToken()
	FbClient, ctx := connect(token)
	if FbClient != nil && ctx != nil {
		fbc := FirebaseContext{
			w, *r, http.Client{}, ctx, *FbClient,
		}
		year := r.URL.Query().Get("year")
		highestNode := getLastPageNumber(fbc, year)
		success := true
		for i := 1; i < highestNode; i++ {
			go func(i int) {
				fmt.Println("Working on " + strconv.Itoa(i))
				rankingsUrl := fmt.Sprintf("https://ctftime.org/stats/%s?page=%d", year, i)
				response := fetch(rankingsUrl, fbc)
				if response != nil {
					err := parseAndStoreRankings(response, i, year, fbc)
					if err != nil {
						fmt.Println(err)
					}
				}
				return
			}(i)
		}
		for success {
			rankingsUrl := fmt.Sprintf("https://ctftime.org/stats/%s?page=%d", year, highestNode)
			response := fetch(rankingsUrl, fbc)
			if response != nil {
				fmt.Println(response.Body)
				err := parseAndStoreRankings(response, highestNode, year, fbc)
				if err != nil {
					http.Error(fbc.w, err.Error(), http.StatusInternalServerError)
				}
				highestNode++
			} else {
				success = false
				updateLastPageNumber(fbc, year, highestNode)
			}
		}
	} else {
		http.Error(w,
			"Failed to connect to Firestore",
			http.StatusInternalServerError)
	}
	w.Write([]byte("Finished loading contents"))
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {}

func main() {
	http.HandleFunc("/favicon.ico", defaultHandler)
	http.HandleFunc("/", updateCurrentRankingsHandler)
	http.ListenAndServe("localhost:8080", nil)
}
