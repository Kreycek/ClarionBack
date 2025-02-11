package auth

import (
	"Clarion/internal/db"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// Variáveis globais
var jwtKey = []byte("my_secret_key") // Chave secreta para assinatura do JWT

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
	client, err := db.ConnectMongoDB("mongodb://admin:secret@localhost:27017")
	if err != nil {
		http.Error(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
		return
	}
	defer db.CloseMongoDB(client)

	// Obter a coleção de usuários
	collection := db.GetCollection(client, "clarion", "user")

	// Verificar se o usuário existe no banco
	var result bson.M
	err = collection.FindOne(context.Background(), bson.D{{Key: "username", Value: user.Username}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		http.Error(w, "Usuário não encontrado", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar o usuário: %v", err), http.StatusInternalServerError)
		return
	}

	// Comparar a senha fornecida com a senha armazenada no banco de dados
	storedPassword, ok := result["password"].(string)
	if !ok {
		http.Error(w, "Erro na senha armazenada", http.StatusInternalServerError)
		return
	}

	// Verificar se a senha fornecida corresponde à senha armazenada
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password))
	if err != nil {
		http.Error(w, "Senha incorreta", http.StatusUnauthorized)
		return
	}

	// Criar o token JWT
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Criar o token
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Erro ao criar o token", http.StatusInternalServerError)
		return
	}

	// Retornar o token para o cliente
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// Função para validar o token JWT
func ValidateToken(w http.ResponseWriter, r *http.Request) {
	// Obter o token do cabeçalho da requisição
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Token não fornecido", http.StatusUnauthorized)
		return
	}

	// Verificar e validar o token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validar o método de assinatura
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Método de assinatura inválido")
		}
		return jwtKey, nil
	})

	if err != nil {
		http.Error(w, "Token inválido", http.StatusUnauthorized)
		return
	}

	// Validar o token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Token válido. Usuário:", claims["username"])
		w.Write([]byte("Token válido"))
	} else {
		http.Error(w, "Token inválido", http.StatusUnauthorized)
	}
}

func CreateUser(_user User) {
	// Cria uma senha hasheada

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(_user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Erro ao gerar senha hasheada: %v", err)
	}

	// Conectar ao MongoDB
	client, err := db.ConnectMongoDB("mongodb://admin:secret@localhost:27017")
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
	}
	defer db.CloseMongoDB(client)

	// Obter a coleção de usuários
	collection := db.GetCollection(client, "clarion", "user")

	// Inserir o usuário com a senha hasheada
	user := bson.M{
		"name":           _user.Name,
		"lastName":       _user.LastName,
		"email":          _user.Email,
		"passportNumber": _user.PassportNumber,
		"username":       _user.Username,
		"password":       string(hashedPassword),
		"perfil":         []int{1, 2, 3},
	}
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatalf("Erro ao inserir usuário no banco: %v", err)
	}

	log.Println("Usuário inserido com sucesso!")
}

func updateUser(_user User) {
	// Conectar ao MongoDB
	client, err := db.ConnectMongoDB("mongodb://admin:secret@localhost:27017")
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Selecionar o banco de dados e a coleção
	collection := client.Database("meu_banco").Collection("user")

	// Convertendo o ID do usuário para o tipo ObjectID do MongoDB
	objectID, err := primitive.ObjectIDFromHex(_user.Id)
	if err != nil {
		log.Fatalf("Erro ao converter o ID para ObjectID: %v", err)
	}

	// Definir os dados que serão atualizados
	update := bson.M{
		"$set": bson.M{
			"name":           _user.Name,
			"lastName":       _user.LastName,
			"passportNumber": _user.PassportNumber,
			"perfil":         []int{1, 2, 3},
		},
	}

	// Atualizar o documento no MongoDB
	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		log.Fatalf("Erro ao atualizar o documento: %v", err)
	}

	// Exibir o número de documentos modificados
	fmt.Printf("Documento atualizado: %v", result.ModifiedCount)
}
