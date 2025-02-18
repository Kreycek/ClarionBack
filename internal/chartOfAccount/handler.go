// Função de handler para a rota GET /users
package chartofaccount

import (
	clarion "Clarion"
	"Clarion/internal/db"
	"Clarion/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetChartOfAccountByIdHandler(w http.ResponseWriter, r *http.Request) {
	// Validar token
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

	// Extrair o ID da URL
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID não fornecido na URL", http.StatusBadRequest)
		return
	}

	// Verifica se o ID fornecido é válido
	_, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Buscar o usuário no banco de dados pelo ID
	user, err := GetChartOfAccountByID(client, clarion.DBName, "chartOfAccount", id)
	if err != nil {
		http.Error(w, "Erro ao buscar plano de contas", http.StatusInternalServerError)
		return
	}

	// Configurar o cabeçalho da resposta como JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Enviar o usuário como resposta JSON
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
		http.Error(w, "Erro ao codificar resposta", http.StatusInternalServerError)
	}
}

func GetAllChartOfAccountsHandler(w http.ResponseWriter, r *http.Request) {
	status, msg := clarion.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("erro ao buscar plano de contas: %v", msg), http.StatusUnauthorized)
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
		http.Error(w, fmt.Sprintf("erro ao buscar plano de contas: %v", err), http.StatusInternalServerError)
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

func InsertChartOfAccountHandler(w http.ResponseWriter, r *http.Request) {
	// Validar o token de autenticação
	status, msg := clarion.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// Ler o corpo da requisição
	var coa models.ChartOfAccount
	err := json.NewDecoder(r.Body).Decode(&coa)
	if err != nil {
		http.Error(w, "erro ao decodificar corpo da requisição", http.StatusBadRequest)
		return
	}

	if coa.Active == false {
		coa.Active = true
	}

	if coa.CreatedAt.IsZero() {
		coa.CreatedAt = time.Now()
	}
	// Conectar ao MongoDB
	client, err := db.ConnectMongoDB(clarion.ConectionString)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao conectar ao MongoDB: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.CloseMongoDB(client)

	// Inserir o usuário no MongoDB
	err = InsertChartOfAccount(client, clarion.DBName, "chartOfAccount", coa)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao inserir usuário: %v", err), http.StatusInternalServerError)
		return
	}

	// Retornar a resposta com os dados dos usuários
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(coa); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}

}

func UpdateChartOfAccountHandler(w http.ResponseWriter, r *http.Request) {
	// Validar o token de autenticação
	status, msg := clarion.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("Erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// Decodificar o JSON recebido
	var coa models.ChartOfAccount
	if err := json.NewDecoder(r.Body).Decode(&coa); err != nil {
		// http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		clarion.FormataRetornoHTTP(w, "Erro ao decodificar JSON", http.StatusBadRequest)

		return
	}

	// Verifica se o ID é válido
	if coa.ID.IsZero() {

		clarion.FormataRetornoHTTP(w, "ID do plano de contas inválido", http.StatusBadRequest)

		// http.Error(w, "ID do usuário inválido", http.StatusBadRequest)
		return
	}

	if coa.UpdatedAt.IsZero() {
		coa.UpdatedAt = time.Now()
	}

	// Criar o objeto de atualização
	update := bson.M{
		"$set": bson.M{
			"codAccount":   coa.CodAccount,
			"description":  coa.Description,
			"type":         coa.Type,
			"year":         coa.Year,
			"updatedAt":    coa.UpdatedAt,
			"idUserUpdate": coa.ID.Hex(),
			"active":       coa.Active,
		},
	}

	// Conectar ao MongoDB e atualizar o usuário
	client, err := db.ConnectMongoDB(clarion.ConectionString)
	if err != nil {
		clarion.FormataRetornoHTTP(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)

		// http.Error(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(clarion.DBName).Collection("chartOfAccount")
	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": coa.ID}, update)
	if err != nil {
		clarion.FormataRetornoHTTP(w, "Erro ao atualizar plano de contas, Erro ao atualizar plano de contas", http.StatusInternalServerError)

		// log.Println("Erro ao atualizar usuário:", err)
		// http.Error(w, "Erro ao atualizar usuário", http.StatusInternalServerError)
		return
	}

	// Verifica se algum documento foi modificado
	if result.ModifiedCount == 0 {
		clarion.FormataRetornoHTTP(w, "Nenhuma alteração realizada", http.StatusOK)

		// http.Error(w, "Nenhuma alteração realizada", http.StatusNotModified)
		return
	}

	// Responder com sucesso
	clarion.FormataRetornoHTTP(w, "Plano de contas atualizado com sucesso! Documento modificado", http.StatusOK)

}

func UpdateYearForAllDocumentsHandler(w http.ResponseWriter, r *http.Request) {
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
		Year int `json:"year"`
	}

	// Decodificar o corpo da requisição JSON
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	// Buscar usuários com paginação
	err = UpdateYearForAllDocuments(client, clarion.DBName, "chartOfAccount", request.Year)
	if err != nil {
		clarion.FormataRetornoHTTP(w, false, http.StatusOK)
		return
	}

	// Retornar resposta JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(true); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}

// Função para verificar o nome de usuário e senha
func VerifyExistChartOfAccountHandler(w http.ResponseWriter, r *http.Request) {

	status, msg := clarion.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// Parse o corpo da requisição
	var account models.ChartOfAccountVerifyExistRequest

	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da requisição", http.StatusBadRequest)
		return
	}

	// fmt.Println("email", email)
	// Conectar ao MongoDB
	client, err := db.ConnectMongoDB(clarion.ConectionString)
	if err != nil {
		http.Error(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
		return
	}
	defer db.CloseMongoDB(client)

	// Obter a coleção de usuários
	collection := db.GetCollection(client, "clarion", "chartOfAccount")
	// filter := bson.D{
	// 	{Key: "$or", Value: bson.A{
	// 		bson.D{{Key: "email", Value: userName}},
	// 	}},
	// }

	fmt.Println("account", account)

	filter := bson.D{
		{Key: "codAccount", Value: bson.M{
			"$regex":   account.CodAccount,
			"$options": "i", // Ignora diferença entre maiúsculas e minúsculas
		}},
	}

	var result bson.M
	err = collection.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			clarion.FormataRetornoHTTP(w, false, http.StatusOK)
			return
		}
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	// Se encontrou um documento, retorna true
	clarion.FormataRetornoHTTP(w, true, http.StatusOK)

}
