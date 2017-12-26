package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

func connect(token option.ClientOption) (*firestore.Client, context.Context) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "ctf-time-for-android", token)
	if err != nil {
		fmt.Println(err.Error())
		return nil, nil
	} else {
		return client, ctx
	}
}

func generateToken() option.ClientOption {
	return option.WithCredentialsFile(os.Getenv("CTF_TIME_KEY"))
}

func getLastPageNumber(fbc FirebaseContext, year string) int {
	collectionString := fmt.Sprintf("%s_Rankings", year)
	lastPageNumberDoc, _ := fbc.fb.Collection(collectionString).Doc("LastPageNumber").Get(fbc.ctx)
	lastPageNumber, _ := lastPageNumberDoc.DataAt("lastPageNumber")
	if lastPageNumberInt, ok := lastPageNumber.(int64); ok {
		return int(lastPageNumberInt)
	} else {
		return 0
	}
}

func updateLastPageNumber(fbc FirebaseContext, year string, newPageNumber int) {
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
			http.Error(fbc.w, err.Error(), http.StatusInternalServerError)
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

func hashDiff(resultsHash string, pageNumber int, year string, fbc FirebaseContext) (bool, error) {
	pageNumDoc := fmt.Sprintf("Page%dHash", pageNumber)
	collectionString := fmt.Sprintf("%s_Rankings", year)
	hashDoc, err := fbc.fb.Collection(collectionString).Doc(pageNumDoc).Get(fbc.ctx)
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
		fmt.Println(hashDocValue)
		return true, nil
	} else {
		return false, nil
	}
}
