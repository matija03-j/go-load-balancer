package loadbalancer

type ServerPool struct {
	servers []*Server
}

func (pool *ServerPool) addServer(server *Server) {
	pool.servers = append(pool.servers, server)
}

func (pool *ServerPool) healthCheckAllServers() {
	for _, server := range pool.servers {
		go server.healthCheckOneServer()
	}
}
