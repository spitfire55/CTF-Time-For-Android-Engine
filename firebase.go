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

func storeRankings(rankings []Ranking, fbc FirebaseContext) {
	for _, ranking := range rankings {
		_, err := fbc.fb.Collection("2017_Rankings").Doc(strconv.Itoa(ranking.Rank)).Set(fbc.ctx, ranking)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(fbc.w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func storeRankingsHash(hash string, pageNumDoc string, fbc FirebaseContext) {
	fmt.Println("New doc at ", pageNumDoc)
	fmt.Println(hash)
	_, err := fbc.fb.Collection("2017_Rankings").Doc(pageNumDoc).Set(fbc.ctx, map[string]string{
		"hash": hash,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}

func hashDiff(resultsHash string, pageNumber int, fbc FirebaseContext) (bool, error) {
	pageNumDoc := fmt.Sprintf("Page%dHash", pageNumber)
	hashDoc, err := fbc.fb.Collection("2017_Rankings").Doc(pageNumDoc).Get(fbc.ctx)
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
