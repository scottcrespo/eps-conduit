/*
main.go

Description:
	Eps-Conduit is a light-weight load balancer.

Source Code:
	https://github.com/orlandogolang/eps-conduit
*/
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	LB "github.com/scottcrespo/eps-conduit/load-balancer"
)

// Handling user flags
// User flags must be package globals they can be easily worked on by Config member functions
// and avoid passing each command line option as a parameter.
var configFile = flag.String("config", "/etc/conduit.conf", "Path to config file. Default is /etc/conduit.conf")
var backendStr = flag.String("b", "", "Host strings for the backend services (comma separated)")
var bind = flag.String("bind", "", "The port the load balancer should listen to")
var mode = flag.String("mode", "", "Balancing Mode")
var certFile = flag.String("cert", "", "Path to rsa private key")
var keyFile = flag.String("key", "", "Path to rsa public key")

func main() {
	flag.Parse()

	input := LB.LoadBalancer{
		Bind: *bind,
		Mode: *mode,
		Certfile: *certFile,
		Keyfile: *keyFile,
	}
	input.BackendsFromStr(*backendStr)

	lb := LB.GetLoadBalancer(*configFile, &input)
	// send requests to proxies via lb.handle
	http.HandleFunc("/", lb.Handle)

	// Start the http(s) listener depending on user's selected mode
	if lb.Mode == "http" {
		http.ListenAndServe(":"+lb.Bind, nil)
	} else if lb.Mode == "https" {
		http.ListenAndServeTLS(":"+lb.Bind, lb.Certfile, lb.Keyfile, nil)
	} else {
		fmt.Fprintf(os.Stderr, "unknown mode or mode not set")
		os.Exit(1)
	}
}
