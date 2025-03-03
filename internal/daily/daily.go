package daily

import (
	"Clarion/internal/db"
	"Clarion/internal/models"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Função para obter todos os diários para carregar o drop de buscar
func GetDailysActive(client *mongo.Client, dbName, collectionName string) ([]any, error) {
	collection := db.GetCollection(client, dbName, collectionName)

	filter := bson.M{"active": true}
	// Consultar todos os documentos
	cursor, err := collection.Find(context.Background(), filter)
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

func GetDailyByID(client *mongo.Client, dbName, collectionName, dailyId string) (map[string]any, error) {

	collection := client.Database(dbName).Collection(collectionName)

	// Criar filtro para buscar um usuário pelo ID
	fmt.Println("id", dailyId)

	objectID, erroId := primitive.ObjectIDFromHex(dailyId)
	if erroId != nil {
		log.Fatalf("Erro ao converter string para ObjectID: %v", erroId)
	}

	filter := bson.M{"_id": objectID}

	// Variável para armazenar o usuário retornado
	var daily models.Daily

	// Usar FindOne para pegar apenas um único registro
	err := collection.FindOne(context.Background(), filter).Decode(&daily)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("plano de contas não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar plano de contas: %v", err)
	}

	// Converter o _id para string
	dailyId = daily.ID.Hex()

	// Retornar o usuário como um mapa
	dailys := map[string]any{
		"ID":          dailyId, // Agora o campo ID é uma string
		"CodDaily":    daily.CodDaily,
		"Description": daily.Description,
		"Active":      daily.Active,
		"Documents":   daily.Documents,
	}

	fmt.Println("COAData", dailys)

	return dailys, nil
}

// Função para inserir um usuário na coleção "user"
func InsertDaily(client *mongo.Client, dbName, collectionName string, daily models.Daily) error {
	collection := client.Database(dbName).Collection(collectionName)

	// Criar um contexto para a operação de inserção
	ctx := context.Background()

	// Inserir o documento
	_, err := collection.InsertOne(ctx, daily)
	if err != nil {
		return fmt.Errorf("erro ao inserir plano de contas: %v", err)
	}

	return nil
}

func SearchDailys(client *mongo.Client, dbName, collectionName string, codDaily, description *string, documents []int, page, limit int64) ([]any, int64, error) {
	collection := client.Database(dbName).Collection(collectionName)

	// Criando o filtro dinâmico
	filter := bson.M{}
	if codDaily != nil && *codDaily != "" {
		filter["codDaily"] = bson.M{"$regex": *codDaily, "$options": "i"}
	}
	if description != nil && *description != "" {
		filter["description"] = bson.M{"$regex": *description, "$options": "i"}
	}
	if len(documents) > 0 {
		filter["documents"] = bson.M{"$in": documents}
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
	var dailys []any
	for cursor.Next(context.Background()) {
		var daily models.Daily
		if err := cursor.Decode(&daily); err != nil {
			return nil, 0, fmt.Errorf("erro ao decodificar usuário: %v", err)
		}

		dailys = append(dailys, map[string]any{
			"ID":          daily.ID.Hex(), // Convertendo para string
			"CodDaily":    daily.CodDaily,
			"Description": daily.Description,
			"Documents":   daily.Documents,
			"Active":      daily.Active,
		})
	}

	// Retorna usuários e total de registros
	return dailys, total, nil
}
