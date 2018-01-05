package engine

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
)

// Fetch executes a GET request to a URL. If the response status code is not 200, an error is returned specifying which URL
// failed and the status code of the request.
func Fetch(url string) (*http.Response, error) {
	client := http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		errorVal := fmt.Sprintf("Web page status for %s is %d", url, resp.StatusCode)
		return nil, errors.New(errorVal)
	}
	return resp, nil
}

func CalculateHash(data interface{}) string {
	sha256Hash := sha256.New()
	sha256Hash.Write([]byte(fmt.Sprintf("%#v", data)))
	return base64.StdEncoding.EncodeToString(sha256Hash.Sum(nil))
}
