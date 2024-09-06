package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var (
	requestCount int32
	serverUrls   = [3]string{
		"http://localhost:8080/",
		"http://localhost:8081/",
		"http://localhost:8082/",
	}

	healthyUrls []string
	mu          sync.RWMutex
)

func main() {

	var healthCheckPeriod int
	flag.IntVar(&healthCheckPeriod, "healthCheckPeriod", 10, "specify the health check in seconds")

	flag.Parse()

	go func() {
		for range time.Tick(time.Second * time.Duration(healthCheckPeriod)) {
			healthCheck()
		}
	}()
	http.HandleFunc("/", reqhandler)
	fmt.Println("Load balancer is running at http://localhost:8000")

	log.Fatal(http.ListenAndServe(":8000", nil))
}

func reqhandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()

	currentCount := atomic.LoadInt32(&requestCount)
	atomic.AddInt32(&requestCount, 1)

	req, err := http.NewRequest(r.Method, healthyUrls[currentCount%int32(len(healthyUrls))], r.Body)

	mu.RUnlock()

	if err != nil {
		http.Error(w, "Failed to create a new request", http.StatusInternalServerError)
		return
	}

	currentServer := currentCount % int32(len(healthyUrls))
	log.Printf("Forwarding request #%d to: %s .....\n", currentServer+1, healthyUrls[currentServer])

	for key, values := range r.Header {
		for _, val := range values {
			req.Header.Add(key, val)
		}
	}

	client := &http.Client{}
	res, error := client.Do(req)

	if error != nil {
		log.Printf("Failed to forward request: %v\n", error)
		http.Error(w, "Failed to forward the request", http.StatusInternalServerError)
		return
	}

	defer res.Body.Close()

	for name, values := range res.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	log.Printf("Response from server: %s %d\n=========================", res.Proto, res.StatusCode)
	io.Copy(w, res.Body)
}

func healthCheck() {

	mu.Lock()

	healthyUrls = nil

	for _, url := range serverUrls {
		res, err := http.Get(url)

		if err != nil || res.StatusCode != http.StatusOK {
			continue
		}
		healthyUrls = append(healthyUrls, url)
	}
	mu.Unlock()

	log.Printf("Healthy Urls:  %v\n", healthyUrls)
}
