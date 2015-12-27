package tests

import (
  "net/http"
  "testing"
  "io/ioutil"
  "log"
)

func TestListening(t *testing.T) {

  for i:=0; i <= 4; i = i+1 {
    req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
    if err != nil{
      t.Errorf("Failed to create new request")
    }
    //req.Host = "http://localhost:8000"

    client := &http.Client{}
  //  resp, err := http.Get("http://localhost:8000")
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
