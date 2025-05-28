package movement

import (
	clarion "Clarion"
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

func GetAllMovements(client *mongo.Client, dbName, collectionName string, page, limit int) ([]any, int, error) {
	collection := db.GetCollection(client, dbName, collectionName)

	// Criar o filtro (por enquanto vazio, pode ser expandido)
	filter := bson.M{}

	// Obter a contagem total de usuários antes da paginação
	total, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao contar documentos: %v", err)
	}

	// Definindo ordenação por codAccount (1 = ascendente, -1 = descendente)
	sort := bson.D{{Key: "date", Value: -1}}
	// Definir opções de busca com paginação
	options := options.Find()
	options.SetSort(sort)
	options.SetLimit(int64(limit))
	options.SetSkip(int64((page - 1) * limit))

	// Buscar usuários com paginação
	cursor, err := collection.Find(context.Background(), filter, options)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao buscar movimentos: %v", err)
	}
	defer cursor.Close(context.Background())

	var movements []any
	for cursor.Next(context.Background()) {
		var movement models.Movement
		if err := cursor.Decode(&movement); err != nil {
			return nil, 0, fmt.Errorf("erro ao decodificar movimento: %v", err)
		}

		// Adiciona os usuários formatados
		movements = append(movements, map[string]any{
			"ID":              movement.ID.Hex(), // Agora o campo ID é uma string
			"CodDaily":        movement.CodDaily,
			"CodDocument":     movement.CodDocument,
			"Date":            movement.Date,
			"Active":          movement.Active,
			"Month":           movement.Month,
			"Year":            movement.Year,
			"Movements":       movement.Movements,
			"CompanyFullData": movement.CompanyFullData,
			"CompanyId":       movement.CompanyId,
			"CompanyDocument": movement.CompanyDocument,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("erro ao iterar no cursor: %v", err)
	}

	return movements, int(total), nil
}

func GetMovementByID(client *mongo.Client, dbName, collectionName, movementId string) (map[string]any, error) {

	collection := client.Database(dbName).Collection(collectionName)

	// Criar filtro para buscar um usuário pelo ID

	objectID, erroId := primitive.ObjectIDFromHex(movementId)
	if erroId != nil {
		log.Fatalf("Erro ao converter string para ObjectID: %v", erroId)
	}

	filter := bson.M{"_id": objectID}

	// Variável para armazenar o usuário retornado
	var movement models.Movement

	// Usar FindOne para pegar apenas um único registro
	err := collection.FindOne(context.Background(), filter).Decode(&movement)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("plano de contas não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar plano de contas: %v", err)
	}

	// Retornar o usuário como um mapa
	movements := map[string]any{
		"ID":              movement.ID.Hex(), // Agora o campo ID é uma string
		"CodDaily":        movement.CodDaily,
		"CodDocument":     movement.CodDocument,
		"Movements":       movement.Movements,
		"Date":            movement.Date,
		"Active":          movement.Active,
		"CompanyFullData": movement.CompanyFullData,
		"CompanyId":       movement.CompanyId,
		"CompanyDocument": movement.CompanyDocument,
		"Year":            movement.Year,
		"Month":           movement.Month,
	}

	return movements, nil
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 28/05/2025 11:51
Data Final da criação :  28/05/2025 11:54
*/
func GetMovementByCompanyId(client *mongo.Client, dbName, collectionName, companyId string, page, limit int) ([]any, int, error) {
	collection := db.GetCollection(client, dbName, collectionName)

	// Criar o filtro (por enquanto vazio, pode ser expandido)

	objectID, erroId := primitive.ObjectIDFromHex(companyId)
	if erroId != nil {
		log.Fatalf("Erro ao converter string para ObjectID: %v", erroId)
	}

	filter := bson.M{"companyId": objectID}

	// Obter a contagem total de usuários antes da paginação
	total, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao contar documentos: %v", err)
	}

	// Definindo ordenação por codAccount (1 = ascendente, -1 = descendente)
	sort := bson.D{{Key: "date", Value: -1}}
	// Definir opções de busca com paginação
	options := options.Find()
	options.SetSort(sort)
	options.SetLimit(int64(limit))
	options.SetSkip(int64((page - 1) * limit))

	// Buscar usuários com paginação
	cursor, err := collection.Find(context.Background(), filter, options)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao buscar movimentos: %v", err)
	}
	defer cursor.Close(context.Background())

	var movements []any
	for cursor.Next(context.Background()) {
		var movement models.Movement
		if err := cursor.Decode(&movement); err != nil {
			return nil, 0, fmt.Errorf("erro ao decodificar movimento: %v", err)
		}

		// Adiciona os usuários formatados
		movements = append(movements, map[string]any{
			"ID":              movement.ID.Hex(), // Agora o campo ID é uma string
			"CodDaily":        movement.CodDaily,
			"CodDocument":     movement.CodDocument,
			"Date":            movement.Date,
			"Active":          movement.Active,
			"Month":           movement.Month,
			"Year":            movement.Year,
			"Movements":       movement.Movements,
			"CompanyFullData": movement.CompanyFullData,
			"CompanyId":       movement.CompanyId,
			"CompanyDocument": movement.CompanyDocument,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("erro ao iterar no cursor: %v", err)
	}

	return movements, int(total), nil
}

// Função para inserir um usuário na coleção "user"
func InsertMovement(client *mongo.Client, dbName, collectionName string, movement models.Movement) error {
	collection := client.Database(dbName).Collection(collectionName)

	// Criar um contexto para a operação de inserção
	ctx := context.Background()

	// Inserir o documento
	_, err := collection.InsertOne(ctx, movement)
	if err != nil {
		return fmt.Errorf("erro ao inserir plano de contas: %v", err)
	}

	return nil
}

func SearchMovements(client *mongo.Client, dbName, collectionName string, month *int, year *int, CodDaily, CodDocument, date *string, page, limit int64) ([]any, int64, error) {
	collection := client.Database(dbName).Collection(collectionName)

	// Criando o filtro dinâmico
	filter := bson.M{}
	if CodDaily != nil && *CodDaily != "" {
		filter["codDaily"] = *CodDaily // Apenas atribui diretamente
	}

	if CodDocument != nil && *CodDocument != "" {
		// filter["codDocument"] = bson.M{"$regex": *CodDocument, "$options": "i"}
		filter["codDocument"] = *CodDocument // Apenas atribui diretamente
	}

	if month != nil && *month != 0 {

		filter["month"] = *month // Apenas atribui diretamente
	}

	if year != nil && *year != 0 {

		filter["year"] = *year // Apenas atribui diretamente
	}

	// var parsedTime time.Time // Declara parsedTime fora do if
	var err error // Declara err fora do if

	// Contar total de usuários antes da paginação
	total, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}

	sort := bson.D{{Key: "date", Value: -1}}

	// Executa a consulta com paginação
	cursor, err := collection.Find(
		context.Background(),
		filter,
		options.Find().SetSort(sort).SetSkip(int64((page-1)*limit)).SetLimit(int64(limit)),
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	// Processa os resultados
	var movements []any
	for cursor.Next(context.Background()) {
		var movement models.Movement
		if err := cursor.Decode(&movement); err != nil {
			return nil, 0, fmt.Errorf("erro ao decodificar usuário: %v", err)
		}

		movements = append(movements, map[string]any{
			"ID":              movement.ID.Hex(), // Agora o campo ID é uma string
			"CodDaily":        movement.CodDaily,
			"CodDocument":     movement.CodDocument,
			"Movements":       movement.Movements,
			"Date":            movement.Date,
			"Active":          movement.Active,
			"CompanyFullData": movement.CompanyFullData,
			"CompanyId":       movement.CompanyId,
			"CompanyDocument": movement.CompanyDocument,
			"Year":            movement.Year,
			"Month":           movement.Month,
		})
	}

	// Retorna usuários e total de registros
	return movements, total, nil
}

func SearchRecordsBetweenMonths(startYear, startMonth, endYear, endMonth int) ([]models.Movement, error) {

	client, err := db.ConnectMongoDB(clarion.ConectionString)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao MongoDB: %v", err)
	}

	// Acessando a coleção do MongoDB
	collection := client.Database(clarion.DBName).Collection("movement")

	// Consulta para incluir registros dentro do intervalo de meses e anos
	query := bson.M{
		"$or": []bson.M{
			// Caso o intervalo abranja um único ano
			{
				"year":  startYear,
				"month": bson.M{"$gte": startMonth},
			},
			// Caso o intervalo abranja o ano final (final do intervalo de pesquisa)
			{
				"year":  endYear,
				"month": bson.M{"$lte": endMonth},
			},
			// Caso o intervalo se estenda por múltiplos anos (anos intermediários)
			{
				"year": bson.M{"$gt": startYear, "$lt": endYear},
			},
		},
	}

	// Encontrando os registros que atendem à consulta
	cursor, err := collection.Find(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar registros: %v", err)
	}
	defer cursor.Close(context.Background())

	var records []models.Movement
	for cursor.Next(context.Background()) {
		var record models.Movement
		if err := cursor.Decode(&record); err != nil {
			return nil, fmt.Errorf("erro ao decodificar registro: %v", err)
		}

		records = append(records, record)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar os resultados: %v", err)
	}

	return records, nil
}
