package arukas

import (
	"fmt"
	"log"
	"os"
	"testing"
)

var realHTTPClient *httpClient

func TestMain(m *testing.M) {

	token := os.Getenv("ARUKAS_JSON_API_TOKEN")
	secret := os.Getenv("ARUKAS_JSON_API_SECRET")

	if token == "" || secret == "" {
		log.Println("[WARN] Please Set ENV 'ARUKAS_JSON_API_TOKEN' and 'ARUKAS_JSON_API_SECRET'")
	} else {

		baseURL := os.Getenv("ARUKAS_JSON_API_URL")
		trace := false
		if os.Getenv("ARUKAS_DEBUG") != "" {
			trace = true
		}

		c, err := NewClient(&ClientParam{
			APIBaseURL: baseURL,
			Token:      token,
			Secret:     secret,
			UserAgent:  fmt.Sprintf("go-arukas-test/v%s", Version),
			Trace:      trace,
			TraceOut:   os.Stderr,
		})

		if err != nil {
			log.Fatal(err)
		}

		realHTTPClient = c.(*client).httpAPI.(*httpClient)
	}
	ret := m.Run()
	os.Exit(ret)
}

func isAccTest() bool {
	return os.Getenv("TEST_ACC") != ""
}
