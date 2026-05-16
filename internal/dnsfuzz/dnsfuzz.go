package dnsfuzz

import (
	"net"
	"sync"
)

type FuzzResult struct {
	Host string   `json:"host"`
	IPs  []string `json:"ips"`
}

var defaultWordlist = []string{"www", "mail", "ftp", "localhost", "webmail", "smtp", "pop", "ns1", "web", "ns2", "dns", "ns", "test", "dev", "vpn", "m", "admin", "api", "blog", "secure"}

func BruteForce(domain string, wordlist []string, concurrency int) ([]FuzzResult, error) {
	if len(wordlist) == 0 {
		wordlist = defaultWordlist
	}

	var results []FuzzResult
	var mu sync.Mutex

	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for _, word := range wordlist {
		wg.Add(1)
		sem <- struct{}{}
		go func(sub string) {
			defer wg.Done()
			defer func() { <-sem }()
			
			fqdn := sub + "." + domain
			ips, err := net.LookupHost(fqdn)
			if err == nil && len(ips) > 0 {
				mu.Lock()
				results = append(results, FuzzResult{Host: fqdn, IPs: ips})
				mu.Unlock()
			}
		}(word)
	}

	wg.Wait()
	return results, nil
}
