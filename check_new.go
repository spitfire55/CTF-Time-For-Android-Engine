package main

import (
	"strconv"
	"net/http"
	"fmt"
	"sync"
	"time"
)

func updateAllTeams(fbc *FirebaseContext) {
	highestNode := int(getHighestNode(fbc))
	if highestNode != 0 { // make sure getHighestNode didn't fail
		increment := 0 // add 100 to each time in inner for loop
		interval := 20 // stays at 100 or however many should be scanned in instance
		for increment < highestNode {
			var wg sync.WaitGroup
			fmt.Println(increment)
			for x := increment; (x < increment + interval) && (x < highestNode); x++ {
				wg.Add(1)
				go func(x int, rw http.ResponseWriter, req http.Request) {
					defer wg.Done()
					updateSingleTeam(x, &rw, &req)
				}(x, fbc.w, fbc.r)
			}
			wg.Wait()
			time.Sleep(1 * time.Second)
			increment += 20 // added value must equal value of interval
		}
		//TODO: Implement check_new feature here
	} else {
		http.Error(fbc.w,
			"Unable to acquire highest node",
				http.StatusInternalServerError)
	}
}

func updateSingleTeam(node int, w *http.ResponseWriter, r *http.Request) {
	FbClient, ctx := connect()
	if FbClient != nil && ctx != nil {
		// create new firebase context w/ same ResponseWriter & Request
		fbc := &FirebaseContext{
			*w, *r, http.Client{}, ctx, *FbClient,
		}
		body := checkNewTeam(node, fbc)
		if body != nil {
			bodyKeyed := getSingleTeam(body)
			saveNewTeam(node, bodyKeyed, fbc)
			fbc.fb.Close()
		}
	} else {
		fmt.Printf("Failed on %d\n", node)
	}
}

func checkNewTeam(node int, fbc *FirebaseContext) []byte {
	baseUrl := "https://ctftime.org/api/v1/teams/"
	return fetch(baseUrl+strconv.Itoa(node)+"/", fbc)
}
