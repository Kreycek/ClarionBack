package company

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

func GetAllAutoComplete(client *mongo.Client, dbName, collectionName, name string) ([]any, error) {
	collection := db.GetCollection(client, dbName, collectionName)

	// Criar o filtro para pesquisar pelo nome
	filter := bson.M{}
	if name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"} // Busca insensível a maiúsculas e minúsculas
	}

	// Buscar todas as empresas sem paginação
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar empresas: %v", err)
	}
	defer cursor.Close(context.Background())

	var companies []any
	for cursor.Next(context.Background()) {
		var company models.Company
		if err := cursor.Decode(&company); err != nil {
			return nil, fmt.Errorf("erro ao decodificar empresa: %v", err)
		}

		// Adiciona as empresas formatadas
		companies = append(companies, map[string]any{
			"ID":         company.ID.Hex(), // Convertendo para string
			"CodCompany": company.CodCompany,
			"Name":       company.Name,
			"Documents":  company.Documents,
			"Active":     company.Active,
			"Exercise":   company.Exercise,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar no cursor: %v", err)
	}

	return companies, nil
}

func GetAllCompany(client *mongo.Client, dbName, collectionName string, page, limit int) ([]any, int, error) {
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

	var companys []any
	for cursor.Next(context.Background()) {
		var company models.Company
		if err := cursor.Decode(&company); err != nil {
			return nil, 0, fmt.Errorf("erro ao decodificar usuário: %v", err)
		}

		// Adiciona os usuários formatados
		companys = append(companys, map[string]any{
			"ID":              company.ID.Hex(), // Convertendo para string
			"CodCompany":      company.CodCompany,
			"Name":            company.Name,
			"CAE":             company.CAE,
			"Documents":       company.Documents,
			"MainActivity":    company.MainActivity,
			"OtherActivities": company.OtherActivities,
			"LegalNature":     company.LegalNature,
			"SocialCapital":   company.SocialCapital,
			"NationalCapital": company.NationalCapital,
			"ExtraCapital":    company.ExtraCapital,
			"PublicCapital":   company.PublicCapital,
			"VATRegime":       company.VATRegime,
			"Email":           company.Email,
			"WebSite":         company.WebSite,
			"Active":          company.Active,
			"Phone":           company.Phone,
			"Exercise":        company.Exercise,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("erro ao iterar no cursor: %v", err)
	}

	return companys, int(total), nil
}

func GetCompanyByID(client *mongo.Client, dbName, collectionName, companyId string) (map[string]any, error) {

	collection := client.Database(dbName).Collection(collectionName)

	objectID, erroId := primitive.ObjectIDFromHex(companyId)
	if erroId != nil {
		log.Fatalf("Erro ao converter string para ObjectID: %v", erroId)
	}

	filter := bson.M{"_id": objectID}

	// Variável para armazenar o usuário retornado
	var company models.Company

	// Usar FindOne para pegar apenas um único registro
	err := collection.FindOne(context.Background(), filter).Decode(&company)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("plano de contas não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar plano de contas: %v", err)
	}

	// Converter o _id para string
	companyId = company.ID.Hex()

	// Retornar o usuário como um mapa
	companys := map[string]any{
		"ID":              companyId, // Agora o campo ID é uma string
		"CodCompany":      company.CodCompany,
		"Name":            company.Name,
		"CAE":             company.CAE,
		"Documents":       company.Documents,
		"MainActivity":    company.MainActivity,
		"OtherActivities": company.OtherActivities,
		"LegalNature":     company.LegalNature,
		"SocialCapital":   company.SocialCapital,
		"NationalCapital": company.NationalCapital,
		"ExtraCapital":    company.ExtraCapital,
		"PublicCapital":   company.PublicCapital,
		"VATRegime":       company.VATRegime,
		"Email":           company.Email,
		"WebSite":         company.WebSite,
		"Active":          company.Active,
		"Phone":           company.Phone,
		"Exercise":        company.Exercise,
	}

	return companys, nil
}

// Função para inserir uma empresa na coleção "user"
func InsertCompany(client *mongo.Client, dbName, collectionName string, company models.Company) error {
	collection := client.Database(dbName).Collection(collectionName)

	// Criar um contexto para a operação de inserção
	ctx := context.Background()

	// Inserir o documento
	_, err := collection.InsertOne(ctx, company)
	if err != nil {
		return fmt.Errorf("erro ao inserir a empresa: %v", err)
	}

	return nil
}

func SearchCompany(
	client *mongo.Client,
	dbName, collectionName string,
	codCompany *string,
	document *string,
	address *string,
	page,
	limit int64) ([]any, int64, error) {

	collection := client.Database(dbName).Collection(collectionName)

	// Criando o filtro dinâmico
	filter := bson.M{}
	if codCompany != nil && *codCompany != "" {
		filter["codCompany"] = bson.M{"$regex": *codCompany, "$options": "i"}
	}
	if document != nil && *document != "" {
		filter["documents.documentNumber"] = bson.M{"$regex": *document, "$options": "i"}
	}

	if address != nil && *address != "" {
		filter["documents.address"] = bson.M{"$regex": *address, "$options": "i"}
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
		var company models.Company
		if err := cursor.Decode(&company); err != nil {
			return nil, 0, fmt.Errorf("erro ao decodificar usuário: %v", err)
		}

		dailys = append(dailys, map[string]any{
			"ID":              company.ID.Hex(), // Convertendo para string
			"CodCompany":      company.CodCompany,
			"Name":            company.Name,
			"CAE":             company.CAE,
			"Documents":       company.Documents,
			"MainActivity":    company.MainActivity,
			"OtherActivities": company.OtherActivities,
			"LegalNature":     company.LegalNature,
			"SocialCapital":   company.SocialCapital,
			"NationalCapital": company.NationalCapital,
			"ExtraCapital":    company.ExtraCapital,
			"PublicCapital":   company.PublicCapital,
			"vatRegime":       company.VATRegime,
			"Email":           company.Email,
			"WebSite":         company.WebSite,
			"Active":          company.Active,
			"Phone":           company.Phone,
			"Exercise":        company.Exercise,
		})
	}

	// Retorna usuários e total de registros
	return dailys, total, nil
}
