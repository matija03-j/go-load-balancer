package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

func main() {
	var (
		target   = flag.String("target", "http://localhost:8000", "Load balancer address")
		requests = flag.Int("n", 100, "Total number of requests to send")
		ips      = flag.String("ips", "", "Comma-separated source IPs to bind to (e.g. 127.0.0.1,127.0.0.2,127.0.0.3)")
	)
	flag.Parse()

	clients := buildClients(*ips)

	results := make([]userResult, len(clients))

	var wg sync.WaitGroup
	for i, client := range clients {
		wg.Add(1)
		go func(i int, client *ipClient) {
			defer wg.Done()
			counts, failures := runUser(client, *target, *requests)
			results[i] = userResult{counts: counts, failures: failures}
		}(i, client)
	}
	wg.Wait()

	fmt.Printf("Sent %d requests per user to %s (%d users)\n", *requests, *target, len(clients))
	for i, result := range results {
		keys := make([]string, 0, len(result.counts))
		for k := range result.counts {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		fmt.Printf("%s: %d failed\n", clients[i].ip, result.failures)
		for _, k := range keys {
			fmt.Printf("  %6d  %s\n", result.counts[k], k)
		}
	}
	log.Println("done")
}

type ipClient struct {
	ip   string
	http *http.Client
}

type userResult struct {
	counts   map[string]int
	failures int
}

func buildClients(ips string) []*ipClient {
	if ips == "" {
		return []*ipClient{{ip: "(default)", http: &http.Client{Timeout: 5 * time.Second}}}
	}

	clients := []*ipClient{}
	for _, ip := range strings.Split(ips, ",") {
		ip = strings.TrimSpace(ip)
		parsed := net.ParseIP(ip)
		if parsed == nil {
			log.Fatalf("Invalid source IP: %q", ip)
		}

		transport := &http.Transport{
			DialContext: (&net.Dialer{
				LocalAddr: &net.TCPAddr{IP: parsed},
				Timeout:   5 * time.Second,
			}).DialContext,
		}
		clients = append(clients, &ipClient{
			ip:   ip,
			http: &http.Client{Transport: transport, Timeout: 5 * time.Second},
		})
	}
	return clients
}

func runUser(client *ipClient, target string, n int) (map[string]int, int) {
	counts := map[string]int{}
	failures := 0

	for i := 0; i < n; i++ {
		time.Sleep(1 * time.Second)
		resp, err := client.http.Get(target)
		if err != nil {
			failures++
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		key := strings.TrimSpace(string(body))
		if resp.StatusCode != http.StatusOK {
			key = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, key)
		}

		counts[key]++
	}

	return counts, failures
}
