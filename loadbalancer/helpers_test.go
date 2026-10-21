package loadbalancer

import (
	"fmt"
	"net/url"
)

func fakeURLs(n int) []*url.URL {
	urls := make([]*url.URL, n)
	for i := range urls {
		urls[i], _ = url.Parse(fmt.Sprintf("http://127.0.0.1:%d", 9000+i))
	}
	return urls
}
