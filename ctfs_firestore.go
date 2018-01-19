// Copyright 2017-2018 Dale Lakes <spitfire@spitfy.re>. All rights reserved.
// Use of this source code is governed by the MIT license located in the LICENSE file.

package goctftime

import (
	"fmt"
	"strconv"
	"strings"
)

func GetLastCtfId(fbc FirebaseContext) int {
	lastPageNumberDoc, err := fbc.Fb.Collection("CTFs").Doc("LastCtfId").Get(fbc.Ctx)
	if err != nil {
		return 0
	}
	lastPageNumber, err := lastPageNumberDoc.DataAt("lastCtfId")
	if err != nil {
		return 0
	}
	if lastPageNumberInt, ok := lastPageNumber.(int64); ok {
		return int(lastPageNumberInt)
	} else {
		return 0
	}
}

func UpdateLastCtfId(fbc FirebaseContext, newCtfId int) {
	_, err := fbc.Fb.Collection("CTFs").Doc("LastCtfId").Set(fbc.Ctx, map[string]int{
		"lastCtfId": newCtfId,
	})
	if err != nil {
		fmt.Println(err)
	}
}

func StoreCtf(ctfId int, ctf Ctf, fbc FirebaseContext) error {
	_, err := fbc.Fb.Collection("CTFs").Doc(strconv.Itoa(ctfId)).Set(fbc.Ctx, ctf)
	if err != nil {
		return err
	}
	return nil
}

func CompareCtfHash(id int, ctf Ctf, fbc FirebaseContext) (bool, error) {
	hashDoc, err := fbc.Fb.Collection("CTFs").Doc(strconv.Itoa(id)).Get(fbc.Ctx)
	if err != nil {
		// Ctf not found, return true to create it
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
	if ctf.Hash != hashDocValue {
		// Hashes are different, so return true
		return true, nil
	} else {
		// Hashes are same, so return false to prevent unnecessary write
		return false, nil
	}
}
