package tests

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestEndToEnd(t *testing.T) {

	for i := 0; i <= 4; i++ {
		req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
		if err != nil {
			t.Errorf("Failed to create new request")
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("client failed to issue request")
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("failed to read response body")
		}
		log.Println(string(body))
	}
}
