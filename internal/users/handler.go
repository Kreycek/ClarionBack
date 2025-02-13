package users

import (
	clarion "Clarion"
	"Clarion/internal/db"
	"Clarion/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	// var user models.User2
	// adas := json.NewDecoder(r.Body).Decode(&user)
	// fmt.Println(adas)
	// if adas != nil {
	// 	fmt.Println(w, "erro ao decodificar corpo da requisição", http.StatusBadRequest)
	// 	return
	// }
	// fmt.Println("AASAS", adas)

	// Ler o corpo da requisição
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "erro ao decodificar corpo da requisição", http.StatusBadRequest)
		return
	}

	// Conectar ao MongoDB
	client, err := db.ConnectMongoDB(clarion.ConectionString)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao conectar ao MongoDB: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.CloseMongoDB(client)

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
