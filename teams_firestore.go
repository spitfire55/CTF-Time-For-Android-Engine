// Copyright 2017-2018 Dale Lakes <spitfire@spitfy.re>. All rights reserved.
// Use of this source code is governed by the MIT license located in the LICENSE file.

package goctftime

import (
	"fmt"
	"strconv"
	"strings"
)

func GetLastTeamId(fbc FirebaseContext) int {
	lastPageNumberDoc, _ := fbc.Fb.Collection("Teams").Doc("LastTeamId").Get(fbc.Ctx)
	lastPageNumber, _ := lastPageNumberDoc.DataAt("lastTeamId")
	if lastPageNumberInt, ok := lastPageNumber.(int64); ok {
		return int(lastPageNumberInt)
	} else {
		return 0
	}
}

func UpdateLastTeamId(fbc FirebaseContext, newPageNumber int) {
	_, err := fbc.Fb.Collection("Teams").Doc("LastTeamId").Set(fbc.Ctx, map[string]int{
		"lastTeamId": newPageNumber,
	})
	if err != nil {
		fmt.Println(err)
	}
}

func StoreTeam(teamId int, team Team, fbc FirebaseContext) error {
	_, err := fbc.Fb.Collection("Teams").Doc(strconv.Itoa(teamId)).Set(fbc.Ctx, team)
	if err != nil {
		return err
	}
	return nil
}

func CompareTeamHash(id int, team Team, fbc FirebaseContext) (bool, error) {
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
