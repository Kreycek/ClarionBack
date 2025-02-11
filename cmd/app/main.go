package main

import (
	"Clarion/internal/auth"
	"encoding/json"
	"log"
	"net/http"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Definir um mapa com os dados
	data := map[string]string{
		"name": "ricardo",
	}

	// Definir o header como JSON
	w.Header().Set("Content-Type", "application/json")

	// Codificar o mapa de dados para JSON e enviar como resposta
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func main() {
	// Configura as rotas para autenticação e validação de token
	http.HandleFunc("/login", auth.VerifyUser) // Rota de login (gera o JWT)
	// http.HandleFunc("/validate", auth.ValidateToken) // Rota de validação do token

	http.HandleFunc("/teste", loginHandler)

	// http.HandleFunc("/createUser", auth.createUser) // Rota de validação do token

	// auth.CreateUser("rico", "654321")

	// Inicia o servidor na porta 8080
	log.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
