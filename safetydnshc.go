package main

import (
	"github.com/miekg/dns"
	"log"
	"sync"
	"time"
)

var (
	upstreamDNS        = "8.8.8.8:53" // Example: Google's Public DNS
	isUpstreamAvailable bool         // Global flag to indicate upstream DNS status
	cache              *DNSCache
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
		StalePeriod: now.Add(24 * time.Hour), // Using a fixed 24-hour period for stale data
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
	if now.Before(entry.Expires) || (now.Before(entry.StalePeriod) && !isUpstreamAvailable) {
		return entry.Response, true
	}
	return nil, false
}

func handleDNSQuery(w dns.ResponseWriter, r *dns.Msg) {
	key := r.Question[0].Name + dns.TypeToString[r.Question[0].Qtype]
	if response, found := cache.Get(key); found {
		response.SetReply(r)
		w.WriteMsg(response)
		return
	}

	if !isUpstreamAvailable {
		dns.HandleFailed(w, r)
		return
	}

	// Forward the query to the upstream DNS server
	c := new(dns.Client)
	response, _, err := c.Exchange(r, upstreamDNS)
	if err != nil {
		log.Printf("Failed to reach upstream DNS: %v", err)
		isUpstreamAvailable = false
		dns.HandleFailed(w, r)
		return
	}

	cache.Set(key, response, 5*time.Minute) // Cache with a TTL of 5 minutes
	w.WriteMsg(response)
}

func checkUpstreamDNS() {
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn("example.com"), dns.TypeA)
	_, _, err := c.Exchange(m, upstreamDNS)
	if err != nil {
		log.Println("Upstream DNS check failed:", err)
		isUpstreamAvailable = false
	} else {
		log.Println("Upstream DNS is available")
		isUpstreamAvailable = true
	}
}

func scheduleHealthCheck(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				checkUpstreamDNS()
			}
		}
	}()
}

func main() {
	// Initialize the DNS cache
	cache = NewDNSCache()

	// Schedule upstream DNS health checks every 5 minutes
	scheduleHealthCheck(5 * time.Minute)

	// Setup DNS server
	dns.HandleFunc(".", handleDNSQuery)
	server := &dns.Server{Addr: "127.0.0.1:53", Net: "udp"}

	log.Println("Starting DNS proxy server...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
