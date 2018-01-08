package engine

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

type Ctf struct { // key = CTF ID
	hash  string
	Name  string
	Url   string
	Image string // relative Url to image
}

func ParseAndStoreCtf(ctfId int, resp *http.Response, fbc FirebaseContext) error {
	var ctf Ctf
	rootSel, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}
	ctf.Name = rootSel.Find("h2").Text()
	urlText := rootSel.Find(".row .span10 p").Text()
	if urlText != "" {
		ctf.Url = strings.Join(strings.Split(urlText, " ")[1:], " ")
	}
	ctf.Image, _ = rootSel.Find(".span2 img").Attr("src")

	ctfHash := CalculateHash(ctf)
	ctf.hash = ctfHash
	hashDiff, err := CtfHashDiff(ctfId, ctf, fbc)
	if err != nil {
		fmt.Println("booty")
		return err
	}
	if hashDiff {
		err := StoreCtf(ctfId, ctf, fbc)
		if err != nil {
			fmt.Println("ballz")
			return err
		}
	}
	return nil
}
