package loadbalancer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

const (
	MAX_FAILURES = 3
)

type Server struct {
	url          *url.URL
	alive        bool
	mux          sync.RWMutex
	reverseProxy *httputil.ReverseProxy
	failureCount int
}

func newServer(url *url.URL) *Server {
	server := &Server{
		url:          url,
		alive:        true,
		failureCount: 0,
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("[%s] %s\n", url.Host, err.Error())
		server.incrementFailureCount()
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
	}
	proxy.ModifyResponse = func(response *http.Response) error {
		server.resetFailureCount()
		return nil
	}

	server.reverseProxy = proxy
	return server
}

func (server *Server) isAlive() bool {
	server.mux.RLock()
	defer server.mux.RUnlock()

	return server.alive
}

func (server *Server) setAlive(alive bool) {
	server.mux.Lock()
	defer server.mux.Unlock()

	server.alive = alive
}

func (server *Server) healthCheckOneServer() {
	for {
		if isServerAlive(server.url) {
			server.setAlive(true)
			log.Printf("%s is alive\n", server.url)
		} else {
			server.setAlive(false)
			log.Printf("%s is dead\n", server.url)
		}
		time.Sleep(20 * time.Second)
	}
}

func isServerAlive(url *url.URL) bool {
	healthURL := url.String() + "/health"

	client := http.Client{Timeout: time.Second}

	resp, err := client.Get(healthURL)
	if err != nil {
		return false
	}

	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func (server *Server) incrementFailureCount() {
	server.mux.Lock()
	defer server.mux.Unlock()

	server.failureCount++
	if server.failureCount >= MAX_FAILURES {
		server.alive = false
		server.failureCount = 0
	}
}

func (server *Server) resetFailureCount() {
	server.mux.Lock()
	defer server.mux.Unlock()

	server.failureCount = 0
}
