package loadbalancer

import "net/url"

type ServerPool struct {
	servers []*Server
}

func newPool(servers []*url.URL) *ServerPool {
	pool := &ServerPool{
		servers: []*Server{},
	}

	for _, server := range servers {
		newServer := newServer(server)
		pool.addServer(newServer)
	}

	return pool
}

func (pool *ServerPool) addServer(server *Server) {
	pool.servers = append(pool.servers, server)
}

func (pool *ServerPool) healthCheckAllServers() {
	for _, server := range pool.servers {
		go server.healthCheckOneServer()
	}
}

func (pool *ServerPool) getServer(hash uint32) *Server {
	numServers := len(pool.servers)
	if numServers == 0 {
		return nil
	}
	serverIndex := int(hash % uint32(numServers))

	for offset := 0; offset < numServers; offset++ {
		currentIndex := (serverIndex + offset) % numServers
		server := pool.servers[currentIndex]
		if server.isAlive() {
			return server
		}
	}
	return nil
}
