package engine

import (
	"fmt"
	"strconv"
	"strings"
)

// GetLastRankingsPageNumber retrieves the final rankings page number for a certain year from the Firestore database.
// Specifically, it queries the rankings collection of the specific year for the lastPageNumber field in the LastPageNumber
// document. If the collection, document, or field is not found or the value is not an integer, the return value is zero.
func GetLastRankingsPageNumber(fbc FirebaseContext, year string) int {
	collectionString := fmt.Sprintf("%s_Rankings", year)
	lastPageNumberDoc, err := fbc.Fb.Collection(collectionString).Doc("LastPageNumber").Get(fbc.Ctx)
	if err != nil {
		return 0
	}
	lastPageNumber, err := lastPageNumberDoc.DataAt("lastPageNumber")
	if err != nil {
		return 0
	}
	if lastPageNumberInt, ok := lastPageNumber.(int64); ok {
		return int(lastPageNumberInt)
	} else {
		return 0
	}
}

// UpdateLastRankingsPageNumber updates the final rankings page number for a certain year to the Firestore database.
// The upload path is: <year>_Rankings/LastPageNumber/lastPageNumber
func UpdateLastRankingsPageNumber(fbc FirebaseContext, year string, newPageNumber int) {
	collectionString := fmt.Sprintf("%s_Rankings", year)
	_, err := fbc.Fb.Collection(collectionString).
		Doc("LastPageNumber").Set(fbc.Ctx, map[string]int{
		"lastPageNumber": newPageNumber,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}

// StoreRankings stores an array of Ranking types to Firestore database. A new document is either created or updated for each
// Ranking. The Document key is the Rank field of the Ranking struct. If any of the individual Ranking structs fails to
// upload, the subsequent rankings will not be uploaded and an error is returned.
// The upload path for a ranking struct is: <year>_Rankings/<rank>
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

// StoreRankingsHash stores a sha256 hash of an array of fifty rankings. The fifty rankings corresponds to the fifty rankings
// displayed by ctftime.org when a specific page is requested.
// Page 1 = 1st - 50th place teams, page 2 = 51st - 100th place teams, etc. If the page requested is the final page, there
// might be less than 50 teams displayed on the page.
// The upload path for a hash is: <year>_Rankings/<pageNumber>/hash
func StoreRankingsHash(hash string, pageNumDoc string, year string, fbc FirebaseContext) {
	collectionString := fmt.Sprintf("%s_Rankings", year)
	_, err := fbc.Fb.Collection(collectionString).Doc(pageNumDoc).Set(fbc.Ctx, map[string]string{
		"hash": hash,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}

// RankingsHashDiff compares the hash computed from a rankings page on ctftime.org to the hash stored of that page in the
// Firestore database. If the hashes are not equal, then that rankings page has changed in some way. This check is used to
// prevent unnecessary writes to the Firestore database if a rankings page has not changed between scrape iterations.
// If the page is new, a new document is created to store the hash for that page number.
// NOTE: The hash is computed from an array of Ranking types (i.e. after parsing), not from the response body of the request.
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
