package auth

import (
	clarion "Clarion"
	"Clarion/internal/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id             string `json:"id"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Name           string `json:"name"`
	LastName       string `json:"lastName"`
	Email          string `json:"email"`
	PassportNumber string `json:"passportNumber"`
	Perfil         []int  `json:"perfil"`
}

func formataRetornoHTTP(w http.ResponseWriter, mensagem string, codHttp int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codHttp) // Código 200 OK
	return json.NewEncoder(w).Encode(map[string]string{"message": mensagem})
}

// Variáveis globais

// Função para verificar o nome de usuário e senha
func VerifyUser(w http.ResponseWriter, r *http.Request) {
	// Parse o corpo da requisição
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da requisição", http.StatusBadRequest)
		return
	}

	// Conectar ao MongoDB
	client, err := db.ConnectMongoDB(clarion.ConectionString)
	if err != nil {
		http.Error(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
		return
	}
	defer db.CloseMongoDB(client)

	// Obter a coleção de usuários
	collection := db.GetCollection(client, "clarion", "user")

	filter := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "username", Value: user.Username}},
			bson.D{{Key: "email", Value: user.Username}},
		}},
	}

	// Verificar se o usuário existe no banco
	var result bson.M
	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		formataRetornoHTTP(w, "Erro geral", http.StatusUnauthorized)
		return
		// log.Fatal(err)
	}
	if err == mongo.ErrNoDocuments {
		formataRetornoHTTP(w, "Usuário não encontrado", http.StatusUnauthorized)
		// http.Error(w, "Usuário não encontrado", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar o usuário: %v", err), http.StatusInternalServerError)
		return
	}

	// Comparar a senha fornecida com a senha armazenada no banco de dados
	storedPassword, ok := result["password"].(string)
	if !ok {
		formataRetornoHTTP(w, "Erro na senha armazenada", http.StatusUnauthorized)
		// http.Error(w, "Erro na senha armazenada", http.StatusInternalServerError)
		return
	}

	// Verificar se a senha fornecida corresponde à senha armazenada
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password))
	if err != nil {
		formataRetornoHTTP(w, "Senha incorreta", http.StatusUnauthorized)
		// http.Error(w, "Senha incorreta", http.StatusUnauthorized)
		return
	}

	// Criar o token JWT
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Criar o token
	tokenString, err := token.SignedString(clarion.SecretKey)
	if err != nil {
		formataRetornoHTTP(w, "Erro ao criar o token", 401)
		// http.Error(w, "Erro ao criar o token", http.StatusInternalServerError)
		return
	}

	// Retornar o token para o cliente
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// Função para validar o token JWT
func ValidateToken(w http.ResponseWriter, r *http.Request) {
	// Recuperar o token do cabeçalho 'Authorization'
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		formataRetornoHTTP(w, "Token não fornecido", http.StatusUnauthorized)
		// http.Error(w, "Token não fornecido", http.StatusUnauthorized)
		return
	}

	// Remover o prefixo 'Bearer ' caso esteja presente
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	// Parse e validação do token JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Método de assinatura inválido")
		}
		return clarion.SecretKey, nil // jwtKey é a chave secreta para validação
	})

	if err != nil {
		// formataRetornoHTTP(w, "Token inválido")
		http.Error(w, "Token inválido", http.StatusUnauthorized)
		return
	}

	// Validar o token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Token válido. Usuário:", claims["username"])

		// Enviar uma resposta com o status de sucesso e a mensagem "Token válido"

		formataRetornoHTTP(w, "Token válido", http.StatusOK)
	} else {
		http.Error(w, "Token inválido", http.StatusUnauthorized)
	}
}
