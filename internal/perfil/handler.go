package perfil

import (
	clarion "Clarion"
	"Clarion/internal/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Função de handler para a rota GET /users
func GetAllPerfilsHandler(w http.ResponseWriter, r *http.Request) {

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
	users, err := GetAllPerfil(client, clarion.DBName, "perfil")
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao buscar usuários: %v", err), http.StatusInternalServerError)
		return
	}

	// Retornar a resposta com os dados dos usuários
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Printf("erro ao codificar resposta JSON: %v", err)
	}
}
