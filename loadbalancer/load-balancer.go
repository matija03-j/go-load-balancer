package loadbalancer

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

// TODO: implement the logic to receive request
func (lb *LoadBalancer) ListenAndServe(address string) error {
	return nil
}
