package main

import (
	lb "loadbalancer/loadbalancer"
)

func main() {
	//TODO: Add more serves and put them in the pool

	lb := lb.NewLB()
	go lb.HealthCheck()
	lb.ReceiveRequest()
}
