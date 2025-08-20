package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

// BaselineStats armazena a média e o desvio padrão para status.
type BaselineStats struct {
	Mean float64
	Std  float64
}

// MonitoringService encapsula a lógina do monitoramento
type MonitoringService struct {
	baselineStats map[string]BaselineStats
}

type App struct {
	service *MonitoringService
}

// NewMonitoringService cria e inicializa o serviço de monitoramento.
func NewMonitoringService(csvFilePath string) (*MonitoringService, error) {
	service := &MonitoringService{
		baselineStats: make(map[string]BaselineStats),
	}

	err := service.calculateBaseline(csvFilePath)
	if err != nil {
		return nil, err
	}
	return service, nil
}

// calculateBaseline lê o CSV e calcula as métricas da base.
func (s *MonitoringService) calculateBaseline(filepath string) error {
	// 1. Abrir e ler o arquivo CSV
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Read() // Pula o cabeçalho

	// 2. Agregar dados por minuto
	transactionsPerMinute := make(map[time.Time]map[string]int)
	layout := "2006-01-02 15:04:05"

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		t, err := time.Parse(layout, record[0])
		if err != nil {
			log.Printf("Aviso: Falha ao parsear a data '%s' : %v", record[0], err)
			continue
		}
		minute := t.Truncate(time.Minute)
		status := record[1]
		count, _ := strconv.Atoi(record[2])

		if _, ok := transactionsPerMinute[minute]; !ok {
			transactionsPerMinute[minute] = make(map[string]int)
		}
		transactionsPerMinute[minute][status] += count
	}

	// 3. Calcular Média e Descio Padrão
	statsData := map[string][]float64{
		"denied":         {},
		"failed":         {},
		"reversed_total": {},
	}

	for _, data := range transactionsPerMinute {
		statsData["denied"] = append(statsData["denied"], float64(data["denied"]))
		statsData["failed"] = append(statsData["failed"], float64(data["failed"]))
		reversed := float64(data["reversed"] + data["backend_reversed"])
		statsData["reversed_total"] = append(statsData["reversed_total"], reversed)
	}

	for status, values := range statsData {
		var sum, mean, std float64
		n := float64(len(values))

		for _, v := range values {
			sum += v
		}

		mean = sum / n

		for _, v := range values {
			std += math.Pow(v-mean, 2)
		}

		std = math.Sqrt(std / n)

		s.baselineStats[status] = BaselineStats{Mean: mean, Std: std}
		log.Printf("Baseline para '%s' : Média=%.2f, Desvio Padrão=%.2f", status, mean, std)
	}

	return nil
}

// CheckForAnomalies verifica os dados de uma transação contra o baseline.
func (s *MonitoringService) CheckForAnomalies(data map[string]int) []string {
	const stdMultiplier = 3.0
	var alerts []string

	for status, stats := range s.baselineStats {
		count := data[status]
		threshold := stats.Mean + (stdMultiplier * stats.Std)

		if float64(count) > threshold && count > 1 {
			alertMsg := fmt.Sprintf(
				"Status '%s' acima do normal. Contagem: %d, Limiar: %.2f",
				status, count, threshold,
			)
			alerts = append(alerts, alertMsg)
		}
	}
	return alerts
}

// monitorHandler é o handler HTTP para endpoing /monitor
func (app *App) monitorHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		Counts map[string]int `json:"counts"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	alerts := app.service.CheckForAnomalies(requestData.Counts)

	response := make(map[string]interface{})
	if len(alerts) == 0 {
		response["recommendation"] = "ok"
	} else {
		response["recommendation"] = "alertar"
		response["details"] = alerts
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Inicializa o serviço com os dados CSV.
	service, err := NewMonitoringService("transactions.csv")
	if err != nil {
		log.Fatalf("Falha ao inicializar o serviço de monitoramento: %v", err)
	}

	// Cria uma instância do nosso handler, injetando o serviço nele.
	app := &App{
		service: service,
	}

	// Registra o método do nosso handler para a rota "/monitor".
	// Agora a chamada está correta e clara.
	http.HandleFunc("/monitor", app.monitorHandler)

	// Adiciona uma rota raiz para fácil verificação
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Serviço de monitoramento está no ar. Use o endpoint POST /monitor para enviar dados.")
	})

	log.Println("Serviço de monitoramento iniciado em http://localhost:8080")
	// Inicia o servidor HTTP
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
