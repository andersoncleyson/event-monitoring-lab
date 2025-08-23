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

var (
	transactionStats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "transaction_stats_total",
			Help: "Transaction statistics by status.",
		},
		[]string{"status"},
	)

	checkoutStats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "checkout_stats_count",
			Help: "Hourly checkout statistics.",
		},
		[]string{"hour", "comparison"},
	)
)

func init() {
	prometheus.MustRegister(transactionStats)
	prometheus.MustRegister(checkoutStats)
}

func processTransactionsCSV(filePath string) {
	log.Println("Starting CSV data simulation...")

	transactionStats.WithLabelValues("approved").Set(0)
	transactionStats.WithLabelValues("denied").Set(0)
	transactionStats.WithLabelValues("failed").Set(0)
	transactionStats.WithLabelValues("reversed").Set(0)
	transactionStats.WithLabelValues("refunded").Set(0)

	// Infinite loop to play the data continuously
	for {
		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("Error opening CSV %v. Trying again in 10s", err)
			time.Sleep(10 * time.Second)
			continue
		}

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

			status := strings.TrimSpace(strings.ToLower(record[1]))
			count, _ := strconv.ParseFloat(record[2], 64)

			transactionStats.WithLabelValues(status).Set(count)

			time.Sleep(10 * time.Millisecond)

		}
		file.Close()
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
