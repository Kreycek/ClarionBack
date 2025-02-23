package daily

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

// Função de handler para a rota GET /dailys apenas diários sem documentos
func GetAllOnlyDailysHandler(w http.ResponseWriter, r *http.Request) {

	status, msg := clarion.TokenValido(w, r)

	if !status {
		http.Error(w, fmt.Sprintf("erro ao buscar usuários: %v", msg), http.StatusUnauthorized)
		return
	}

	// Conectar ao MongoDB
	client, err := db.ConnectMongoDB(clarion.ConectionString)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao conectar ao MongoDB: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.CloseMongoDB(client)

	// Obter todos os usuários
	dailys, err := GetDailys(client, clarion.DBName, "daily")
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao buscar usuários: %v", err), http.StatusInternalServerError)
		return
	}

	// Retornar a resposta com os dados dos usuários
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dailys); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}

func GetAllDailysHandler(w http.ResponseWriter, r *http.Request) {
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
	dailys, total, err := GetAllDaily(client, clarion.DBName, "daily", page, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao buscar diários: %v", err), http.StatusInternalServerError)
		return
	}

	// Criar resposta JSON com paginação
	response := map[string]any{
		"total":  total,
		"page":   page,
		"limit":  limit,
		"pages":  (total + limit - 1) / limit, // Calcula o número total de páginas
		"dailys": dailys,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}

func GetDailyByIdHandler(w http.ResponseWriter, r *http.Request) {
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
	user, err := GetDailyByID(client, clarion.DBName, "daily", id)
	if err != nil {
		http.Error(w, "Erro ao buscar diários", http.StatusInternalServerError)
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

func InsertDailyHandler(w http.ResponseWriter, r *http.Request) {
	// Validar o token de autenticação
	status, msg := clarion.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// Ler o corpo da requisição
	var daily models.Daily
	err := json.NewDecoder(r.Body).Decode(&daily)
	if err != nil {
		http.Error(w, "erro ao decodificar corpo da requisição", http.StatusBadRequest)
		return
	}

	if daily.Active == false {
		daily.Active = true
	}

	if daily.CreatedAt.IsZero() {
		daily.CreatedAt = time.Now()
	}
	// Conectar ao MongoDB
	client, err := db.ConnectMongoDB(clarion.ConectionString)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao conectar ao MongoDB: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.CloseMongoDB(client)

	// Inserir o usuário no MongoDB
	err = InsertDaily(client, clarion.DBName, "daily", daily)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao inserir Diário: %v", err), http.StatusInternalServerError)
		return
	}

	// Retornar a resposta com os dados dos usuários
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(daily); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}

}

func UpdateDailyHandler(w http.ResponseWriter, r *http.Request) {
	// Validar o token de autenticação
	status, msg := clarion.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("Erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// Decodificar o JSON recebido
	var daily models.Daily
	if err := json.NewDecoder(r.Body).Decode(&daily); err != nil {
		// http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		clarion.FormataRetornoHTTP(w, "Erro ao decodificar JSON", http.StatusBadRequest)

		return
	}

	// Verifica se o ID é válido
	if daily.ID.IsZero() {

		clarion.FormataRetornoHTTP(w, "ID do plano de contas inválido", http.StatusBadRequest)

		// http.Error(w, "ID do usuário inválido", http.StatusBadRequest)
		return
	}

	if daily.UpdatedAt.IsZero() {
		daily.UpdatedAt = time.Now()
	}

	// Criar o objeto de atualização
	update := bson.M{
		"$set": bson.M{
			"codDaily":     daily.CodDaily,
			"description":  daily.Description,
			"documents":    daily.Documents,
			"updatedAt":    daily.UpdatedAt,
			"idUserUpdate": daily.ID.Hex(),
			"active":       daily.Active,
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

	collection := client.Database(clarion.DBName).Collection("daily")
	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": daily.ID}, update)
	if err != nil {
		clarion.FormataRetornoHTTP(w, "Erro ao atualizar diário, Erro ao atualizar diário", http.StatusInternalServerError)

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
	clarion.FormataRetornoHTTP(w, "Diário atualizado com sucesso! Documento modificado", http.StatusOK)

}

// Função para verificar o nome de usuário e senha
func VerifyExistDailyHandler(w http.ResponseWriter, r *http.Request) {

	status, msg := clarion.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// Parse o corpo da requisição
	var daily models.DailyVerifyExistRequest

	err := json.NewDecoder(r.Body).Decode(&daily)
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
	collection := db.GetCollection(client, "clarion", "daily")
	// filter := bson.D{
	// 	{Key: "$or", Value: bson.A{
	// 		bson.D{{Key: "email", Value: userName}},
	// 	}},
	// }

	fmt.Println("dailyCod", daily)

	filter := bson.D{{Key: "codDaily", Value: daily.CodDaily}}

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

func SearchDailysHandler(w http.ResponseWriter, r *http.Request) {
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
		CodDaily    *string `json:"codDaily"`
		Description *string `json:"description"`
		Documents   []int   `json:"documents"`
		Page        int64   `json:"page"`
		Limit       int64   `json:"limit"`
	}

	// Decodificar o corpo da requisição JSON
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	fmt.Println(request)

	// Definir valores padrão para paginação
	if request.Page < 1 {
		request.Page = 1
	}
	if request.Limit < 1 {
		request.Limit = 10
	}

	// Buscar usuários com paginação
	dailys, total, err := SearchDailys(client, clarion.DBName, "daily", request.CodDaily, request.Description, request.Documents, request.Page, request.Limit)
	if err != nil {
		http.Error(w, "Erro ao buscar usuários", http.StatusInternalServerError)
		return
	}

	// Criar resposta JSON com paginação
	response := map[string]any{
		"total":  total,
		"page":   request.Page,
		"limit":  request.Limit,
		"pages":  (total + request.Limit - 1) / request.Limit, // Número total de páginas
		"dailys": dailys,
	}

	// Retornar resposta JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}
