package chartofaccount

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

func GetChartOfAccountByID(client *mongo.Client, dbName, collectionName, chartOfAccountID string) (map[string]any, error) {
	collection := client.Database(dbName).Collection(collectionName)

	// Criar filtro para buscar um usuário pelo ID
	fmt.Println("id", chartOfAccountID)

	objectID, erroId := primitive.ObjectIDFromHex(chartOfAccountID)
	if erroId != nil {
		log.Fatalf("Erro ao converter string para ObjectID: %v", erroId)
	}

	filter := bson.M{"_id": objectID}

	// Variável para armazenar o usuário retornado
	var coa models.ChartOfAccount

	// Usar FindOne para pegar apenas um único registro
	err := collection.FindOne(context.Background(), filter).Decode(&coa)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("plano de contas não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar plano de contas: %v", err)
	}

	// Converter o _id para string
	chartOfAccountID = coa.ID.Hex()

	// Retornar o usuário como um mapa
	COAData := map[string]any{
		"ID":          chartOfAccountID, // Agora o campo ID é uma string
		"CodAccount":  coa.CodAccount,
		"Description": coa.Description,
		"Active":      coa.Active,
		"Type":        coa.Type,
		"Year":        coa.Year,
	}

	fmt.Println("COAData", COAData)

	return COAData, nil
}

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

// Função para inserir um usuário na coleção "user"
func InsertChartOfAccount(client *mongo.Client, dbName, collectionName string, user models.ChartOfAccount) error {
	collection := client.Database(dbName).Collection(collectionName)

	// Criar um contexto para a operação de inserção
	ctx := context.Background()

	// Inserir o documento
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("erro ao inserir plano de contas: %v", err)
	}

	return nil
}

// UpdateYearForAllDocuments adiciona um ano ao campo "year" de todos os documentos da coleção,
// caso o ano ainda não esteja presente no array "year".
func UpdateYearForAllDocuments(client *mongo.Client, dbName, collectionName string, newYear int) error {
	collection := client.Database(dbName).Collection(collectionName)

	// Operação de atualização para adicionar o novo ano ao campo "year" de todos os documentos
	update := bson.M{
		"$addToSet": bson.M{
			"year": newYear, // Adiciona o novo ano ao array "year" se não estiver presente
		},
	}

	// Executa a atualização em todos os documentos
	result, err := collection.UpdateMany(context.Background(), bson.M{}, update)
	if err != nil {
		return fmt.Errorf("erro ao adicionar item ao campo 'year' em todos os documentos: %v", err)
	}

	// Verifica se algum documento foi modificado
	if result.ModifiedCount == 0 {
		return fmt.Errorf("nenhum documento foi modificado, possivelmente todos já possuem o ano %d", newYear)
	}

	// Retorna nil se tudo ocorreu com sucesso
	return nil
}
