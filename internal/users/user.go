package users

import (
	"Clarion/internal/db"
	"Clarion/internal/models"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Modelo de Usuário com campos do MongoDB

func GetUserByID(client *mongo.Client, dbName, collectionName, userID string) (map[string]any, error) {
	collection := client.Database(dbName).Collection(collectionName)

	// Criar filtro para buscar um usuário pelo ID
	fmt.Println("id", userID)

	objectID, erroId := primitive.ObjectIDFromHex(userID)
	if erroId != nil {
		log.Fatalf("Erro ao converter string para ObjectID: %v", erroId)
	}

	filter := bson.M{"_id": objectID}

	// Variável para armazenar o usuário retornado
	var user models.User

	// Usar FindOne para pegar apenas um único registro
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("usuário não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar usuário: %v", err)
	}

	// Converter o _id para string
	userID = user.ID.Hex()

	// Retornar o usuário como um mapa
	userData := map[string]any{
		"ID":             userID, // Agora o campo ID é uma string
		"Name":           user.Name,
		"LastName":       user.LastName,
		"Email":          user.Email,
		"PassportNumber": user.PassportNumber,
		"Perfil":         user.Perfil,
		"Username":       user.Username,
		"Active":         user.Active,
	}

	fmt.Println("userData", userData)

	return userData, nil
}

func SearchUsers(client *mongo.Client, dbName, collectionName string, name, email *string, perfis []int) ([]any, error) {
	collection := client.Database(dbName).Collection(collectionName)

	// Criando o filtro dinâmico
	fmt.Println("*name", *name)
	filter := bson.M{}
	if name != nil && *name != "" {
		fmt.Println("name", name)
		filter["name"] = bson.M{"$regex": *name, "$options": "i"}
	}
	if email != nil && *email != "" {
		filter["email"] = bson.M{"$regex": *email, "$options": "i"}
	}
	if len(perfis) > 0 {
		fmt.Println("perfil", perfis)
		filter["perfil"] = bson.M{"$in": perfis} // Busca usuários que tenham pelo menos um dos perfis informados
	}

	// Executa a consulta no MongoDB
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Converte os resultados em uma slice de usuários
	var users []any
	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, fmt.Errorf("erro ao decodificar usuário: %v", err)
		}

		// Converter o _id do MongoDB para string para retorno
		userID := user.ID
		// Preenche o usuário com o ID convertido em string
		users = append(users, map[string]any{
			"ID":             userID, // Agora o campo ID é uma string
			"Name":           user.Name,
			"LastName":       user.LastName,
			"Email":          user.Email,
			"PassportNumber": user.PassportNumber,
			"Perfil":         user.Perfil,
			"Username":       user.Username,
		})

	}

	// var users []models.User
	// if err = cursor.All(context.Background(), &users); err != nil {
	// 	return nil, err
	// }

	return users, nil
}

// Função para obter todos os usuários do banco de dados
func GetAllUsers(client *mongo.Client, dbName, collectionName string) ([]any, error) {
	collection := db.GetCollection(client, dbName, collectionName)

	// Consultar todos os documentos
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuários: %v", err)
	}
	defer cursor.Close(context.Background())

	var users []any
	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, fmt.Errorf("erro ao decodificar usuário: %v", err)
		}

		// Converter o _id do MongoDB para string para retorno
		userID := user.ID
		// Preenche o usuário com o ID convertido em string
		users = append(users, map[string]any{
			"ID":             userID, // Agora o campo ID é uma string
			"Name":           user.Name,
			"LastName":       user.LastName,
			"Email":          user.Email,
			"PassportNumber": user.PassportNumber,
			"Perfil":         user.Perfil,
			"Username":       user.Username,
		})

	}

	// Verifica se houve algum erro durante a iteração do cursor
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar no cursor: %v", err)
	}

	// Retorna os usuários
	return users, nil
}

// Função para inserir um usuário na coleção "user"
func InsertUser(client *mongo.Client, dbName, collectionName string, user models.User) error {
	collection := client.Database(dbName).Collection(collectionName)

	// Criar um contexto para a operação de inserção
	ctx := context.Background()

	// Inserir o documento
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("erro ao inserir usuário: %v", err)
	}

	return nil
}
