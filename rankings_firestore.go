// Copyright 2017-2018 Dale Lakes <spitfire@spitfy.re>. All rights reserved.
// Use of this source code is governed by the MIT license located in the LICENSE file.

package goctftime

import (
	"fmt"
	"strconv"
	"strings"
)

func GetLastRankingsPageNumber(fbc FirebaseContext, year string) int {
	collectionString := fmt.Sprintf("%s_Rankings", year)
	lastPageNumberDoc, err := fbc.Fb.Collection(collectionString).Doc("LastRankingsPage").Get(fbc.Ctx)
	if err != nil {
		return 0
	}
	lastPageNumber, err := lastPageNumberDoc.DataAt("lastRankingsPage")
	if err != nil {
		return 0
	}
	if lastPageNumberInt, ok := lastPageNumber.(int64); ok {
		return int(lastPageNumberInt)
	} else {
		return 0
	}
}

func UpdateLastRankingsPageNumber(fbc FirebaseContext, year string, newPageNumber int) {
	collectionString := fmt.Sprintf("%s_Rankings", year)
	_, err := fbc.Fb.Collection(collectionString).Doc("LastRankingsPage").Set(fbc.Ctx, map[string]int{
		"lastRankingsPage": newPageNumber,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}

func StoreRankings(rankings []Ranking, year string, fbc FirebaseContext) error {
	collectionString := fmt.Sprintf("%s_Rankings", year)
	for _, ranking := range rankings {
		_, err := fbc.Fb.Collection(collectionString).Doc(strconv.Itoa(ranking.Rank)).Set(fbc.Ctx, ranking)
		if err != nil {
			return err
		}
	}
	return nil
}

// Stores a hash of fifty rankings, which is the number of rankings displayed on a single page
func StoreRankingsHash(hash string, pageNumDoc string, year string, fbc FirebaseContext) {
	collectionString := fmt.Sprintf("%s_Rankings", year)
	_, err := fbc.Fb.Collection(collectionString).Doc(pageNumDoc).Set(fbc.Ctx, map[string]string{
		"hash": hash,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}

func RankingsHashDiff(resultsHash string, pageNumber int, year string, fbc FirebaseContext) (bool, error) {
	collectionPath := fmt.Sprintf("%s_Rankings", year)
	pageNumDoc := fmt.Sprintf("Page%dHash", pageNumber)
	hashDoc, err := fbc.Fb.Collection(collectionPath).Doc(pageNumDoc).Get(fbc.Ctx)
	if err != nil {
		// Document doesn't exists, so error is nil to create it
		if strings.Contains(err.Error(), "NotFound") {
			return true, nil
		}
		// Some other error, so return error
		return false, err
	}
	hashDocValue, err := hashDoc.DataAt("hash")
	if err != nil {
		// Document doesn't have hash field or we can't read it, so return error
		return false, err
	}
	if resultsHash != hashDocValue {
		// Hashes are different, so no error
		return true, nil
	} else {
		// Hashes are same, so no error but return false to prevent unnecessary write
		return false, nil
	}
}
