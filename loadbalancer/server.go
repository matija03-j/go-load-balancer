package loadbalancer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type Server struct {
	url          *url.URL
	alive        bool
	mux          sync.RWMutex
	reverseProxy *httputil.ReverseProxy
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
		time.Sleep(5 * time.Second)
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
