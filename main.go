package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	StatusCode int
	Duration   time.Duration
}

func worker(id int, wg *sync.WaitGroup, jobs <-chan int, url string, results chan<- Result) {
	defer wg.Done()
	client := &http.Client{}
	for range jobs {
		start := time.Now()
		resp, err := client.Get(url)
		duration := time.Since(start)
		if err != nil {
			results <- Result{StatusCode: 0, Duration: duration}
			continue
		}
		results <- Result{StatusCode: resp.StatusCode, Duration: duration}
		resp.Body.Close()
	}
}

func main() {
	url := flag.String("url", "", "URL do serviço que sera testado")
	totalRequests := flag.Int("requests", 10, "Número total de requisições")
	concurrency := flag.Int("concurrency", 1, "Número de chamadas paralelas")
	flag.Parse()
	if *url == "" {
		fmt.Println("A URL é obrigatória. Use --url para definir.")
		return
	}
	start := time.Now()
	jobs := make(chan int, *totalRequests)
	results := make(chan Result, *totalRequests)
	var wg sync.WaitGroup
	for w := 0; w < *concurrency; w++ {
		wg.Add(1)
		go worker(w, &wg, jobs, *url, results)
	}
	for j := 0; j < *totalRequests; j++ {
		jobs <- j
	}
	close(jobs)
	wg.Wait()
	close(results)
	statusCount := make(map[int]int)
	var totalDuration time.Duration
	var minDuration, maxDuration time.Duration
	minDuration = time.Duration(1<<63 - 1)
	for r := range results {
		statusCount[r.StatusCode]++
		totalDuration += r.Duration
		if r.Duration < minDuration {
			minDuration = r.Duration
		}
		if r.Duration > maxDuration {
			maxDuration = r.Duration
		}
	}
	avgDuration := totalDuration / time.Duration(*totalRequests)
	duration := time.Since(start)
	fmt.Println("-------------------")
	fmt.Println("Relatório de Teste de Carga")
	fmt.Printf("Tempo total: %v\n", duration)
	fmt.Printf("Total de requisições: %d\n", *totalRequests)
	fmt.Printf("Respostas HTTP 200: %d\n", statusCount[200])
	fmt.Println("Distribuição de outros códigos de status HTTP:")
	for code, count := range statusCount {
		if code != 200 {
			fmt.Printf("HTTP %d: %d\n", code, count)
		}
	}
	fmt.Println("-------------------")
	fmt.Printf("\nLatencia total:\n")
	fmt.Printf("Minima: %v\n", minDuration)
	fmt.Printf("Maxima: %v\n", maxDuration)
	fmt.Printf("Media: %v\n", avgDuration)
}
