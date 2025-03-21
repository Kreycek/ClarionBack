package main

import (
	clarion "Clarion"
	"Clarion/internal/auth"
	"Clarion/internal/balancete"
	chartofaccount "Clarion/internal/chartOfAccount"
	"Clarion/internal/company"
	costcenter "Clarion/internal/costCenter"
	"Clarion/internal/daily"

	"Clarion/internal/movement"
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
	http.HandleFunc("/GetAllCoaAutoComplete", chartofaccount.GetAllCoaAutoCompleteHandler)

	//DIÁRIO
	http.HandleFunc("/getAllDailys", daily.GetAllDailysHandler)
	http.HandleFunc("/getAllOnlyDailys", daily.GetAllOnlyDailysHandler)
	http.HandleFunc("/getDailyById", daily.GetDailyByIdHandler)
	http.HandleFunc("/InsertDaily", daily.InsertDailyHandler)
	http.HandleFunc("/UpdateDaily", daily.UpdateDailyHandler)
	http.HandleFunc("/VerifyExistDaily", daily.VerifyExistDailyHandler)
	http.HandleFunc("/SearchDailys", daily.SearchDailysHandler)
	http.HandleFunc("/GetAllOnlyDailyActive", daily.GetAllOnlyDailyActiveHandler)

	//	COMPANYS
	http.HandleFunc("/GetAllCompanys", company.GetAllCompanysHandler)
	http.HandleFunc("/GetCompanyById", company.GetCompanyByIdHandler)
	http.HandleFunc("/InsertCompany", company.InsertCompanyHandler)
	http.HandleFunc("/UpdateCompany", company.UpdateCompanyHandler)
	http.HandleFunc("/VerifyExistCompany", company.VerifyExistCompanyHandler)
	http.HandleFunc("/SearchCompanys", company.SearchCompanysHandler)
	http.HandleFunc("/GetAllCompanyAutoComple", company.GetAllCompanyAutoCompleteHandler)

	//MOVIMENTO
	http.HandleFunc("/GetAllMovements", movement.GetAllMovementsHandler)
	http.HandleFunc("/GetMovementById", movement.GetMovementByIdHandler)
	http.HandleFunc("/InsertMovement", movement.InsertMovementHandler)
	http.HandleFunc("/UpdateMovement", movement.UpdateMovementHandler)
	http.HandleFunc("/SearchMovements", movement.SearchMovementsHandler)

	//COST CENTER

	http.HandleFunc("/GetAllCostCenters", costcenter.GetAllCostCentersHandler)
	http.HandleFunc("/GetAllOnlyCostCenters", costcenter.GetAllOnlyCostCentersHandler)
	http.HandleFunc("/GetCostCenerById", costcenter.GetCostCenerByIdHandler)
	http.HandleFunc("/InsertCostCenter", costcenter.InsertCostCenterHandler)
	http.HandleFunc("/UpdateCostCenter", costcenter.UpdateCostCenterHandler)
	http.HandleFunc("/VerifyExistCostCenter", costcenter.VerifyExistCostCenterHandler)
	http.HandleFunc("/SearchCostCenters", costcenter.SearchCostCentersHandler)

	//TESTE
	http.HandleFunc("/teste", loginHandler)
	handler := c.Handler(http.DefaultServeMux)

	balancete.GenerateBalanceteReport(2025, 01, 2025, 06)

	// http.HandleFunc("/createUser", auth.createUser) // Rota de validação do token

	// auth.CreateUser("rico", "654321")

	// Inicia o servidor na porta 8080
	log.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
