package chartofaccount

import (
	"Clarion/internal/db"
	"Clarion/internal/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Função para obter todos os usuários do banco de dados
func GetAllChartOfAccount(client *mongo.Client, dbName, collectionName string, page, limit int) ([]any, int, error) {
	collection := db.GetCollection(client, dbName, collectionName)

	// Criar o filtro (por enquanto vazio, pode ser expandido)
	filter := bson.M{}

	// Obter a contagem total de usuários antes da paginação
	total, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao contar documentos: %v", err)
	}

	// Definir opções de busca com paginação
	options := options.Find()
	options.SetLimit(int64(limit))
	options.SetSkip(int64((page - 1) * limit))

	// Buscar usuários com paginação
	cursor, err := collection.Find(context.Background(), filter, options)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao buscar usuários: %v", err)
	}
	defer cursor.Close(context.Background())

	var charOfAccounts []any
	for cursor.Next(context.Background()) {
		var charOfAccount models.ChartOfAccount
		if err := cursor.Decode(&charOfAccount); err != nil {
			return nil, 0, fmt.Errorf("erro ao decodificar usuário: %v", err)
		}

		// Adiciona os usuários formatados
		charOfAccounts = append(charOfAccounts, map[string]any{
			"ID":          charOfAccount.ID.Hex(), // Convertendo para string
			"CodAccount":  charOfAccount.CodAccount,
			"Description": charOfAccount.Description,
			"Year":        charOfAccount.Year,
			"Type":        charOfAccount.Type,
			"Active":      charOfAccount.Active,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("erro ao iterar no cursor: %v", err)
	}

	return charOfAccounts, int(total), nil
}

func SearchChartOfAccounts(client *mongo.Client, dbName, collectionName string, codAccount, description *string, _type *string, years []int, page, limit int64) ([]any, int64, error) {
	collection := client.Database(dbName).Collection(collectionName)

	// Criando o filtro dinâmico
	filter := bson.M{}
	if codAccount != nil && *codAccount != "" {
		filter["codAccount"] = bson.M{"$regex": *codAccount, "$options": "i"}
	}
	if description != nil && *description != "" {
		filter["description"] = bson.M{"$regex": *description, "$options": "i"}
	}

	if _type != nil && *_type != "" {
		filter["type"] = bson.M{"$regex": *_type, "$options": "i"}
	}
	if len(years) > 0 {
		filter["year"] = bson.M{"$in": years}
	}

	// Contar total de usuários antes da paginação
	total, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}

	// Executa a consulta com paginação
	cursor, err := collection.Find(
		context.Background(),
		filter,
		options.Find().SetSkip(int64((page-1)*limit)).SetLimit(int64(limit)),
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	// Processa os resultados
	var charOfAccounts []any
	for cursor.Next(context.Background()) {
		var charOfAccount models.ChartOfAccount
		if err := cursor.Decode(&charOfAccount); err != nil {
			return nil, 0, fmt.Errorf("erro ao decodificar usuário: %v", err)
		}

		charOfAccounts = append(charOfAccounts, map[string]any{
			"ID":          charOfAccount.ID.Hex(), // Convertendo para string
			"CodAccount":  charOfAccount.CodAccount,
			"Description": charOfAccount.Description,
			"Year":        charOfAccount.Year,
			"Type":        charOfAccount.Type,
			"Active":      charOfAccount.Active,
		})
	}

	// Retorna usuários e total de registros
	return charOfAccounts, total, nil
}
