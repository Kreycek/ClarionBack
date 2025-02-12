package users

import (
	"Clarion/internal/db"
	"Clarion/internal/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Modelo de Usuário com campos do MongoDB

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
