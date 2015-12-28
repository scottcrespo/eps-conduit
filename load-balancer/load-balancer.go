package loadbalancer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config Struct represents the load balancer's configuration
type LoadBalancer struct {

	// the backend services to balance
	Backends []string `toml:"backends"`

	// The port the load balancer is bound to
	Bind string `toml:"bind"`

	// Secure or unsecure http protocol
	Mode string `toml:"mode"`

	// Path to certificate file
	Certfile string `toml:certFile`

	// Path to private key file related to certificate
	Keyfile string `toml:keyFile`

	// Revserse Proxies to forward requests to
	Proxies []*httputil.ReverseProxy

	// Number of proxies available
	HostCount int

	// The index of the next proxy to forward a request to
	NextHost int
}

// singleton Config instance initially set to nil
var lb *LoadBalancer = nil

// GetLoadBalancer implements a singleton pattern to access the LoadBalancer singleton
func GetLoadBalancer(configFile string, input *LoadBalancer) *LoadBalancer {
	if lb == nil {
		lb = new(LoadBalancer)
		lb.init(configFile, input)
	}
	return lb
}

// init initializes a new Config instance by reading from the config file
// It will unmarshal the toml file into the Config struct
func (lb *LoadBalancer) init(configFile string, input *LoadBalancer) {
	_, err := os.Stat(configFile)

	if err != nil {
		log.Fatal("Config file not found: ", configFile)
	}
	if _, err := toml.DecodeFile(configFile, lb); err != nil {
		log.Fatal(err)
	}
	lb.handleUserInput(input)
	lb.printConfigInfo()
	lb.makeProxies()
	lb.HostCount = len(lb.Backends)
	lb.NextHost = 0
}

// handleUserInput checks command line input and overrides config file settings
// Backends is parsed from a raw string to a slice of strings
// TODO: Better input validation
func (lb *LoadBalancer) handleUserInput(input *LoadBalancer) {
	if len(input.Backends) > 0 {
		lb.Backends = input.Backends
	}
	if input.Bind != "" {
		lb.Bind = input.Bind
	}
	if input.Mode != "" {
		lb.Mode = input.Mode
	}
	if input.Certfile != "" {
		lb.Certfile = input.Certfile
	}
	if input.Keyfile != "" {
		lb.Keyfile = input.Keyfile
	}
}

// BackendsFromStr creates a slice of backends from a user-inputted string, and copies the result
// into the LoadBalancer singleton instance.
func (lb *LoadBalancer) BackendsFromStr(backendStr string) {
	if backendStr != "" {
		// Remove whitespace from backends
		backendStr = strings.Replace(backendStr, " ", "", -1)
		// Throwing backends into an array
		lb.Backends = strings.Split(backendStr, ",")
	}
}

// printConfigInfo prints to stderr host and port settings applied to current process
func (lb *LoadBalancer) printConfigInfo() {
	// tell the user what info the load balancer is using
	for _, v := range lb.Backends {
		log.Println("using " + v + " as a backend")
	}
	log.Println("listening on port " + lb.Bind)
}

// makeProxies creates slice of ReverseProxies based on the LoadBalancer backend hosts
// It returns a slice of httputil.ReverseProxy
func (lb *LoadBalancer) makeProxies() {
	// Create a proxy for each backend
	lb.Proxies = make([]*httputil.ReverseProxy, len(lb.Backends))
	for i := range lb.Backends {
		// host must be defined here, and not within the anonymous function.
		// Otherwise, you'll run into scoping issues
		host := lb.Backends[i]
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = host
		}
		lb.Proxies[i] = &httputil.ReverseProxy{Director: director}
	}
}

// handle forwards request by calling ServeHTTP() on the next Proxy
func (lb *LoadBalancer) Handle(w http.ResponseWriter, r *http.Request) {
	lb.pickHost()
	lb.Proxies[lb.NextHost].ServeHTTP(w, r)
}

// pickHost determines the next backend host to forward the request to - according to round-robin
// It returns an integer, which represents the host's index in lb.Backends
func (lb *LoadBalancer) pickHost() {
	nextHost := lb.NextHost + 1
	if nextHost >= lb.HostCount {
		lb.NextHost = 0
	} else {
		lb.NextHost = nextHost
	}
}