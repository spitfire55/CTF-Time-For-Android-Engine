package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
	"net/http"
)

func connect(token option.ClientOption, r *http.Request) *firestore.Client {
	ctx := appengine.NewContext(r)
	client, err := firestore.NewClient(ctx, "ctf-time-for-android", token)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	} else {
		return client
	}
}

func generateToken() option.ClientOption {
	return option.WithCredentialsFile(os.Getenv("CTF_TIME_KEY"))
}

func getLastTeamId(fbc FirebaseContext) int {
	lastPageNumberDoc, _ := fbc.fb.Collection("Teams").Doc("LastTeamId").Get(fbc.ctx)
	lastPageNumber, _ := lastPageNumberDoc.DataAt("lastPageNumber")
	if lastPageNumberInt, ok := lastPageNumber.(int64); ok {
		return int(lastPageNumberInt)
	} else {
		return 0
	}
}

func updateLastTeamId(fbc FirebaseContext, newPageNumber int) {
	_, err := fbc.fb.Collection("Teams").Doc("LastPageNumber").Set(fbc.ctx, map[string]int{
		"lastPageNumber": newPageNumber,
	})
	if err != nil {
		fmt.Println(err)
	}
}

func storeTeam(fbc FirebaseContext, team Team) {

}

func getLastRankingsPageNumber(fbc FirebaseContext, year string) int {
	collectionString := fmt.Sprintf("%s_Rankings", year)
	lastPageNumberDoc, _ := fbc.fb.Collection(collectionString).Doc("LastPageNumber").Get(fbc.ctx)
	lastPageNumber, _ := lastPageNumberDoc.DataAt("lastPageNumber")
	if lastPageNumberInt, ok := lastPageNumber.(int64); ok {
		return int(lastPageNumberInt)
	} else {
		return 0
	}
}

func updateLastRankingsPageNumber(fbc FirebaseContext, year string, newPageNumber int) {
	collectionString := fmt.Sprintf("%s_Rankings", year)
	_, err := fbc.fb.Collection(collectionString).Doc("LastPageNumber").Set(fbc.ctx, map[string]int{
		"lastPageNumber": newPageNumber,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}

func storeRankings(rankings []Ranking, year string, fbc FirebaseContext) {
	collectionString := fmt.Sprintf("%s_Rankings", year)
	for _, ranking := range rankings {
		_, err := fbc.fb.Collection(collectionString).Doc(strconv.Itoa(ranking.Rank)).Set(fbc.ctx, ranking)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func storeRankingsHash(hash string, pageNumDoc string, year string, fbc FirebaseContext) {
	collectionString := fmt.Sprintf("%s_Rankings", year)
	_, err := fbc.fb.Collection(collectionString).Doc(pageNumDoc).Set(fbc.ctx, map[string]string{
		"hash": hash,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}

func rankingsHashDiff(resultsHash string, pageNumber int, year string, fbc FirebaseContext) (bool, error) {
	collectionPath := fmt.Sprintf("%s_Rankings", year)
	pageNumDoc := fmt.Sprintf("Page%dHash", pageNumber)
	hashDoc, err := fbc.fb.Collection(collectionPath).Doc(pageNumDoc).Get(fbc.ctx)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") { // document doesn't exist, create it
			return true, err
		}
		return false, err
	}
	hashDocValue, err := hashDoc.DataAt("hash")
	if err != nil {
		return false, err
	}
	if resultsHash != hashDocValue {
		return true, nil
	} else {
		return false, nil
	}
}
