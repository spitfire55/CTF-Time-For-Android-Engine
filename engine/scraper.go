package engine

import (
	"errors"
	"net/http"
)

func Fetch(url string, fbc FirebaseContext) (*http.Response, error) {
	resp, err := fbc.C.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Web page status is: " + resp.Status)
	}
}
