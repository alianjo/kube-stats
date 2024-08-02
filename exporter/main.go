package main

import (
	"fmt"
	"log"
	"net/http"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	totalSpace = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "filesystem_total_bytes",
		Help: "Total filesystem space in bytes.",
	})
	usedSpace = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "filesystem_used_bytes",
		Help: "Used filesystem space in bytes.",
	})
	availableSpace = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "filesystem_available_bytes",
		Help: "Available filesystem space in bytes.",
	})
)

func init() {
	prometheus.MustRegister(totalSpace)
	prometheus.MustRegister(usedSpace)
	prometheus.MustRegister(availableSpace)
}

func updateMetrics() {
	var stat syscall.Statfs_t
	path := "/" // Path to the filesystem you want to check

	err := syscall.Statfs(path, &stat)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	available := stat.Bavail * uint64(stat.Bsize)
	total := stat.Blocks * uint64(stat.Bsize)
	used := total - available

	totalSpace.Set(float64(total))
	usedSpace.Set(float64(used))
	availableSpace.Set(float64(available))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("Request: %s %s from %s took %v", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
	})
}

func main() {
	updateMetrics()

	promHandler := promhttp.Handler()
	http.Handle("/metrics", loggingMiddleware(promHandler))

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
