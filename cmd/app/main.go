package main

import (
	clarion "Clarion"
	"Clarion/internal/auth"
	chartofaccount "Clarion/internal/chartOfAccount"
	"Clarion/internal/daily"
	"Clarion/internal/perfil"
	"Clarion/internal/users"
	"encoding/json"
	"log"
	"net/http"

	"github.com/rs/cors"
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

	c := cors.New(cors.Options{
		AllowedOrigins: []string{clarion.UrlSite}, // Permitindo o domínio de onde vem a requisição
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})
	// Configura as rotas para autenticação e validação de token
	http.HandleFunc("/login", auth.VerifyUser)       // Rota de login (gera o JWT)
	http.HandleFunc("/validate", auth.ValidateToken) // Rota de validação do token
	http.HandleFunc("/getPerfis", perfil.GetAllPerfilsHandler)

	//USUÁRIOS
	http.HandleFunc("/addUser", users.InsertUserHandler)
	http.HandleFunc("/getAllUsers", users.GetAllUsersHandler)
	http.HandleFunc("/verifyExistUser", users.VerifyExistUser)
	http.HandleFunc("/searchUsers", users.SearchUsersHandler)
	http.HandleFunc("/getUserById", users.GetUserByIdHandler)
	http.HandleFunc("/updateUser", users.UpdateUserHandler)

	//PLANO DE CONTAS
	http.HandleFunc("/getAllChartOfAccount", chartofaccount.GetAllChartOfAccountsHandler)
	http.HandleFunc("/searchChartOfAccounts", chartofaccount.SearchChartOfAccountsHandler)
	http.HandleFunc("/getChartOfAccountById", chartofaccount.GetChartOfAccountByIdHandler)
	http.HandleFunc("/insertChartOfAccount", chartofaccount.InsertChartOfAccountHandler)
	http.HandleFunc("/updateChartOfAccount", chartofaccount.UpdateChartOfAccountHandler)
	http.HandleFunc("/updateAllYearOfAccounts", chartofaccount.UpdateYearForAllDocumentsHandler)
	http.HandleFunc("/VerifyExistChartOfAccount", chartofaccount.VerifyExistChartOfAccountHandler)

	//DIÁRIO
	http.HandleFunc("/getAllDailys", daily.GetAllDailysHandler)
	http.HandleFunc("/getAllOnlyDailys", daily.GetAllOnlyDailysHandler)
	http.HandleFunc("/getDailyById", daily.GetDailyByIdHandler)
	http.HandleFunc("/InsertDaily", daily.InsertDailyHandler)
	http.HandleFunc("/UpdateDaily", daily.UpdateDailyHandler)
	http.HandleFunc("/VerifyExistDaily", daily.VerifyExistDailyHandler)
	http.HandleFunc("/SearchDailys", daily.SearchDailysHandler)

	//TESTE
	http.HandleFunc("/teste", loginHandler)
	handler := c.Handler(http.DefaultServeMux)

	// http.HandleFunc("/createUser", auth.createUser) // Rota de validação do token

	// auth.CreateUser("rico", "654321")

	// Inicia o servidor na porta 8080
	log.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
