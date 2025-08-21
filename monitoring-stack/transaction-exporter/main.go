package main

import (
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// We define Prometheus's "gauges".
// A Gauge is a numeric value that can go up and down.
var (
	transactionsApproved = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "transactions_approved_total",
		Help: "Number of transactions approved in the last minute.",
	})
	transactionsDenied = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "transactions_denied_total",
		Help: "Number of transactions declined in the last minute.",
	})
	transactionsFailed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "transactions_failed_total",
		Help: "Number of failed transactions in the last minute.",
	})
	transactionsReversed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "transactions_reversed_total",
		Help: "Number of transactions reversed in the last",
	})
)

var (
	checkoutStats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "checkout_stats_count",
			Help: "Hourly checkout statistics.",
		},
		[]string{"hour", "comparison"},
	)
)

func init() {
	prometheus.MustRegister(transactionsApproved)
	prometheus.MustRegister(transactionsDenied)
	prometheus.MustRegister(transactionsFailed)
	prometheus.MustRegister(transactionsReversed)
	prometheus.MustRegister(checkoutStats)
}

func processTransactionsCSV(filePath string) {
	log.Println("Starting CSV data simulation...")

	// Infinite loop to play the data continuously
	for {
		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("Error opening CSV %v. Trying again in 10s", err)
			time.Sleep(10 * time.Second)
			continue
		}
		defer file.Close()

		reader := csv.NewReader(file)
		reader.Read()

		for {
			record, err := reader.Read()
			if err == io.EOF {
				log.Println("End of CSV file. Restarting the simulation.")
				break
			}
			if err != nil {
				log.Printf("Error read CSV line: %v", err)
				continue
			}

			status := strings.ToLower(record[1])
			count, _ := strconv.ParseFloat(record[2], 64)

			switch status {
			case "approved":
				transactionsApproved.Set(count)
			case "denied":
				transactionsDenied.Set(count)
			case "failed":
				transactionsFailed.Set(count)
			case "reversed", "backend_reversed":
				transactionsReversed.Add(count)
			}

			if status != "approved" {
				transactionsApproved.Set(0)
			}
			if status != "denied" {
				transactionsDenied.Set(0)
			}
			if status != "failed" {
				transactionsFailed.Set(0)
			}
			if !strings.Contains(status, "reversed") {
				transactionsReversed.Set(0)
			}

			// Wait 1 second to simulate the passing of 1 minute in the log.

			log.Printf("Update métric: %s = %.0f", status, count)
			time.Sleep(1 * time.Second)

		}
	}

}

func processCheckoutCSV(filePath string) {
	log.Println("Starting data simulation of CHECKOUT...")
	for {
		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("CHECKOUT: Error when open CSV: %v, Trying again.", err)
			time.Sleep(10 * time.Second)
			continue
		}

		reader := csv.NewReader(file)
		headers, _ := reader.Read()

		for {
			record, err := reader.Read()
			if err == io.EOF {
				log.Println("CHECKOUT: End of day, restarting checkout simulation")
				break
			}
			if err != nil {
				continue
			}

			hour := record[0]
			log.Printf("CHECKOUT: Processing data for hour: %s", hour)

			for i := 1; i < len(record); i++ {
				comparisonType := headers[i]
				value, _ := strconv.ParseFloat(record[i], 64)
				checkoutStats.WithLabelValues(hour, comparisonType).Set(value)
			}
			time.Sleep(5 * time.Second)
		}
		file.Close()
	}
}

func main() {
	go processTransactionsCSV("/data/transactions.csv")
	go processCheckoutCSV("/data/checkout_2.csv")

	http.Handle("/metrics", promhttp.Handler())
	log.Println("Exporter initialized in port :8081. Metrics endpoints in /metrics")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
