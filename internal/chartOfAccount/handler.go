// Função de handler para a rota GET /users
package chartofaccount

import (
	clarion "Clarion"
	"Clarion/internal/db"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func GetAllChartOfAccountsHandler(w http.ResponseWriter, r *http.Request) {
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
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page < 1 {
		page = 1 // Padrão: primeira página
	}

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit < 1 {
		limit = 10 // Padrão: 10 registros por página
	}

	// Obter usuários paginados
	chartOfAccounts, total, err := GetAllChartOfAccount(client, clarion.DBName, "chartOfAccount", page, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao buscar perfis: %v", err), http.StatusInternalServerError)
		return
	}

	// Criar resposta JSON com paginação
	response := map[string]any{
		"total":           total,
		"page":            page,
		"limit":           limit,
		"pages":           (total + limit - 1) / limit, // Calcula o número total de páginas
		"chartOfAccounts": chartOfAccounts,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}

func SearchChartOfAccountsHandler(w http.ResponseWriter, r *http.Request) {
	// Verificar se a requisição é do tipo POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Validar Token
	status, msg := clarion.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// Conectar ao MongoDB
	client, err := db.ConnectMongoDB(clarion.ConectionString)
	if err != nil {
		http.Error(w, "Erro ao conectar ao MongoDB", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Definir estrutura para receber os parâmetros
	var request struct {
		CodAccount  *string `json:"codAccount"`
		Description *string `json:"description"`
		Year        []int   `json:"year"`
		Type        *string `json:"type"`
		Page        int64   `json:"page"`
		Limit       int64   `json:"limit"`
	}

	// Decodificar o corpo da requisição JSON
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	// Definir valores padrão para paginação
	if request.Page < 1 {
		request.Page = 1
	}
	if request.Limit < 1 {
		request.Limit = 10
	}

	// Buscar usuários com paginação
	charOfAccounts, total, err := SearchChartOfAccounts(client, clarion.DBName, "chartOfAccount", request.CodAccount, request.Description, request.Type, request.Year, request.Page, request.Limit)
	if err != nil {
		http.Error(w, "Erro ao buscar usuários", http.StatusInternalServerError)
		return
	}

	// Criar resposta JSON com paginação
	response := map[string]any{
		"total":           total,
		"page":            request.Page,
		"limit":           request.Limit,
		"pages":           (total + request.Limit - 1) / request.Limit, // Número total de páginas
		"chartOfAccounts": charOfAccounts,
	}

	// Retornar resposta JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}
