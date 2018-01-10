// Copyright 2017-2018 Dale Lakes <spitfire@spitfy.re>. All rights reserved.
// Use of this source code is governed by the MIT license located in the LICENSE file.

package goctftime

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

type Ctf struct { // key = CTF ID
	Hash  string
	Image string // relative Url to image
	Name  string
	Url   string
}

func ParseAndStoreCtf(ctfId int, resp *http.Response, fbc FirebaseContext) error {
	var ctf Ctf
	rootSel, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}
	ctf.Name = rootSel.Find("h2").Text()
	urlText := strings.Split(rootSel.Find(".row .span10 p").First().Text(), " ")
	if len(urlText) != 0 {
		ctf.Url = urlText[len(urlText)-1]
	}
	ctf.Image, _ = rootSel.Find(".span2 img").Attr("src")

	ctfHash := CalculateHash(ctf)
	ctf.Hash = ctfHash
	hashDiff, err := CtfHashDiff(ctfId, ctf, fbc)
	if err != nil {
		return err
	}
	if hashDiff {
		err := StoreCtf(ctfId, ctf, fbc)
		if err != nil {
			return err
		}
	}
	return nil
}
