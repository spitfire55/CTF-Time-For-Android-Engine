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

func storeRankings(rankings []Ranking, year int, fbc FirebaseContext) {
	for _, ranking := range rankings {
		collectionPath := fmt.Sprintf("%d_Rankings", year)
		_, err := fbc.fb.Collection(collectionPath).Doc(strconv.Itoa(ranking.Rank)).Set(fbc.ctx, ranking)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(fbc.w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func storeRankingsHash(hash string, pageNumDoc string, year int, fbc FirebaseContext) {
	collectionPath := fmt.Sprintf("%d_Rankings", year)
	fmt.Println("New doc at ", pageNumDoc)
	fmt.Println(hash)
	_, err := fbc.fb.Collection(collectionPath).Doc(pageNumDoc).Set(fbc.ctx, map[string]string{
		"hash": hash,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}

func hashDiff(resultsHash string, pageNumber int, year int, fbc FirebaseContext) (bool, error) {
	collectionPath := fmt.Sprintf("%d_Rankings", year)
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
		fmt.Println(hashDocValue)
		return true, nil
	} else {
		return false, nil
	}
}
