package main

import (
	"log"
	"net/http/httputil"
	"net/url"
)

func main() {
	//TODO: Add more serves and put them in the pool
	serverUrl, err := url.Parse("http://localhost:8003")
	if err != nil {
		log.Fatal(err)
	}

	server := &Server{
		url:          serverUrl,
		alive:        true,
		reverseProxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}

	pool := &ServerPool{
		servers: []*Server{server},
	}

	pool.addServer(server)

	pool.healthCheckAllServers()
}
