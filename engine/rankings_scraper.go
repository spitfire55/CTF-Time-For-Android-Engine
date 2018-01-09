package engine

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// A Ranking represents a row in the ctftime.org Rankings
type Ranking struct {
	CountryId string
	Events    int
	Rank      int
	Score     float64
	TeamId    int
	TeamName  string
}

// ParseAndStoreRankings parses the response body of a rankings page for the values needed to create a Ranking struct. Once
// the page is parsed and an array of Rankings is created, the results are stored in the Firestore database IF the sha256
// hash of the Rankings array is different from the hash stored in the Firestore database for the particular page. If the
// final page parsed is a newly created page, no page hash document exists in the Firestore database. We then create a page
// hash document along with the Ranking documents.
// NOTE: The page hash is computed from an array of Rankings (i.e. after parsing), not from the response body of the request.
func ParseAndStoreRankings(response *http.Response, pageNumber int, year string, fbc FirebaseContext) error {
	var rankings []Ranking
	pageNumDoc := fmt.Sprintf("Page%dHash", pageNumber)
	rootSel, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return err
	}

	rootSel.Find(".table.table-striped tr").Each(func(rowIndex int, rowSelection *goquery.Selection) {
		var ranking Ranking
		if rowIndex == 0 {
			return
		}
		rowSelection.Find("td").Each(func(colIndex int, colSelection *goquery.Selection) {
			switch colIndex {
			case 0:
				ranking.Rank, _ = strconv.Atoi(colSelection.Text())
			case 1:
				ranking.TeamName = colSelection.Find("a").Text()
				teamIdUrl, _ := colSelection.Find("a").Attr("href")
				teamIdSplit := strings.Split(teamIdUrl, "/")
				ranking.TeamId, _ = strconv.Atoi(teamIdSplit[len(teamIdSplit)-1])
			case 2:
				ranking.CountryId, _ = colSelection.Find("img").Attr("alt")
			case 3:
				ranking.Score, _ = strconv.ParseFloat(colSelection.Text(), 64)
			case 4:
				ranking.Events, _ = strconv.Atoi(colSelection.Text())
			}
		})
		rankings = append(rankings, ranking)
	})
	fmt.Printf("%#v\n", rankings)

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
