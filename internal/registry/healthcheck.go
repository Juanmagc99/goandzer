package registry

import (
	"juanmagc99/goandzer/internal/balancer"
	"net/http"
	"time"
)

func StartHealthChecker(s ServiceConfig, healthPath string, interval time.Duration, rr *balancer.RoundRobin) {
	doCheck := func() {
		var alive []string
		client := &http.Client{Timeout: 2 * time.Second}
		for _, target := range s.Targets {
			healthURL := target + healthPath
			if resp, err := client.Get(healthURL); err == nil {
				resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					alive = append(alive, target)
				}
			}
		}
		rr.SetTargets(alive)
	}

	doCheck()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		doCheck()
	}
}
