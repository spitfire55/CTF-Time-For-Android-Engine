package engine

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/net/html"
)

// A Ranking contains the fields to be stored in Firestore for a
type Ranking struct {
	Rank        int
	TeamName    string
	TeamUrl     string
	CountryFlag string
	CountryID   string
	Score       float64
}

// ParseAndStoreRankings parses the response body of a rankings page for the values needed to create a Ranking struct. Once
// the page is parsed and an array of Rankings is created, the results are stored in the Firestore database IF the sha256
// hash of Rankings array is different from the hash stored in the Firestore database. If the page is new, no hash exists
// in the Firestore database, which means a page hash document will be created along with the Ranking documents.
func ParseAndStoreRankings(response *http.Response, pageNumber int, year string, fbc FirebaseContext) error {
	//TODO: Migrate from html.Tokenizer to GoQuery
	var rankings []Ranking
	pageNumDoc := fmt.Sprintf("Page%dHash", pageNumber)



hashCheck:
	if len(rankings) > 0 {
		resultsHash := CalculateHash(rankings)
		hashDiff, err := RankingsHashDiff(resultsHash, pageNumber, year, fbc)
		if err != nil {
			return err
		}
		if hashDiff {
			StoreRankingsHash(resultsHash, pageNumDoc, year, fbc)
			StoreRankings(rankings, year, fbc)
		}
	} else {
		return errors.New("length of rankings array is zero")
	}
	return nil
}
