package engine

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"google.golang.org/appengine"
	"net/http"
)

func UpdateCtfsHandler(w http.ResponseWriter, r *http.Request) {
	var fbClient *firestore.Client
	var highestCtfId int
	var fbc FirebaseContext
	var debug bool
	newCtf := true
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
		fbClient, err = Connect(token, r)
		if err != nil {
			http.Error(w, "Unable to connect to Firestore to acquire final page number", http.StatusInternalServerError)
			return
		}
		fbc = FirebaseContext{
			Ctx: appengine.NewContext(r), Fb: *fbClient,
		}
		highestCtfId = GetLastCtfId(fbc)
		if highestCtfId == 0 {
			http.Error(w, "Failed to acquire last rankings page value from Firestore.", http.StatusInternalServerError)
			fbc.Fb.Close()
			return
		}
		fbc.Fb.Close()
	} else {
		highestCtfId = 2
	}

	// Phase One
	for i := 1; i < highestCtfId; i++ {
		guard <- struct{}{}
		go func(i int) {
			FbClient, err := Connect(token, r)
			if err != nil {
				fmt.Printf("Unable to connect to Firestore for ctf id %d", i)
				<-guard
				return
			}
			fbc = FirebaseContext{
				Ctx: appengine.NewContext(r), Fb: *FbClient,
			}
			ctfUrl := fmt.Sprintf("https://ctftime.org/ctf/%d", i)
			response, err := Fetch(ctfUrl)
			if err != nil {
				fmt.Println(err.Error())
			} else if err := ParseAndStoreCtf(i, response, fbc); err != nil {
				fmt.Println("wtf" + err.Error())
			}
			fbc.Fb.Close()
			<-guard
		}(i)
	}

	// Phase Two
	for newCtf && !debug {
		fbClient, err := Connect(token, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fbc = FirebaseContext{
			Ctx: appengine.NewContext(r), Fb: *fbClient,
		}

		teamUrl := fmt.Sprintf("https://ctftime.org/ctf/%d", highestCtfId)
		response, err := Fetch(teamUrl)
		if err != nil {
			newCtf = false
			UpdateLastCtfId(fbc, highestCtfId)
		} else {
			err := ParseAndStoreCtf(highestCtfId, response, fbc)
			if err != nil {
				fmt.Println(err.Error())
			}
			highestCtfId++
		}
		fbc.Fb.Close()
	}
	w.Write([]byte("Finished doing work"))
}
