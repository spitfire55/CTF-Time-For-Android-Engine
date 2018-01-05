package engine

import (
	"fmt"
	"strconv"
	"strings"
)

func GetLastRankingsPageNumber(fbc FirebaseContext, year string) int {
	collectionString := fmt.Sprintf("%s_Rankings", year)
	lastPageNumberDoc, _ := fbc.Fb.Collection(collectionString).Doc("LastPageNumber").Get(fbc.Ctx)
	lastPageNumber, _ := lastPageNumberDoc.DataAt("lastPageNumber")
	if lastPageNumberInt, ok := lastPageNumber.(int64); ok {
		return int(lastPageNumberInt)
	} else {
		return 0
	}
}

func UpdateLastRankingsPageNumber(fbc FirebaseContext, year string, newPageNumber int) {
	collectionString := fmt.Sprintf("%s_Rankings", year)
	_, err := fbc.Fb.Collection(collectionString).Doc("LastPageNumber").Set(fbc.Ctx, map[string]int{
		"lastPageNumber": newPageNumber,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}

func StoreRankings(rankings []Ranking, year string, fbc FirebaseContext) {
	collectionString := fmt.Sprintf("%s_Rankings", year)
	for _, ranking := range rankings {
		_, err := fbc.Fb.Collection(collectionString).Doc(strconv.Itoa(ranking.Rank)).Set(fbc.Ctx, ranking)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

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
