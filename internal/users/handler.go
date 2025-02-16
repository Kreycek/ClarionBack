package users

import (
	clarion "Clarion"
	"Clarion/internal/db"
	"Clarion/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Função de handler para a rota GET /users
func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {

	status, msg := clarion.TokenValido(w, r)

	if !status {
		http.Error(w, fmt.Sprintf("erro ao buscar perfis: %v", msg), http.StatusUnauthorized)
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
	users, err := GetAllUsers(client, clarion.DBName, "user")
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao buscar perfis: %v", err), http.StatusInternalServerError)
		return
	}

	// Retornar a resposta com os dados dos usuários
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}

func InsertUserHandler(w http.ResponseWriter, r *http.Request) {
	// Validar o token de autenticação
	status, msg := clarion.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// Ler o corpo da requisição
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "erro ao decodificar corpo da requisição", http.StatusBadRequest)
		return
	}

	if user.Active == false {
		user.Active = true
	}

	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	// Conectar ao MongoDB
	client, err := db.ConnectMongoDB(clarion.ConectionString)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao conectar ao MongoDB: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.CloseMongoDB(client)

	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

		if err != nil {
			log.Println("Erro ao gerar hash da senha:", err)
			return
		}

		// Atribuindo a senha hashada ao campo Password
		user.Password = string(hashedPassword)
	}

	// Inserir o usuário no MongoDB
	err = InsertUser(client, clarion.DBName, "user", user)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao inserir usuário: %v", err), http.StatusInternalServerError)
		return
	}

	// Retornar a resposta com os dados dos usuários
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}

}

// Variáveis globais

// Função para verificar o nome de usuário e senha
func VerifyExistUser(w http.ResponseWriter, r *http.Request) {

	status, msg := clarion.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// Parse o corpo da requisição
	var email models.EmailRequest

	err := json.NewDecoder(r.Body).Decode(&email)
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
	collection := db.GetCollection(client, "clarion", "user")
	// filter := bson.D{
	// 	{Key: "$or", Value: bson.A{
	// 		bson.D{{Key: "email", Value: userName}},
	// 	}},
	// }

	filter := bson.D{{Key: "email", Value: email.Email}}

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

func SearchUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Configurar o cabeçalho da resposta como JSON
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

	// Obtendo os parâmetros da URL

	var user models.User

	erre := json.NewDecoder(r.Body).Decode(&user)
	if erre != nil {
		http.Error(w, "Erro ao ler o corpo da requisição", http.StatusBadRequest)
		return
	}

	name := user.Name
	email := user.Email
	perfilStr := user.Perfil

	// Convertendo "perfil" para inteiro (se existir)
	// var perfil *int
	// if perfilStr != "" {
	// 	perfilVal, err := strconv.Atoi(perfilStr)
	// 	if err == nil {
	// 		perfil = &perfilVal
	// 	}
	// }

	// Busca usuários no MongoDB
	users, err := SearchUsers(client, clarion.DBName, "user", &name, &email, perfilStr)
	if err != nil {
		http.Error(w, "Erro ao buscar usuários", http.StatusInternalServerError)
		return
	}

	// count := len(users) // Contagem dos documentos retornados
	// fmt.Printf("Total de usuários encontrados: %d\n", users[0].Name)~

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}

}
func GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {
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
	user, err := GetUserByID(client, clarion.DBName, "user", id)
	if err != nil {
		http.Error(w, "Erro ao buscar usuário", http.StatusInternalServerError)
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

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Validar o token de autenticação
	status, msg := clarion.TokenValido(w, r)
	if !status {
		http.Error(w, fmt.Sprintf("Erro ao validar token: %v", msg), http.StatusUnauthorized)
		return
	}

	// Decodificar o JSON recebido
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		// http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		clarion.FormataRetornoHTTP(w, "Erro ao decodificar JSON", http.StatusBadRequest)

		return
	}

	// Verifica se o ID é válido
	if user.ID.IsZero() {

		clarion.FormataRetornoHTTP(w, "ID do usuário inválido", http.StatusBadRequest)

		// http.Error(w, "ID do usuário inválido", http.StatusBadRequest)
		return
	}

	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = time.Now()
	}

	// Criar o objeto de atualização
	update := bson.M{
		"$set": bson.M{
			"name":           user.Name,
			"lastName":       user.LastName,
			"passportNumber": user.PassportNumber,
			"perfil":         user.Perfil,
			"UpdatedAt":      user.UpdatedAt,
			"IdUserUpdate":   user.ID.Hex(),
			"active":         user.Active,
		},
	}

	// Se uma nova senha for fornecida, gerar um hash
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			clarion.FormataRetornoHTTP(w, "Erro ao gerar hash da senha,Erro ao processar a senha ", http.StatusInternalServerError)

			// log.Println("Erro ao gerar hash da senha:", err)
			// http.Error(w, "Erro ao processar a senha", http.StatusInternalServerError)
			return
		}
		update["$set"].(bson.M)["password"] = string(hashedPassword)
	}

	// Conectar ao MongoDB e atualizar o usuário
	client, err := db.ConnectMongoDB(clarion.ConectionString)
	if err != nil {
		clarion.FormataRetornoHTTP(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)

		// http.Error(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(clarion.DBName).Collection("user")
	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": user.ID}, update)
	if err != nil {
		clarion.FormataRetornoHTTP(w, "Erro ao atualizar usuário, Erro ao atualizar usuário", http.StatusInternalServerError)

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
	clarion.FormataRetornoHTTP(w, "Usuário atualizado com sucesso! Documento modificado", http.StatusOK)

}
