# Monitoring Analyze Test

This project implements a complete, virtualized monitoring environment for analyzing transaction data from CSV files. During data analysis, Python was used to understand the anomaly type.

The stack uses Prometheus for collecting and storing metrics, Grafana for visualization and dashboards, and a custom Go Exporter for reading the data and exposing it to Prometheus. The entire environment is orchestrated with Docker and Docker Compose, making it portable and easy to configure.