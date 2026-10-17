package main

import (
	lb "loadbalancer/loadbalancer"
	"log"
)

func main() {
	//TODO: Add more serves and put them in the pool

	lb := lb.NewLB()
	err := lb.Run()
	if err != nil {
		log.Fatal(err)
	}
}
