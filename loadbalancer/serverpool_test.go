package loadbalancer

import "testing"

func TestNewPoolCreatesAllServers(t *testing.T) {
	pool := newPool(fakeURLs(3))

	if len(pool.servers) != 3 {
		t.Fatalf("expected 3 servers, got %d", len(pool.servers))
	}
	for i, s := range pool.servers {
		if !s.isAlive() {
			t.Errorf("server %d should start as alive", i)
		}
	}
}

func TestNewPoolEmpty(t *testing.T) {
	pool := newPool(nil)
	if len(pool.servers) != 0 {
		t.Fatalf("expected empty pool, got %d servers", len(pool.servers))
	}
}

func TestGetServerEmptyPoolReturnsNil(t *testing.T) {
	pool := newPool(nil)
	if got := pool.getServer(123); got != nil {
		t.Fatalf("expected nil from empty pool, got %v", got.url)
	}
}

func TestGetServerDeterministicForSameHash(t *testing.T) {
	pool := newPool(fakeURLs(3))

	for _, hash := range []uint32{0, 1, 7, 42, 1_000_000} {
		first := pool.getServer(hash)
		second := pool.getServer(hash)
		if first == nil || second == nil {
			t.Fatalf("hash %d: got nil server", hash)
		}
		if first != second {
			t.Errorf("hash %d: expected same server on repeated calls, got %s then %s",
				hash, first.url, second.url)
		}
	}
}

func TestGetServerSkipsDeadServer(t *testing.T) {
	pool := newPool(fakeURLs(3))
	dead := pool.servers[0]
	dead.setAlive(false)

	got := pool.getServer(0)
	if got == nil {
		t.Fatal("expected fallback to an alive server, got nil")
	}
	if got == dead {
		t.Fatalf("expected alive server, got the dead one (%s)", got.url)
	}
}

func TestGetServerWalksForwardInOrder(t *testing.T) {
	pool := newPool(fakeURLs(3))
	pool.servers[0].setAlive(false)
	pool.servers[1].setAlive(false)

	got := pool.getServer(0)
	if got != pool.servers[2] {
		t.Fatalf("expected server index 2, got %v", got.url)
	}
}

func TestGetServerWrapsAroundProbe(t *testing.T) {
	pool := newPool(fakeURLs(3))
	pool.servers[2].setAlive(false)

	got := pool.getServer(2)
	if got != pool.servers[0] {
		t.Fatalf("expected probe to wrap around to server index 0, got %v", got.url)
	}
}

func TestGetServerMapsHashModuloPoolSize(t *testing.T) {
	pool := newPool(fakeURLs(3))

	cases := []struct {
		hash  uint32
		index int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 0},
		{4, 1},
		{1_000_000, int(1_000_000 % 3)},
	}
	for _, c := range cases {
		got := pool.getServer(c.hash)
		if got != pool.servers[c.index] {
			t.Errorf("hash %d: expected server index %d, got %v", c.hash, c.index, got.url)
		}
	}
}

func TestGetServerAllDeadReturnsNil(t *testing.T) {
	pool := newPool(fakeURLs(3))
	for _, s := range pool.servers {
		s.setAlive(false)
	}
	if got := pool.getServer(5); got != nil {
		t.Fatalf("expected nil when all servers are dead, got %v", got.url)
	}
}
