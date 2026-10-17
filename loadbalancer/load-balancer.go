package loadbalancer

type LoadBalancer struct {
	pool *ServerPool
}

// TODO: get list of server's IPs, make them and put in the list
func NewLB() *LoadBalancer {
	return &LoadBalancer{}
}

func (lb *LoadBalancer) HealthCheck() {
	lb.pool.healthCheckAllServers()
}

// TODO: implement the logic to receive request
func (lb *LoadBalancer) ReceiveRequest() {

}
