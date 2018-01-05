package engine

import "fmt"

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

func StoreTeam(fbc FirebaseContext, team Team) {

}
