package daily

import (
	"Clarion/internal/db"
	"Clarion/internal/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Função para obter todos os diários para carregar o drop de buscar
func GetDailys(client *mongo.Client, dbName, collectionName string) ([]any, error) {
	collection := db.GetCollection(client, dbName, collectionName)

	// Consultar todos os documentos
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuários: %v", err)
	}
	defer cursor.Close(context.Background())

	var dadosBanco []any
	for cursor.Next(context.Background()) {
		var daily models.Daily
		if err := cursor.Decode(&daily); err != nil {
			return nil, fmt.Errorf("erro ao decodificar usuário: %v", err)
		}

		// Converter o _id do MongoDB para string para retorno
		dailyID := daily.ID
		// Preenche o usuário com o ID convertido em string
		dadosBanco = append(dadosBanco, map[string]any{
			"ID":          dailyID, // Agora o campo ID é uma string
			"codDaily":    daily.CodDaily,
			"description": daily.Description,
			"documents":   daily.Documents,
		})

	}

	// Verifica se houve algum erro durante a iteração do cursor
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar no cursor: %v", err)
	}

	// Retorna os usuários
	return dadosBanco, nil
}

func GetAllDaily(client *mongo.Client, dbName, collectionName string, page, limit int) ([]any, int, error) {
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

	var dailys []any
	for cursor.Next(context.Background()) {
		var daily models.Daily
		if err := cursor.Decode(&daily); err != nil {
			return nil, 0, fmt.Errorf("erro ao decodificar usuário: %v", err)
		}

		// Adiciona os usuários formatados
		dailys = append(dailys, map[string]any{
			"ID":          daily.ID.Hex(), // Convertendo para string
			"CodDaily":    daily.CodDaily,
			"Description": daily.Description,
			"Documents":   daily.Documents,
			"Active":      daily.Active,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("erro ao iterar no cursor: %v", err)
	}

	return dailys, int(total), nil
}
