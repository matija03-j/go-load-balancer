package loadbalancer

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHashIPKnownValues(t *testing.T) {
	expected := map[string]uint32{
		"127.0.0.1":     144953630,
		"10.0.0.5":      3669932097,
		"192.168.1.100": 536723748,
		"8.8.8.8":       672416809,
	}
	for ip, want := range expected {
		if got := hashIP(ip); got != want {
			t.Errorf("hashIP(%q) = %d, want %d", ip, got, want)
		}
	}
}

func TestProcessRequestNoServersReturns503(t *testing.T) {
	lb := &LoadBalancer{pool: newPool(nil)}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	lb.processRequest(rec, req, "10.0.0.1")

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", rec.Code)
	}
	if body := rec.Body.String(); body != "No available servers\n" {
		t.Errorf("unexpected body: %q", body)
	}
}

func TestProcessRequestProxiesToBackend(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from backend"))
	}))
	defer backend.Close()

	u, _ := url.Parse(backend.URL)
	lb := &LoadBalancer{pool: newPool([]*url.URL{u})}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	lb.processRequest(rec, req, "10.0.0.1")

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "hello from backend" {
		t.Errorf("unexpected body: %q", rec.Body.String())
	}
}
