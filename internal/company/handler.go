package company

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

func GetAllAutoCompletesHandler(w http.ResponseWriter, r *http.Request) {
	status, msg := clarion.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("erro ao buscar empresas: %v", msg), http.StatusUnauthorized)
		return
	}

	client, err := db.ConnectMongoDB(clarion.ConectionString)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao conectar ao MongoDB: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.CloseMongoDB(client)

	// Obter parâmetro de pesquisa por nome
	query := r.URL.Query()
	name := query.Get("name")

	// Obter empresas filtradas por nome (ou todas se name estiver vazio)
	companys, err := GetAllAutoComplete(client, clarion.DBName, "company", name)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao buscar empresas: %v", err), http.StatusInternalServerError)
		return
	}

	// Criar resposta JSON
	response := map[string]any{
		"companys": companys,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}

func GetAllCompanysHandler(w http.ResponseWriter, r *http.Request) {
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
	companys, total, err := GetAllCompany(client, clarion.DBName, "company", page, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao buscar diários: %v", err), http.StatusInternalServerError)
		return
	}

	// Criar resposta JSON com paginação
	response := map[string]any{
		"total":    total,
		"page":     page,
		"limit":    limit,
		"pages":    (total + limit - 1) / limit, // Calcula o número total de páginas
		"companys": companys,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}

func GetCompanyByIdHandler(w http.ResponseWriter, r *http.Request) {
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
	user, err := GetCompanyByID(client, clarion.DBName, "company", id)
	if err != nil {
		http.Error(w, "Erro ao buscar empresas", http.StatusInternalServerError)
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

func InsertCompanyHandler(w http.ResponseWriter, r *http.Request) {
	// Validar o token de autenticação
	status, msg := clarion.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// Ler o corpo da requisição
	var company models.Company
	err := json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		http.Error(w, "erro ao decodificar corpo da requisição", http.StatusBadRequest)
		return
	}

	if company.Active == false {
		company.Active = true
	}

	if company.CreatedAt.IsZero() {
		company.CreatedAt = time.Now()
	}
	// Conectar ao MongoDB
	client, err := db.ConnectMongoDB(clarion.ConectionString)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao conectar ao MongoDB: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.CloseMongoDB(client)

	// Inserir o usuário no MongoDB
	err = InsertCompany(client, clarion.DBName, "company", company)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao inserir Empresa: %v", err), http.StatusInternalServerError)
		return
	}

	// Retornar a resposta com os dados dos usuários
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(company); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}

}

func UpdateCompanyHandler(w http.ResponseWriter, r *http.Request) {
	// Validar o token de autenticação
	status, msg := clarion.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("Erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// Decodificar o JSON recebido
	var company models.Company
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		// http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		clarion.FormataRetornoHTTP(w, "Erro ao decodificar JSON", http.StatusBadRequest)

		return
	}

	// Verifica se o ID é válido
	if company.ID.IsZero() {

		clarion.FormataRetornoHTTP(w, "ID do plano de contas inválido", http.StatusBadRequest)

		// http.Error(w, "ID do usuário inválido", http.StatusBadRequest)
		return
	}

	if company.UpdatedAt.IsZero() {
		company.UpdatedAt = time.Now()
	}

	// Criar o objeto de atualização
	update := bson.M{
		"$set": bson.M{
			"codCompany":      company.CodCompany,
			"name":            company.Name,
			"cae":             company.CAE,
			"documents":       company.Documents,
			"mainActivity":    company.MainActivity,
			"otherActivities": company.OtherActivities,
			"legalNature":     company.LegalNature,
			"socialCapital":   company.SocialCapital,
			"nationalCapital": company.NationalCapital,
			"extraCapital":    company.ExtraCapital,
			"publicCapital":   company.PublicCapital,
			"vatRegime":       company.VATRegime,
			"email":           company.Email,
			"webSite":         company.WebSite,
			"active":          company.Active,
			"exercise":        company.Exercise,
			"phone":           company.Phone,
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

	collection := client.Database(clarion.DBName).Collection("company")
	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": company.ID}, update)
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
	clarion.FormataRetornoHTTP(w, "Empresa atualizada com sucesso! Documento modificado", http.StatusOK)

}

// Função para verificar o nome de usuário e senha
func VerifyExistCompanyHandler(w http.ResponseWriter, r *http.Request) {

	status, msg := clarion.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// Parse o corpo da requisição
	var company models.Company

	err := json.NewDecoder(r.Body).Decode(&company)
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
	collection := db.GetCollection(client, "clarion", "company")
	// filter := bson.D{
	// 	{Key: "$or", Value: bson.A{
	// 		bson.D{{Key: "email", Value: userName}},
	// 	}},
	// }

	filter := bson.D{{Key: "codCompany", Value: company.CodCompany}}

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

func SearchCompanysHandler(w http.ResponseWriter, r *http.Request) {
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
		CodCompany string `json:"codCompany"`
		Document   string `json:"document"`
		Address    string `json:"address"`
		Page       int64  `json:"page"`
		Limit      int64  `json:"limit"`
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
	companys, total, err := SearchCompany(client, clarion.DBName, "company", &request.CodCompany, &request.Document, &request.Address, request.Page, request.Limit)
	if err != nil {
		http.Error(w, "Erro ao buscar empresas", http.StatusInternalServerError)
		return
	}

	// Criar resposta JSON com paginação
	response := map[string]any{
		"total":    total,
		"page":     request.Page,
		"limit":    request.Limit,
		"pages":    (total + request.Limit - 1) / request.Limit, // Número total de páginas
		"companys": companys,
	}

	// Retornar resposta JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}
