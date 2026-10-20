package main

import (
	"flag"
	lb "loadbalancer/loadbalancer"
	"log"
	"net/url"
	"strings"
)

func main() {
	var servers string
	flag.StringVar(&servers, "servers", "", "Load balanced backends, use commas to separate")
	flag.Parse()

	if len(servers) == 0 {
		log.Fatal("Please provide one or more backends to load balance")
	}

	serverList := []*url.URL{}
	tokens := strings.Split(servers, ",")
	for _, tok := range tokens {
		serverUrl, err := url.Parse(tok)
		if err != nil {
			log.Fatal(err)
		}
		serverList = append(serverList, serverUrl)
	}

	lb := lb.NewLB(serverList)
	err := lb.Run()
	if err != nil {
		log.Fatal(err)
	}
}
