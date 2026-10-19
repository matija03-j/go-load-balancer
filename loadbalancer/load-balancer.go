package loadbalancer

import (
	"hash/fnv"
	"log"
	"net"
	"net/http"
)

const (
	LOAD_BALANCER_ADDRESS = "localhost:8000"
)

type LoadBalancer struct {
	pool *ServerPool
}

// TODO: get list of server's IPs, make them and put in the list
func NewLB() *LoadBalancer {
	return &LoadBalancer{}
}

func (lb *LoadBalancer) Run() error {
	go lb.HealthCheck()
	return lb.ListenAndServe(LOAD_BALANCER_ADDRESS)
}

func (lb *LoadBalancer) HealthCheck() {
	lb.pool.healthCheckAllServers()
}

func (lb *LoadBalancer) ListenAndServe(address string) error {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sourceIP, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Printf("Cannot extract client IP from %q: %v", r.RemoteAddr, err)
			http.Error(
				w,
				"Cannot determine client IP",
				http.StatusInternalServerError,
			)
			return
		}

		lb.processRequest(w, r, sourceIP)
	})

	log.Printf("Load balancer listening on %s", address)
	return http.ListenAndServe(address, handler)
}

func (lb *LoadBalancer) processRequest(w http.ResponseWriter, r *http.Request, sourceIP string) {
	hashedIP := hashIP(sourceIP)
	server := lb.pool.getServer(hashedIP)
	if server == nil {
		http.Error(
			w,
			"No available servers",
			http.StatusServiceUnavailable,
		)
		return
	}

	sendRequestToServer(w, r, server)
}

func hashIP(ip string) uint32 {
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(ip))
	return hasher.Sum32()
}

func sendRequestToServer(w http.ResponseWriter, r *http.Request, server *Server) {
	server.reverseProxy.ServeHTTP(w, r)
}
