package main

import (
	"log"
	"net/http"
	"net/url"
	"time"
)

func (pool *ServerPool) healthCheckAllServers() {
	for _, server := range pool.servers {
		go server.healthCheckOneServer()
	}
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

func load