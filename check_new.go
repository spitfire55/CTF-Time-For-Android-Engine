package main

import (
	"strconv"
	"net/http"
	"fmt"
	"sync"
)

func updateAllTeams(fbc *FirebaseContext) {
	highestNode := int(getHighestNode(fbc))
	fmt.Println(highestNode)
	if highestNode != 0 {
		increment := 0
		interval := 100
		for increment < highestNode {
			var wg sync.WaitGroup
			for x := increment; (x < increment + interval) && (x < highestNode); x++ {
				wg.Add(1)
				go updateSingleTeam(x, &fbc.w, &fbc.r)
			}
			wg.Wait()
			break
		}
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
		} else {
			fmt.Println("booty")
		}
	}
	fmt.Println(node)
}

func checkNewTeam(node int, fbc *FirebaseContext) []byte {
	baseUrl := "https://ctftime.org/api/v1/teams/"
	return fetch(baseUrl+strconv.Itoa(node)+"/", fbc)
}
