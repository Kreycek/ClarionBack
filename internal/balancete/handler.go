package balancete

import (
	clarion "Clarion"
	"strconv"

	"Clarion/internal/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

/*
Função criada por Ricardo Silva Ferreira

	Inicio da criação 24/03/2025 15:00
	Data Final da criação : 24/03/2025 15:16
*/
func GenerateBalanceteReportHandler(w http.ResponseWriter, r *http.Request) {
	status, msg := clarion.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("erro ao buscar perfis: %v", msg), http.StatusUnauthorized)
		return
	}

	client, err := db.ConnectMongoDB(clarion.ConectionString)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao conectar ao MongoDB: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.CloseMongoDB(client)

	// Obter parâmetros de paginação
	query := r.URL.Query()

	initialYear, err := strconv.Atoi(query.Get("initialYear"))
	initialMonth, err := strconv.Atoi(query.Get("initialMonth"))
	endYear, err := strconv.Atoi(query.Get("endYear"))
	endMonth, err := strconv.Atoi(query.Get("endMonth"))

	// Obter usuários paginados
	balancete := GenerateBalanceteReport(initialYear, initialMonth, endYear, endMonth)

	for i := 0; i < len(balancete); i++ {
		if balancete[i].DebitValue > balancete[i].CreditValue {
			balancete[i].BalanceDebitValue = balancete[i].DebitValue - balancete[i].CreditValue
		}

		if balancete[i].CreditValue > balancete[i].DebitValue {
			balancete[i].BalandeCreditValue = balancete[i].CreditValue - balancete[i].DebitValue
		}

	}

	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao buscar diários: %v", err), http.StatusInternalServerError)
		return
	}

	// Criar resposta JSON com paginação
	response := map[string]any{
		"balancete": balancete,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}
