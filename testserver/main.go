package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	var (
		port = flag.Int("port", 8003, "Port to listen on")
		name = flag.String("name", "", "Server name reported in responses (defaults to backend-<port>)")
	)
	flag.Parse()

	id := *name
	if id == "" {
		id = fmt.Sprintf("backend-%d", *port)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "response from %s (path: %s)\n", id, r.URL.Path)
	})

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	log.Printf("%s listening on %s", id, addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
