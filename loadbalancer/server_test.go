package loadbalancer

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewServerInitialAliveState(t *testing.T) {
	url := fakeURLs(1)[0]
	server := newServer(url)

	if !server.isAlive() {
		t.Errorf("expected server to start as alive, got dead")
	}
}

func TestNewServerInitialFailureCount(t *testing.T) {
	url := fakeURLs(1)[0]
	server := newServer(url)

	if server.failureCount != 0 {
		t.Errorf("expected server to start with 0 failures, got %d", server.failureCount)
	}
}

func TestServerInitialProxy(t *testing.T) {
	url := fakeURLs(1)[0]
	server := newServer(url)

	if server.reverseProxy == nil {
		t.Errorf("expected server to have a reverse proxy, got nil")
	}
}

func TestServerSetAlive(t *testing.T) {
	url := fakeURLs(1)[0]
	server := newServer(url)

	server.setAlive(false)
	if server.isAlive() {
		t.Errorf("expected server to be dead after setAlive(false), got alive")
	}

	server.setAlive(true)
	if !server.isAlive() {
		t.Errorf("expected server to be alive after setAlive(true), got dead")
	}
}

func TestServerIncrementFailureCount(t *testing.T) {
	url := fakeURLs(1)[0]
	server := newServer(url)

	for i := 0; i < MAX_FAILURES-1; i++ {
		server.incrementFailureCount()
		if server.isAlive() == false {
			t.Errorf("expected server to be alive after %d failures, got dead", i+1)
		}
	}

	server.incrementFailureCount()
	if server.isAlive() {
		t.Errorf("expected server to be dead after %d failures, got alive", MAX_FAILURES)
	}
	if server.failureCount != 0 {
		t.Errorf("expected failure count to reset to 0 after reaching max failures, got %d", server.failureCount)
	}
}

func TestResetFailureCount(t *testing.T) {
	url := fakeURLs(1)[0]
	server := newServer(url)
	server.incrementFailureCount()
	server.incrementFailureCount()

	server.resetFailureCount()
	if server.failureCount != 0 {
		t.Errorf("expected failure count to reset to 0 after restart, got %d", server.failureCount)
	}
}

func TestIsServerAliveHealthy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	if !isServerAlive(u) {
		t.Errorf("expected server to be alive, got dead")
	}
}

func TestIsServerAliveUnreachable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Close()

	u, _ := url.Parse(server.URL)
	if isServerAlive(u) {
		t.Errorf("expected server to be dead, got alive")
	}
}
