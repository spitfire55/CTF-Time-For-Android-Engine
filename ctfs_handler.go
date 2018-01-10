// Copyright 2017-2018 Dale Lakes <spitfire@spitfy.re>. All rights reserved.
// Use of this source code is governed by the MIT license located in the LICENSE file.

package goctftime

import (
	"context"
	"fmt"
	"google.golang.org/appengine"
	"net/http"
)

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
