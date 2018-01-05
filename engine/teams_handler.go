package engine

import (
	"fmt"
	"log"
	"net/http"

	"google.golang.org/appengine"
)


func UpdateTeamsHandler(w http.ResponseWriter, r *http.Request) {
	debug := true
	token := GenerateToken()
	fbClient := Connect(token, r)
	var highestTeamId int
	var fbc FirebaseContext
	if fbClient != nil {
		fbc = FirebaseContext{
			W: w, R: *r, C: http.Client{}, Ctx: appengine.NewContext(r), Fb: *fbClient,
		}
		if debug {
			highestTeamId = 5
		} else {
			highestTeamId = GetLastTeamId(fbc)
		}
		fbc.Fb.Close()
	} else {
		log.Fatal("FbClient is nil")
	}
	newTeam := true
	maxRoutines := 5
	guard := make(chan struct{}, maxRoutines)
	for i := 1; i < highestTeamId; i++ {
		guard <- struct{}{}
		fbClient := Connect(token, r)
		if fbClient != nil {
			fbc = FirebaseContext{
				W: w, R: *r, C: http.Client{}, Ctx: appengine.NewContext(r), Fb: *fbClient,
			}
			go func(teamId int, fbc FirebaseContext) {
				team, err := ParseTeam(teamId)
				if err != nil {
					fmt.Println(err.Error())
					return // something went wrong, don't call the store function
				}
				StoreTeam(fbc, team)
				fbc.Fb.Close()
				<-guard // must be last line of goroutine
			}(i, fbc)
		}
	}
	for newTeam && !debug {
		fbClient = Connect(token, r)
		if fbClient != nil {
			fbc = FirebaseContext{
				W: w, R: *r, C: http.Client{}, Ctx: appengine.NewContext(r), Fb: *fbClient,
			}
			teamUrl := fmt.Sprintf("https://ctftime.org/team/%d", highestTeamId)
			response := Fetch(teamUrl, fbc)
			if response != nil {
				team, err := ParseTeam(highestTeamId)
				if err != nil {
					fmt.Println(err)
				}
				StoreTeam(fbc, team)
				highestTeamId++
			} else {
				newTeam = false
				UpdateLastTeamId(fbc, highestTeamId)
			}
			fbc.Fb.Close()
		}
	}
}
