package costcenter

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
func GetCostCenter(client *mongo.Client, dbName, collectionName string) ([]any, error) {
	collection := db.GetCollection(client, dbName, collectionName)

	// Consultar todos os documentos
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuários: %v", err)
	}
	defer cursor.Close(context.Background())

	var dadosBanco []any
	for cursor.Next(context.Background()) {
		var cc models.CostCenter
		if err := cursor.Decode(&cc); err != nil {
			return nil, fmt.Errorf("erro ao decodificar ,centro de custo: %v", err)
		}

		// Converter o _id do MongoDB para string para retorno
		ccID := cc.ID
		// Preenche o usuário com o ID convertido em string
		dadosBanco = append(dadosBanco, map[string]any{
			"ID":            ccID, // Agora o campo ID é uma string
			"codCostCenter": cc.CodCostCenter,
			"description":   cc.Description,
			"costCenterSub": cc.CostCenterSub,
			"active":        cc.Active,
		})

	}

	// Verifica se houve algum erro durante a iteração do cursor
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar no cursor: %v", err)
	}

	// Retorna os usuários
	return dadosBanco, nil
}

func GetAllCostCenter(client *mongo.Client, dbName, collectionName string, page, limit int) ([]any, int, error) {
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

	var ccs []any
	for cursor.Next(context.Background()) {
		var cc models.CostCenter
		if err := cursor.Decode(&cc); err != nil {
			return nil, 0, fmt.Errorf("erro ao decodificar centro de custo: %v", err)
		}

		// Adiciona os usuários formatados
		ccs = append(ccs, map[string]any{
			"ID":            cc.ID.Hex(), // Agora o campo ID é uma string
			"CodCostCenter": cc.CodCostCenter,
			"Description":   cc.Description,
			"CostCenterSub": cc.CostCenterSub,
			"Active":        cc.Active,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("erro ao iterar no cursor: %v", err)
	}

	return ccs, int(total), nil
}

func GetCostCenterByID(client *mongo.Client, dbName, collectionName, centerCostId string) (map[string]any, error) {

	collection := client.Database(dbName).Collection(collectionName)

	// Criar filtro para buscar um usuário pelo ID
	fmt.Println("id", centerCostId)

	objectID, erroId := primitive.ObjectIDFromHex(centerCostId)
	if erroId != nil {
		log.Fatalf("Erro ao converter string para ObjectID: %v", erroId)
	}

	filter := bson.M{"_id": objectID}

	// Variável para armazenar o usuário retornado
	var costCenter models.CostCenter

	// Usar FindOne para pegar apenas um único registro
	err := collection.FindOne(context.Background(), filter).Decode(&costCenter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("plano de contas não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar plano de contas: %v", err)
	}

	// Converter o _id para string

	// Retornar o usuário como um mapa
	constCenters := map[string]any{
		"ID":            costCenter.ID.Hex(), // Agora o campo ID é uma string
		"CodCostCenter": costCenter.CodCostCenter,
		"Description":   costCenter.Description,
		"CostCenterSub": costCenter.CostCenterSub,
		"Active":        costCenter.Active,
	}

	fmt.Println("COAData", constCenters)

	return constCenters, nil
}

// Função para inserir um usuário na coleção "user"
func InsertCostCenter(client *mongo.Client, dbName, collectionName string, costCenter models.CostCenter) error {
	collection := client.Database(dbName).Collection(collectionName)

	// Criar um contexto para a operação de inserção
	ctx := context.Background()

	// Inserir o documento
	_, err := collection.InsertOne(ctx, costCenter)
	if err != nil {
		return fmt.Errorf("erro ao inserir centro de custo: %v", err)
	}

	return nil
}

func SearchCostsCenter(client *mongo.Client, dbName, collectionName string, costCenter *string, page, limit int64) ([]any, int64, error) {
	collection := client.Database(dbName).Collection(collectionName)

	// Criando o filtro dinâmico
	filter := bson.M{}
	if costCenter != nil && *costCenter != "" {
		filter["codCostCenter"] = bson.M{"$regex": *costCenter, "$options": "i"}
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
	var costCenters []any
	for cursor.Next(context.Background()) {
		var costCenter models.CostCenter
		if err := cursor.Decode(&costCenter); err != nil {
			return nil, 0, fmt.Errorf("erro ao decodificar usuário: %v", err)
		}

		costCenters = append(costCenters, map[string]any{
			"ID":            costCenter.ID.Hex(), // Agora o campo ID é uma string
			"CodCostCenter": costCenter.CodCostCenter,
			"Description":   costCenter.Description,
			"CostCenterSub": costCenter.CostCenterSub,
			"Active":        costCenter.Active,
		})
	}

	// Retorna usuários e total de registros
	return costCenters, total, nil
}
