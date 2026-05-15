package httpclient

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// BenchmarkResult holds the results of an HTTP benchmark.
type BenchmarkResult struct {
	TotalRequests int
	TotalTime     time.Duration
	SuccessCount  int
	FailureCount  int
	Latencies     []time.Duration
	RPS           float64
}

// Benchmark runs a concurrent HTTP benchmark against a URL.
func Benchmark(url string, totalRequests int, concurrency int) BenchmarkResult {
	var wg sync.WaitGroup
	results := make(chan time.Duration, totalRequests)
	failures := make(chan bool, totalRequests)
	
	requestsPerWorker := totalRequests / concurrency
	startTime := time.Now()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := &http.Client{
				Timeout: 10 * time.Second,
			}
			for j := 0; j < requestsPerWorker; j++ {
				reqStart := time.Now()
				resp, err := client.Get(url)
				if err != nil {
					failures <- true
					continue
				}
				resp.Body.Close()
				results <- time.Since(reqStart)
			}
		}()
	}

	wg.Wait()
	close(results)
	close(failures)

	totalTime := time.Since(startTime)
	
	var latencies []time.Duration
	for l := range results {
		latencies = append(latencies, l)
	}

	failureCount := 0
	for range failures {
		failureCount++
	}

	return BenchmarkResult{
		TotalRequests: totalRequests,
		TotalTime:     totalTime,
		SuccessCount:  len(latencies),
		FailureCount:  failureCount,
		Latencies:     latencies,
		RPS:           float64(len(latencies)) / totalTime.Seconds(),
	}
}

// PrintBenchmarkResult prints the benchmark statistics.
func PrintBenchmarkResult(res BenchmarkResult) {
	fmt.Printf("\n--- Benchmark Results ---\n")
	fmt.Printf("Total Requests:    %d\n", res.TotalRequests)
	fmt.Printf("Total Time:        %v\n", res.TotalTime)
	fmt.Printf("Successful:        %d\n", res.SuccessCount)
	fmt.Printf("Failed:            %d\n", res.FailureCount)
	fmt.Printf("Requests/sec:      %.2f\n", res.RPS)

	if len(res.Latencies) > 0 {
		var totalLat time.Duration
		min := res.Latencies[0]
		max := res.Latencies[0]
		for _, l := range res.Latencies {
			totalLat += l
			if l < min {
				min = l
			}
			if l > max {
				max = l
			}
		}
		fmt.Printf("Average Latency:   %v\n", totalLat/time.Duration(len(res.Latencies)))
		fmt.Printf("Min Latency:       %v\n", min)
		fmt.Printf("Max Latency:       %v\n", max)
	}
}
