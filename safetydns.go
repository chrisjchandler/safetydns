package main

import (
	"github.com/miekg/dns"
	"sync"
	"time"
)

type CacheEntry struct {
	Response    *dns.Msg
	Expires     time.Time
	StalePeriod time.Time
}

type DNSCache struct {
	cache map[string]*CacheEntry
	mu    sync.RWMutex
}

func NewDNSCache() *DNSCache {
	return &DNSCache{
		cache: make(map[string]*CacheEntry),
	}
}

func (c *DNSCache) Set(key string, response *dns.Msg, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	c.cache[key] = &CacheEntry{
		Response:    response,
		Expires:     now.Add(ttl),
		StalePeriod: now.Add(24 * time.Hour), // Using a fixed 24 hours for stale period
	}
}

func (c *DNSCache) Get(key string) (*dns.Msg, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, found := c.cache[key]
	if !found {
		return nil, false
	}
	now := time.Now()
	if now.Before(entry.Expires) || now.Before(entry.StalePeriod) {
		return entry.Response, true
	}
	return nil, false
}

var cache = NewDNSCache()

func handleDNSQuery(w dns.ResponseWriter, r *dns.Msg) {
	key := r.Question[0].Name + dns.TypeToString[r.Question[0].Qtype]
	cachedResponse, found := cache.Get(key)
	if found {
		cachedResponse.SetReply(r)
		w.WriteMsg(cachedResponse)
		return
	}

	// Forward the query to upstream DNS (example: 8.8.8.8 Google DNS)
	c := new(dns.Client)
	in, _, err := c.Exchange(r, "8.8.8.8:53")
	if err != nil {
		// Handle error, possibly by serving stale data if available
		dns.HandleFailed(w, r)
		return
	}

	// Cache the successful response with a TTL (for simplicity, TTL is assumed)
	cache.Set(key, in, 5*time.Minute) // Example TTL of 5 minutes
	w.WriteMsg(in)
}

func main() {
	// Setup DNS server
	dns.HandleFunc(".", handleDNSQuery)
	server := &dns.Server{Addr: "127.0.0.1:53", Net: "udp"}
	
	// Start DNS server
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
