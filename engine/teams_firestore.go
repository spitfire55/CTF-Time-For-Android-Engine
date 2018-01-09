package engine

import (
	"fmt"
	"strconv"
	"strings"
)

func GetLastTeamId(fbc FirebaseContext) int {
	lastPageNumberDoc, _ := fbc.Fb.Collection("Teams").Doc("LastTeamId").Get(fbc.Ctx)
	lastPageNumber, _ := lastPageNumberDoc.DataAt("lastPageNumber")
	if lastPageNumberInt, ok := lastPageNumber.(int64); ok {
		return int(lastPageNumberInt)
	} else {
		return 0
	}
}

func UpdateLastTeamId(fbc FirebaseContext, newPageNumber int) {
	_, err := fbc.Fb.Collection("Teams").Doc("LastPageNumber").Set(fbc.Ctx, map[string]int{
		"lastPageNumber": newPageNumber,
	})
	if err != nil {
		fmt.Println(err)
	}
}

// StoreTeam stores a Team document to Firestore database. A new document is created or an existing document is updated.
// The Document key is the team ID. The upload path for a Team struct is: Teams/<team_id>
func StoreTeam(teamId int, team Team, fbc FirebaseContext) error {
	_, err := fbc.Fb.Collection("Teams").Doc(strconv.Itoa(teamId)).Set(fbc.Ctx, team)
	if err != nil {
		return err
	}
	return nil
}

// TeamHashDiff compares the hash computed from a team page on ctftime.org to the hash stored of that page in the
// Firestore database. If the hashes are not equal, then that team page has changed in some way. This check is used to
// prevent unnecessary writes to the Firestore database if a team page has not changed between scrape iterations.
// If the team is new, a new document is created to store the hash for that page.
// NOTE: The hash is computed from the Team struct (after parsing), not from the response body of the request.
func TeamHashDiff(id int, team Team, fbc FirebaseContext) (bool, error) {
	hashDoc, err := fbc.Fb.Collection("Teams").Doc(strconv.Itoa(id)).Get(fbc.Ctx)
	if err != nil {
		// Team not found, return true to create it
		if strings.Contains(err.Error(), "NotFound") {
			return true, nil
		}
		// Some other error, so return error
		return false, err
	}
	hashDocValue, err := hashDoc.DataAt("Hash")
	if err != nil {
		// Document doesn't have hash field or we can't read it, so return error
		return false, err
	}
	if team.Hash != hashDocValue {
		// Hashes are different, so return true
		return true, nil
	} else {
		// Hashes are same, so return false to prevent unnecessary write
		return false, nil
	}
}
