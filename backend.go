package main

import (
	"net/http/httputil"
	"net/url"
	"sync"
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

type ServerPool struct {
	servers []*Server
}

func (pool *ServerPool) addServer(server *Server) {
	pool.servers = append(pool.servers, server)
}
