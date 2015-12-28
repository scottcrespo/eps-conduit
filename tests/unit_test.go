package tests

import (
	"io/ioutil"
	"log"
	"net/http"
  "net/http/httptest"
	"testing"

  LB "github.com/scottcrespo/eps-conduit/load-balancer"
)

func TestStandalone(t *testing.T) {

  configFile := "/etc/conduit.conf"
  // initialize the main LoadBalancer Instance using configFile and empty user input
	lb := LB.GetLoadBalancer(configFile, &LB.LoadBalancer{})

  ts := httptest.NewServer(http.HandlerFunc(lb.Handle))
  defer ts.Close()
  //ts.URL = "http://localhost:8000"

  for i := 0; i < 5; i++ {
    resp, err := http.Get(ts.URL)
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
